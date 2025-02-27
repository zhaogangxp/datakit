
{{.CSS}}

- DataKit 版本：{{.Version}}
- 文档发布日期：{{.ReleaseDate}}
- 操作系统支持：全平台

# 简介

DataKit 支持手动更新和自动更新两种方式。

## 前置条件

- 自动更新要求 DataKit 版本 >= 1.1.6-rc1
- 手动更新暂无版本要求

## 手动更新

直接执行如下命令查看当前 DataKit 版本。如果线上有最新版本，则会提示对应的更新命令，如：

```shell
# Linux/Mac
$ datakit --version

       Version: 1.1.6-rc0
        Commit: d1f4604d
        Branch: dk-api
 Build At(UTC): 2021-05-11 11:07:06
Golang Version: go version go1.15.8 darwin/amd64
      Uploader: xxxxxxxxxxxxx/xxxxxxx/xxxxxxx
ReleasedInputs: all
---------------------------------------------------

Online version available: 1.1.8-rc1.1, commit 658339b6eb (release at 2021-08-13 05:52:17)

Upgrade:
    DK_UPGRADE=1 bash -c "$(curl -L https://static.guance.com/datakit/install.sh)"

# Windows

$ datakit.exe --version

       Version: 1.1.6-rc0
        Commit: d1f4604d
        Branch: dk-api
 Build At(UTC): 2021-05-11 11:07:06
Golang Version: go version go1.15.8 darwin/amd64
      Uploader: xxxxxxxxxxxxx/xxxxxxx/yyyyyyy
ReleasedInputs: all
---------------------------------------------------

Online version available: 1.1.8-rc1.1, commit 658339b6eb (release at 2021-08-13 05:52:17)

Upgrade:
    $env:DK_UPGRADE="1"; Set-ExecutionPolicy Bypass -scope Process -Force; Import-Module bitstransfer; start-bitstransfer -source https://static.guance.com/datakit/install.ps1 -destination .install.ps1; powershell .install.ps1;
```

> 如果当前 DataKit 处于被代理模式，自动更新的提示命令中，会自动加上代理设置：

```shell
# Linux/Mac
HTTPS_PROXY=http://10.100.64.198:9530 DK_UPGRADE=1 ...

# Windows
$env:HTTPS_PROXY="http://10.100.64.198:9530"; $env:DK_UPGRADE="1" ...
```

## 自动更新

在 Linux 中，为便于 DataKit 实现自动更新，可通过 crontab 方式添加任务，实现定期更新。

> 注：目前自动更新只支持 Linux，且暂不支持代理模式。

### 准备更新脚本

将如下脚本内容复制到 DataKit 所在机器的安装目录下，保存 `datakit-update.sh`（名称随意）

```bash
#!/bin/bash
# Update DataKit if new version available

otalog=/usr/local/datakit/ota-update.log
installer=https://static.guance.com/datakit/installer-linux-amd64

# 注意：如果不希望更新 RC 版本的 DataKit，可移除 `--accept-rc-version`
/usr/local/datakit/datakit --check-update --accept-rc-version --update-log $otalog

if [[ $? == 42 ]]; then
	echo "update now..."
	DK_UPGRADE=1 bash -c "$(curl -L https://static.guance.com/datakit/install.sh)"
fi
```

### 添加 crontab 任务

执行如下命令，进入 crontab 规则添加界面：

```shell
crontab -u root -e
```

添加如下规则：

```shell
# 意即每天凌晨尝试一下新版本更新
0 0 * * * bash /path/to/datakit-update.sh
```

Tips: crontab 基本语法如下

```
*   *   *   *   *     <command to be execute>
^   ^   ^   ^   ^
|   |   |   |   |
|   |   |   |   +----- day of week(0 - 6) (Sunday=0)
|   |   |   +--------- month (1 - 12)   
|   |   +------------- day of month (1 - 31)
|   +----------------- hour (0 - 23)   
+--------------------- minute (0 - 59)
```

执行如下命令确保 crontab 安装成功：

```shell
crontab -u root -l
```

确保 crontab 服务启动：

```shell
service cron restart
```

如果安装成功且有尝试更新，则在 `update_log` 中能看到类似如下日志：

```
2021-05-10T09:49:06.083+0800 DEBUG	ota-update datakit/main.go:201	get online version...
2021-05-10T09:49:07.728+0800 DEBUG	ota-update datakit/main.go:216	online version: datakit 1.1.6-rc0/9bc4b960, local version: datakit 1.1.6-rc0-62-g7a1d0956/7a1d0956
2021-05-10T09:49:07.728+0800 INFO	ota-update datakit/main.go:224	Up to date(1.1.6-rc0-62-g7a1d0956)
```

如果确实发生了更新，会看到类似如下的更新日志：

```
2021-05-10T09:52:18.352+0800 DEBUG ota-update datakit/main.go:201 get online version...
2021-05-10T09:52:18.391+0800 DEBUG ota-update datakit/main.go:216 online version: datakit 1.1.6-rc0/9bc4b960, local version: datakit 1.0.1/7a1d0956
2021-05-10T09:52:18.391+0800 INFO  ota-update datakit/main.go:219 New online version available: 1.1.6-rc0, commit 9bc4b960 (release at 2021-04-30 14:31:27)
...
``` 
