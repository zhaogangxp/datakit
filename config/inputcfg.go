package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
)

var DisableSelfInput bool

func SearchDir(dir string, suffix string) []string {
	fps := []string{}

	if err := filepath.Walk(dir, func(fp string, f os.FileInfo, err error) error {
		if err != nil {
			l.Errorf("walk on %s failed: %s", fp, err)
			return nil
		}

		if f == nil {
			l.Warnf("nil FileInfo on %s", fp)
			return nil
		}

		if f.IsDir() {
			l.Debugf("ignore dir %s", fp)
			return nil
		}

		if suffix == "" || strings.HasSuffix(f.Name(), suffix) {
			fps = append(fps, fp)
		}
		return nil
	}); err != nil {
		l.Error(err)
	}

	return fps
}

func GetGitRepoDir(cloneDirName string) (string, error) {
	if cloneDirName == "" {
		// you shouldn't be here, check before you call this function.
		return "", fmt.Errorf("git_repo_clone_dir_empty")
	}
	return filepath.Join(datakit.GitReposDir, cloneDirName), nil
}

// LoadInputsConfigEx load all inputs under @InstallDir/conf.d.
func LoadInputsConfigEx(confRootPath string) map[string]*ast.Table {
	availableInputCfgs := map[string]*ast.Table{}

	confs := SearchDir(confRootPath, ".conf")

	l.Debugf("loading %d conf...", len(confs))

	for _, fp := range confs {
		l.Debugf("loading conf %s...", fp)

		if filepath.Base(fp) == "datakit.conf" {
			continue
		}

		tbl, err := ParseCfgFile(fp)
		if err != nil {
			l.Warnf("parse conf %s failed: %s, ignored", fp, err)
			continue
		}

		deprecates := checkDepercatedInputs(tbl, deprecatedInputs)
		if len(deprecates) > 0 {
			for k, v := range deprecates {
				l.Warnf("input `%s' removed, please use %s instead", k, v)
			}
		}

		if len(tbl.Fields) == 0 {
			l.Warnf("no conf available on %s", fp)
			continue
		}

		l.Debugf("parse %s ok", fp)
		availableInputCfgs[fp] = tbl
	}

	return availableInputCfgs
}

// fp == "", add new when not exist, set ConfigPaths empty when exist.
func addConfigInfoPath(inputName string, fp string, loaded int8) {
	if c, ok := inputs.ConfigInfo[inputName]; ok {
		if len(fp) == 0 {
			c.ConfigPaths = []*inputs.ConfigPathStat{} // set empty for reload datakit
			return
		}
		for _, p := range c.ConfigPaths {
			if p.Path == fp {
				p.Loaded = loaded
				return
			}
		}
		c.ConfigPaths = append(c.ConfigPaths, &inputs.ConfigPathStat{Loaded: loaded, Path: fp})
	} else {
		creator, ok := inputs.Inputs[inputName]
		if ok {
			config := &inputs.Config{
				ConfigPaths:  []*inputs.ConfigPathStat{},
				SampleConfig: creator().SampleConfig(),
				Catalog:      creator().Catalog(),
				ConfigDir:    datakit.ConfdDir,
			}
			if len(fp) > 0 {
				config.ConfigPaths = append(config.ConfigPaths, &inputs.ConfigPathStat{Loaded: loaded, Path: fp})
			}
			inputs.ConfigInfo[inputName] = config
		}
	}
}

func doLoadInputConf(name string, creator inputs.Creator, inputcfgs map[string]*ast.Table) error {
	l.Debugf("search input cfg for %s", name)

	list := searchDatakitInputCfg(inputcfgs, name, creator)

	for _, i := range list {
		inputs.AddInput(name, i)
	}

	return nil
}

func searchDatakitInputCfg(inputcfgs map[string]*ast.Table,
	name string,
	creator inputs.Creator) []inputs.Input {
	inputlist := []inputs.Input{}

	addConfigInfoPath(name, "", 0) // init config info

	for fp, tbl := range inputcfgs {
		for field, node := range tbl.Fields {
			switch field {
			case "inputs": //nolint:goconst
				stbl, ok := node.(*ast.Table)
				if !ok {
					l.Warnf("ignore bad toml node for %s within %s", name, fp)
					addConfigInfoPath(name, fp, 0)
				} else {
					for inputName, v := range stbl.Fields {
						if inputName != name {
							continue
						}
						lst, err := TryUnmarshal(v, inputName, creator)
						if err != nil {
							l.Warnf("unmarshal input %s failed within %s: %s", inputName, fp, err.Error())
							addConfigInfoPath(name, fp, 0)
							continue
						}

						l.Infof("load input %s from %s ok", inputName, fp)

						// dca config path
						addConfigInfoPath(name, fp, 1)

						inputlist = append(inputlist, lst...)
					}
				}

			default: // compatible with old version: no [[inputs.xxx]] header
				l.Debugf("ignore field %s in file %s", field, fp)
			}
		}
	}

	return inputlist
}

func isDisabled(wlists, blists []*inputHostList, hostname, name string) bool {
	for _, bl := range blists {
		if bl.MatchHost(hostname) && bl.MatchInput(name) {
			return true // 一旦上榜，无脑屏蔽
		}
	}

	// 如果采集器在白名单中，但对应的 host 不在白名单，则屏蔽掉
	// 如果采集器在白名单中，对应的 host 在白名单，放行
	// 如果采集器不在白名单中，不管 host 情况，一律放行
	if len(wlists) > 0 {
		for _, wl := range wlists {
			if wl.MatchInput(name) { // 说明@name有白名单限制
				if wl.MatchHost(hostname) {
					return false
				} else { // 不在白名单中的 host，屏蔽掉
					return true
				}
			}
		}
	}
	return false
}

