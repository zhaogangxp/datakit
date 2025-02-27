package cmds

import (
	"fmt"
	"reflect"
	"time"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/config"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
)

var (
	failed  = 0
	unknown = 0
	passed  = 0
	ignored = 0
)

func checkInputCfg(tpl *ast.Table, fp string) {
	var err error

	if len(tpl.Fields) == 0 {
		warnf("[W] no content in %s\n", fp)
		ignored++
		return
	}

	for field, node := range tpl.Fields {
		switch field {
		default:
			infof("[I] ignore config %s\n", fp)
			ignored++
			return

		case "inputs": //nolint:goconst
			stbl, ok := node.(*ast.Table)
			if !ok {
				l.Warnf("ignore bad toml node within %s", fp)
			} else {
				for inputName, v := range stbl.Fields {
					if c, ok := inputs.Inputs[inputName]; !ok {
						warnf("[W] unknown input `%s' found in %s\n", inputName, fp)
						unknown++
					} else {
						if _, err = config.TryUnmarshal(v, inputName, c); err != nil {
							errorf("[E] failed to init input %s from %s:\n%s\n", inputName, fp, err.Error())
							failed++
						} else {
							if FlagVVV {
								output("[OK] %s/%s\n", inputName, fp)
							}
							passed++
						}
					}
				}
			}
		}
	}
}

// check samples of every inputs.
func checkSample() error {
	start := time.Now()
	failed = 0
	unknown = 0
	passed = 0
	ignored = 0

	for k, c := range inputs.Inputs {
		i := c()

		if k == datakit.DatakitInputName {
			warnf("[W] ignore self input\n")
			ignored++
			continue
		}

		tpl, err := toml.Parse([]byte(i.SampleConfig()))
		if err != nil {
			errorf("[E] failed to parse %s: %s", k, err.Error())
			failed++
		} else {
			checkInputCfg(tpl, k)
		}
	}

	infof("\n------------------------\n")
	infof("checked %d sample, %d ignored, %d passed, %d failed, %d unknown, ",
		len(inputs.Inputs), ignored, passed, failed, unknown)

	infof("cost %v\n", time.Since(start))

	if failed > 0 {
		return fmt.Errorf("load %d sample failed", failed)
	}
	return nil
}

func checkConfig(dir, suffix string) error {
	start := time.Now()
	fps := config.SearchDir(dir, suffix)

	failed = 0
	unknown = 0
	passed = 0
	ignored = 0

	for _, fp := range fps {
		tpl, err := config.ParseCfgFile(fp)
		if err != nil {
			errorf("[E] failed to parse %s: %s, %s", fp, err.Error(), reflect.TypeOf(err))
			failed++
		} else {
			checkInputCfg(tpl, fp)
		}
	}

	infof("\n------------------------\n")
	infof("checked %d conf, %d ignored, %d passed, %d failed, %d unknown, ",
		len(fps), ignored, passed, failed, unknown)

	infof("cost %v\n", time.Since(start))

	if failed > 0 {
		return fmt.Errorf("load %d conf failed", failed)
	}

	return nil
}
