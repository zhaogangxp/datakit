{{.CSS}}

- DataKit 版本：{{.Version}}
- 文档发布日期：{{.ReleaseDate}}
- 操作系统支持：全平台

# 简介

本文介绍 DataKit 的基本安装。

## 注册/登陆观测云

浏览器访问 [观测云注册入口](https://auth.guance.com/redirectpage/register)，填写对应信息之后，即可[登陆观测云](https://console.guance.com/pageloading/login)

## 获取安装命令

登陆工作空间，点击左侧「集成」选择顶部「Datakit」，即可看到各种平台的安装命令。

> 注意，以下 Linux/Mac/Windows 安装程序，能自动识别硬件平台（arm/x86, 32bit/64bit），无需做硬件平台选择。

### Linux/Mac

命令大概如下：

```shell
DK_DATAWAY=https://openway.guance.com?token=<TOKEN> bash -c "$(curl -L https://static.guance.com/datakit/install.sh)"
```

安装完成后，在终端会看到安装成功的提示。

#### Mac 安装注意事项

Mac 上安装时，如果安装/升级过程中出现

```shell
"launchctl" failed with stderr: /Library/LaunchDaemons/cn.dataflux.datakit.plist: Service is disabled
# 或者
"launchctl" failed with stderr: /Library/LaunchDaemons/com.guance.datakit.plist: Service is disabled
```

执行

```shell
sudo launchctl enable system/datakit
```

然后再执行如下命令即可

```shell
sudo launchctl load -w /Library/LaunchDaemons/cn.dataflux.datakit.plist
# 或者
sudo launchctl load -w /Library/LaunchDaemons/com.guance.datakit.plist
```

### Windows

> Windows 上安装需在 Powershell 命令行安装，且必须以管理员身份运行 Powershell。按下 Windows 键，输入 powershell 即可看到弹出的 powershell 图标，右键选择「以管理员身份运行」即可。

```powershell
$env:DK_DATAWAY="https://openway.guance.com?token=<TOKEN>"; Set-ExecutionPolicy Bypass -scope Process -Force; Import-Module bitstransfer; start-bitstransfer -source https://static.guance.com/datakit/install.ps1 -destination .install.ps1; powershell .install.ps1;
```

### 额外支持的安装变量

安装脚本支持的环境变量如下（全平台支持）：

- `DK_DATAWAY`：指定 dataway 地址，含 `TOKEN`
- `DK_CLOUD_PROVIDER`：支持安装阶段填写云厂商(`aliyun/aws/tencent/hwcloud/azure`)
- `DK_NAMESPACE`：支持安装阶段指定命名空间(选举用)
- `DK_GLOBAL_TAGS`：支持安装阶段填写全局 tag，格式范例：`project=abc,owner=张三`（多个 tag 之间以英文逗号分隔）
- `DK_HTTP_LISTEN`：支持安装阶段指定 DataKit HTTP 服务绑定的网卡（默认 `localhost`）
- `DK_HTTP_PORT`：支持安装阶段指定 DataKit HTTP 服务绑定的端口（默认 `9529`）
- `DK_DCA_ENABLE`：支持安装阶段开启 DCA 服务（默认未开启）
- `DK_DCA_LISTEN`：支持安装阶段自定义配置 DCA 服务的监听地址和端口（默认`0.0.0.0:9531`）
- `DK_DCA_WHITE_LIST`: 支持安装阶段设置访问 DCA 服务白名单，多个白名单以 `,` 分割 (如：`192.168.0.1/24,10.10.0.1/24`)
- `DK_INSTALL_ONLY`：仅安装，不运行
- `DK_HOSTNAME`:支持安装阶段自定义配置主机名
- `DK_DEF_INPUTS`：默认开启的采集器列表，格式范例：`input1,input2,input3`
- `DK_UPGRADE`：升级到最新版本（注：一旦开启该选项，其它选项均无效）
- `DK_INSTALLER_BASE_URL`：可选择不同环境的安装脚本，默认为 `https://static.guance.com/datakit`
- `HTTPS_PROXY`：通过 Datakit 代理安装
- `DK_PROXY_TYPE`：代理类型。选项有: "datakit" 或 "nginx"，均为小写
- `DK_NGINX_IP`：代理服务器 IP 地址（只需要填 IP 不需要填端口）。这个与上面的 "HTTP_PROXY" 和 "HTTPS_PROXY" 互斥，而且优先级最高，会覆盖以上两者
- `DK_GIT_URL`: 管理配置文件的远程 git repo 地址。（如 `http://username:password@github.com/username/repository.git`）
- `DK_GIT_KEY_PATH`: 本地 PrivateKey 的全路径。（如 `/Users/username/.ssh/id_rsa`）
- `DK_GIT_KEY_PW`: 本地 PrivateKey 的使用密码。（如 `passwd`）
- `DK_GIT_BRANCH`: 指定拉取的分支。<stong>为空则是默认</strong>，默认是远程指定的主分支，一般是 `master`。
- `DK_GIT_INTERVAL`: 定时拉取的间隔。（如 `1m`）
- `DK_LIMIT_CPUMAX`:Linux系统下支持CPU的最大功率,默认30.0。（仅Linux平台支持）
- `DK_LIMIT_CPUMIN`:Linux系统下支持CPU的最小功率,默认5.0。（仅Linux平台支持）
- `DK_CGROUP_ENABLED`:Linux系统下开启Cgroup功能。（如`false`，默认为true，仅Linux平台支持）

如果需要增加环境变量，在 `DK_DATAWAY` 前面追加即可。如追加 `DK_NAMESPACE` 设置：

```shell
# Linux/Mac
DK_NAMESPACE="<namespace>" DK_DATAWAY="https://openway.guance.com?token=<TOKEN>" bash -c "$(curl -L https://static.guance.com/datakit/install.sh)"

# Windows
$env:DK_NAMESPACE="<namespace>"; $env:DK_DATAWAY="https://openway.guance.com?token=<TOKEN>"; Set-ExecutionPolicy Bypass -scope Process -Force; Import-Module bitstransfer; start-bitstransfer -source https://static.guance.com/datakit/install.ps1 -destination .install.ps1; powershell .install.ps1;
```

俩种环境变量的设置格式为：

```shell
# Windows: 多个环境变量之间以分号分割
$env:NAME1="value1"; $env:Name2="value2"

# Linux/Mac: 多个环境变量之间以空格分割
NAME1="value1" NAME2="value2"
```

#### 如何应付不友好的主机名

由于 DataKit 使用主机名（Hostname）作为数据串联的依据，某些情况下，一些主机名取得不是很友好，比如 `iZbp141ahn....`，但由于某些原因，又不能修改这些主机名，这给使用带来一定的困扰。在 DataKit 中，可在主配置中覆盖这个不友好的主机名。

在 `datakit.conf` 中，修改如下配置，DataKit 将读取 `ENV_HOSTNAME` 来覆盖当前的真实主机名：

```toml
[environments]
	ENV_HOSTNAME = "your-fake-hostname-for-datakit"
```

> 注意：如果之前某个主机已经采集了一段时间的数据，更改主机名后，这些历史数据将不再跟新的主机名关联。更改主机名，相当于新增了一台全新的主机。

其它相关链接：

- 关于 DataKit 的基本使用，参考 [DataKit 使用入门](datakit-service-how-to)