func TryUnmarshal(tbl interface{}, name string, creator inputs.Creator) (inputList []inputs.Input, err error) {
	var tbls []*ast.Table

	switch t := tbl.(type) {
	case []*ast.Table:
		tbls = t
	case *ast.Table:
		tbls = append(tbls, t)
	default:
		err = fmt.Errorf("invalid toml format on %s: %v", name, t)
		return
	}

	for _, t := range tbls {
		input := creator()
		err = toml.UnmarshalTable(t, input)
		if err != nil {
			l.Warnf("toml unmarshal %s failed: %v", name, err)
			continue
		}

		inputList = append(inputList, input)
	}

	return
}

//nolint:lll
var confsampleFingerprint = append([]byte(fmt.Sprintf(`# {"version": "%s", "desc": "do NOT edit this line"}`, datakit.Version)), byte('\n'))

func initDatakitConfSample(name string, c inputs.Creator) error {
	if name == datakit.DatakitInputName {
		return nil
	}

	input := c()
	catalog := input.Catalog()

	cfgpath := filepath.Join(datakit.ConfdDir, catalog, name+".conf.sample")
	l.Debugf("create datakit conf path %s", filepath.Join(datakit.ConfdDir, catalog))
	if err := os.MkdirAll(filepath.Join(datakit.ConfdDir, catalog), datakit.ConfPerm); err != nil {
		l.Errorf("create catalog dir %s failed: %s", catalog, err.Error())
		return err
	}

	sample := input.SampleConfig()
	if sample == "" {
		return fmt.Errorf("no sample available on collector %s", name)
	}

	// 在 conf-sample 头部增加一些指纹信息.
	// 一般用户在编辑 conf 时，都是 copy 这个 sample 的。如果 sample 中带上指纹，
	// 那么最终的配置上也会带上这可能便于后续的升级，即升级程序能识别某个 conf
	// 的版本，进而进行指定的升级
	if err := ioutil.WriteFile(cfgpath, append(confsampleFingerprint, []byte(sample)...), datakit.ConfPerm); err != nil {
		l.Errorf("failed to create sample configure for collector %s: %s", name, err.Error())
		return err
	}

	return nil
}

// Creata datakit input plugin's configures if not exists.
func initPluginSamples() error {
	for name, create := range inputs.Inputs {
		if !datakit.Enabled(name) {
			l.Debugf("initPluginSamples: ignore unchecked input %s", name)
			continue
		}

		if err := initDatakitConfSample(name, create); err != nil {
			return err
		}
	}
	return nil
}

func initDefaultEnabledPlugins(c *Config) {
	if len(c.DefaultEnabledInputs) == 0 {
		l.Debug("no default inputs enabled")
		return
	}

	for _, name := range c.DefaultEnabledInputs {
		l.Debugf("init default input %s conf...", name)

		var fpath, sample string

		if c, ok := inputs.Inputs[name]; ok {
			i := c()
			sample = i.SampleConfig()

			fpath = filepath.Join(datakit.ConfdDir, i.Catalog(), name+".conf")
		} else {
			l.Warnf("input %s not found, ignored", name)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), datakit.ConfPerm); err != nil {
			l.Errorf("mkdir failed: %s, ignored", err.Error())
			continue
		}

		// check exist
		if _, err := os.Stat(fpath); err == nil {
			continue
		}

		if err := ioutil.WriteFile(fpath, []byte(sample), datakit.ConfPerm); err != nil {
			l.Errorf("write input %s config failed: %s, ignored", name, err.Error())
			continue
		}

		l.Infof("enable input %s ok", name)
	}
}

func LoadInputConfigFile(f string, creator inputs.Creator) ([]inputs.Input, error) {
	tbl, err := ParseCfgFile(f)
	if err != nil {
		return nil, fmt.Errorf("parse conf failed: %w", err)
	}

	return parseTableToInputs(tbl, creator)
}

func LoadInputConfig(data string, creator inputs.Creator) ([]inputs.Input, error) {
	tbl, err := toml.Parse([]byte(data))
	if err != nil {
		l.Errorf("parse toml %s failed", data)
		return nil, fmt.Errorf("[error] parse conf failed: %w", err)
	}

	return parseTableToInputs(tbl, creator)
}

func parseTableToInputs(tbl *ast.Table, creator inputs.Creator) ([]inputs.Input, error) {
	var inputlist []inputs.Input
	var err error

	for field, node := range tbl.Fields {
		switch field {
		case "inputs": //nolint:goconst
			stbl, ok := node.(*ast.Table)
			if ok {
				for inputName, v := range stbl.Fields {
					inputlist, err = TryUnmarshal(v, inputName, creator)
					if err != nil {
						return nil, fmt.Errorf("unmarshal input failed, %w", err)
					}
				}
			}

		default: // compatible with old version: no [[inputs.xxx]] header
			inputlist, err = TryUnmarshal(node, "", creator)
			if err != nil {
				return nil, fmt.Errorf("unmarshal input failed: %w", err)
			}
		}
	}

	return inputlist, nil
}

var deprecatedInputs = map[string]string{
	"dockerlog":         "docker",
	"docker_containers": "docker",
	"traceSkywalking":   "skywalking",
	"traceJaeger":       "jaeger",
}

func checkDepercatedInputs(tbl *ast.Table, entries map[string]string) (res map[string]string) {
	res = map[string]string{}

	for _, node := range tbl.Fields {
		stbl, ok := node.(*ast.Table)
		if !ok {
			continue
		}
		for inputName := range stbl.Fields {
			if x, ok := entries[inputName]; ok {
				res[inputName] = x
			}
		}
	}
	return
}
