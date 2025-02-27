// Package datakit defined all datakit's global settings
package datakit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"gitlab.jiagouyun.com/cloudcare-tools/cliutils"
	"gitlab.jiagouyun.com/cloudcare-tools/cliutils/logger"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/git"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/goroutine"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/service"
)

const (
	OSWindows = `windows`
	OSLinux   = `linux`
	OSDarwin  = `darwin`

	OSArchWinAmd64    = "windows/amd64"
	OSArchWin386      = "windows/386"
	OSArchLinuxArm    = "linux/arm"
	OSArchLinuxArm64  = "linux/arm64"
	OSArchLinux386    = "linux/386"
	OSArchLinuxAmd64  = "linux/amd64"
	OSArchDarwinAmd64 = "darwin/amd64"

	CommonChanCap = 32

	// data category, aka API /v1/write/category.
	MetricDeprecated = "/v1/write/metrics"
	Metric           = "/v1/write/metric"
	Network          = "/v1/write/network"
	KeyEvent         = "/v1/write/keyevent"
	Object           = "/v1/write/object"
	CustomObject     = "/v1/write/custom_object"
	Logging          = "/v1/write/logging"
	Tracing          = "/v1/write/tracing"
	Rum              = "/v1/write/rum"
	Security         = "/v1/write/security"

	// other APIS.
	HeartBeat         = "/v1/write/heartbeat"
	Election          = "/v1/election"
	ElectionHeartbeat = "/v1/election/heartbeat"
	QueryRaw          = "/v1/query/raw"
	Workspace         = "/v1/workspace"
	ObjectLabel       = "/v1/object/labels" // object label
	LogUpload         = "/v1/log"
	LogFilter         = "/v2/logfilter/pull"
	ListDataWay       = "/v2/list/dataway"
	DatakitInputName  = "self"
)

var (
	Exit = cliutils.NewSem()
	WG   = sync.WaitGroup{}

	Docker     = false
	Version    = git.Version
	AutoUpdate = false

	InstallDir = optionalInstallDir[runtime.GOOS+"/"+runtime.GOARCH]

	DatakitHostName = "" // 默认为空

	optionalInstallDir = map[string]string{
		OSArchWinAmd64: `C:\Program Files\datakit`,
		OSArchWin386:   `C:\Program Files (x86)\datakit`,

		OSArchLinuxArm:    `/usr/local/datakit`,
		OSArchLinuxArm64:  `/usr/local/datakit`,
		OSArchLinuxAmd64:  `/usr/local/datakit`,
		OSArchLinux386:    `/usr/local/datakit`,
		OSArchDarwinAmd64: `/usr/local/datakit`,
	}

	AllOS   = []string{OSWindows, OSLinux, OSDarwin}
	AllArch = []string{
		OSArchWinAmd64, OSArchWin386, OSArchLinuxArm,
		OSArchLinuxArm64, OSArchLinux386, OSArchLinuxAmd64, OSArchDarwinAmd64,
	}

	UnknownOS   = []string{"unknown"}
	UnknownArch = []string{"unknown"}

	DataDir  = filepath.Join(InstallDir, "data")
	ConfdDir = filepath.Join(InstallDir, "conf.d")

	GitReposDir      = filepath.Join(InstallDir, "gitrepos")
	GetReposConfDirs []string // git repos conf search dirs

	PythonDDir    = filepath.Join(InstallDir, "python.d")
	PythonCoreDir = filepath.Join(PythonDDir, "core")

	MainConfPathDeprecated = filepath.Join(InstallDir, "datakit.conf")
	MainConfPath           = filepath.Join(ConfdDir, "datakit.conf")

	pidFile = filepath.Join(InstallDir, ".pid")

	PipelineDir        = filepath.Join(InstallDir, "pipeline")
	PipelinePatternDir = filepath.Join(PipelineDir, "pattern")
	CacheDir           = filepath.Join(InstallDir, "cache")
	GRPCDomainSock     = filepath.Join(InstallDir, "datakit.sock")
	GRPCSock           = ""
)

func SetWorkDir(dir string) {
	InstallDir = dir

	DataDir = filepath.Join(InstallDir, "data")
	ConfdDir = filepath.Join(InstallDir, "conf.d")

	MainConfPathDeprecated = filepath.Join(InstallDir, "datakit.conf")
	MainConfPath = filepath.Join(ConfdDir, "datakit.conf")

	PipelineDir = filepath.Join(InstallDir, "pipeline")
	PipelinePatternDir = filepath.Join(PipelineDir, "pattern")
	CacheDir = filepath.Join(InstallDir, "cache")
	GRPCDomainSock = filepath.Join(InstallDir, "datakit.sock")
	pidFile = filepath.Join(InstallDir, ".pid")

	GitReposDir = filepath.Join(InstallDir, "gitrepos")
	PythonDDir = filepath.Join(InstallDir, "python.d")
	PythonCoreDir = filepath.Join(PythonDDir, "core")

	InitDirs()
}

func InitDirs() {
	for _, dir := range []string{
		DataDir,
		ConfdDir,
		PipelineDir,
		PipelinePatternDir,
		GitReposDir,
	} {
		if err := os.MkdirAll(dir, ConfPerm); err != nil {
			l.Fatalf("create %s failed: %s", dir, err)
		}
	}
}

const (
	ConfPerm = os.ModePerm
)

var (
	// goroutines caches  goroutine.
	goroutines = []*goroutine.Group{}

	l = logger.DefaultSLogger("datakit")
)

func SetLog() {
	l = logger.SLogger("datakit")
}

// G create a goroutine group, with namespace datakit.
func G(name string) *goroutine.Group {
	panicCb := func(b []byte) bool {
		l.Errorf("%s", b)
		select {
		case <-Exit.Wait(): // don't continue when exit
			return false
		default:
			return true
		}
	}

	gName := "datakit_" + name
	opt := goroutine.Option{Name: gName, PanicTimes: 6, PanicCb: panicCb, PanicTimeout: 10 * time.Millisecond}
	g := goroutine.NewGroup(opt)
	var mu sync.Mutex
	mu.Lock()
	goroutines = append(goroutines, g)
	mu.Unlock()
	return g
}

// GWait wait all goroutine group exit.
func GWait() {
	for _, g := range goroutines {
		// just ignore error
		_ = g.Wait()
		l.Infof("goroutine Group %s exit", g.Name())
	}
	l.Info("all goroutine group exit successfully")
}

func Quit() {
	_ = os.Remove(pidFile)

	Exit.Close()
	WG.Wait()
	GWait()
	service.Stop()
}

func PID() (int, error) {
	if x, err := ioutil.ReadFile(filepath.Clean(pidFile)); err != nil {
		return -1, err
	} else {
		if pid, err := strconv.ParseInt(string(x), 10, 32); err != nil {
			return -1, err
		} else {
			return int(pid), nil
		}
	}
}

func SavePid() error {
	if isRuning() {
		return fmt.Errorf("datakit still running, PID: %s", pidFile)
	}

	pid := os.Getpid()
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), os.ModePerm)
}

func isRuning() bool {
	var oidPid int64
	var name string
	var p *process.Process

	cont, err := ioutil.ReadFile(filepath.Clean(pidFile))
	// pid文件不存在
	if err != nil {
		return false
	}

	oidPid, err = strconv.ParseInt(string(cont), 10, 32)
	if err != nil {
		return false
	}

	p, _ = process.NewProcess(int32(oidPid))
	name, _ = p.Name()

	return name == getBinName()
}

func getBinName() string {
	bin := "datakit"

	if runtime.GOOS == OSWindows {
		bin += ".exe"
	}

	return bin
}
