{{.CSS}}

- DataKit 版本：{{.Version}}
- 文档发布日期：{{.ReleaseDate}}
- 操作系统支持：`{{.AvailableArchs}}`

# {{.InputName}}

采集 container 指标数据、对象数据和容器日志，以及当前主机上的 kubelet Pod 指标和对象，上报到观测云。

## 前置条件

- 目前 container 会默认连接 Docker 服务，需安装 Docker v17.04 及以上版本。

## 配置

进入 DataKit 安装目录下的 `conf.d/{{.Catalog}}` 目录，复制 `{{.InputName}}.conf.sample` 并命名为 `{{.InputName}}.conf`。示例如下：

```toml
{{.InputSample}} 
```

注意，默认不开启指标采集。如有需要，请将 `enable_metric` 改为 `true` 并重启 DataKit。

支持以环境变量的方式修改配置参数（只在 DataKit 以 K8s daemonset 方式运行时生效，主机部署的 DataKit 不支持此功能）：

| 环境变量名                                             | 对应的配置参数项                    | 参数示例                                                     |
| :---                                                   | ---                                 | ---                                                          |
| `ENV_INPUT_CONTAINER_ENABLE_METRIC`                    | `enable_metric`                     | `true`/`false`                                               |
| `ENV_INPUT_CONTAINER_ENABLE_OBJECT`                    | `enable_object`                     | `true`/`false`                                               |
| `ENV_INPUT_CONTAINER_ENABLE_LOGGING`                   | `enable_logging`                    | `true`/`false`                                               |
| `ENV_INPUT_CONTAINER_LOGGING_REMOVE_ANSI_ESCAPE_CODES` | `logging_remove_ansi_escape_codes ` | `true`/`false`                                               |
| `ENV_INPUT_CONTAINER_TAGS`                             | `tags`                              | `tag1=value1,tag2=value2` 如果配置文件中有同名 tag，会覆盖它 |

## 指标集

以下所有数据采集，默认会追加名为 `host` 的全局 tag（tag 值为 DataKit 所在主机名），也可以在配置中通过 `[inputs.{{.InputName}}.tags]` 指定其它标签：

``` toml
 [inputs.{{.InputName}}.tags]
  # some_tag = "some_value"
  # more_tag = "some_other_value"
  # ...
```

## 指标

{{ range $i, $m := .Measurements }}

{{if eq $m.Type "metric"}}

### `{{$m.Name}}`
{{$m.Desc}}

-  标签

{{$m.TagsMarkdownTable}}

- 字段列表

{{$m.FieldsMarkdownTable}}
{{end}}

{{ end }}

## 对象

{{ range $i, $m := .Measurements }}

{{if eq $m.Type "object"}}

### `{{$m.Name}}`
{{$m.Desc}}

-  标签

{{$m.TagsMarkdownTable}}

- 字段列表

{{$m.FieldsMarkdownTable}}
{{end}}

{{ end }}

## 日志

{{ range $i, $m := .Measurements }}

{{if eq $m.Type "logging"}}

### `{{$m.Name}}`
{{$m.Desc}}

-  标签

{{$m.TagsMarkdownTable}}

- 字段列表

{{$m.FieldsMarkdownTable}}
{{end}}

{{ end }}

### 标签定制和删除

- `drop_tags`：对于某些 Tag，收集的意义不大，但会导致时间线暴涨。目前将 `contaienr_id` 这个 tag 移除掉了。
- `ignore_image_name`：配置正则表达式列表[正则表达式参见这里](https://golang.org/pkg/regexp/syntax/#hdr-Syntax)，如果容器的 `image_name` 能够匹配任意其一，则忽略此条容器的数据
- `ignore_container_name`：配置正则表达式列表，如果容器的 `container_name` 能够匹配任意其一，则忽略此条容器的数据
- `kubelet.ignore_container_name`：配置正则表达式列表，如果 kubelet pod 的 `pod_name` 能够匹配任意其一，则忽略此条数据。适用于 pod 的指标和对象

### 容器日志采集

对于容器日志采集，有着更细致的配置，主要解决了区分日志 `source` 和使用 pipeline 的问题。

日志采集配置项为 `[[inputs.container.log]]`，该项是数组配置，意即可以有多个 log 来处理采集到的容器日志，比如某个容器中既有 MySQL 日志，也有 Redis 日志，那么此时可能需要两个 log 来分别处理它们。

- `match_by`：
   - 匹配类型，主要配置项，只支持填写 `contianer-name` 和 `deployment-name`（注意是中横线）。
   - 例如该项为 `container_name`，则会以容器名进行后续的正则匹配。当匹配成功时使用此 `log` 各项参数配置（`source`、`pipeline`、`ignore_status`、`character_encoding` 和 `multiline_match`）
- `match`：
   - 匹配日志来源的正则表达式，如果 `match_by` 是 `container_name`，这里需要添加容器名的正则，用以在所有容器中找到一个或多个指定的容器。
   - 该参数类型是字符串数组，只要任意一个正则匹配成功即可。
   - 未匹配的容器，其日志将执行默认处理方式。
   - [正则表达式参见这里](https://golang.org/pkg/regexp/syntax/#hdr-Syntax)

>Tips：为保证此处正则表达式的正确书写，请务必将正则表达式用 `'''这里是一个正则表达式'''` 这种形式来配置（即两边用三个单引号来包围正则文本），否则可能导致正则转义问题。
- `source`：
   - 指定数据来源，其值不可为空
- `service`：
   - 指定该条日志的服务名，如果为空值，则使用 `source` 字段值
- `pipeline`：
   - 指定 pipeline 文件，只需写文件名即可，不需要写全路径，使用方式见 [Pipeline 文档](pipeline)。
   - 当此值为空值或该文件不存在时，将不使用 pipeline 功能
- `ignore_status`：
   - 忽略对应 `status` 数据，只能是 `"emerg","alert","critical","error","warning","info","debug","OK"`
   - 例如填写 `["info"]`，在经过 pipeline 处理后，将丢弃所有 `status` 字段为 `info` 的数据
- `character_encoding`：
   - 选择字符编码，如果编码有误会导致数据无法查看。默认为空即可。
   - 支持的编码有 `"utf-8", "utf-16le", "utf-16le", "gbk", "gb18030" or ""`
- `multiline_match`：
   - 设置正则表达式用以配置多行，例如 `^\d{4}-\d{2}-\d{2}` 行首匹配 `YYYY-MM-DD` 时间格式，和 `logging` 采集器的同名配置字段用法相同

如果一个容器的 `container name` 和 `deployment` 分别匹配两个 log，会优先使用 `deployment` 所匹配的 log。例如容器的 `container name` 为 `containerAAA`，`deployment` 为 `deploymentAAA`，且配置如下：

```toml
[[inputs.container.log]]
  match_by = "container-name"
  match = ['''container*''']
  source = "dummy1"
  service = "dummy1"
  pipeline = "dummy1.p"
  ignore_status = []
  character_encoding = ""
  # multiline_match = '''^\S'''

[[inputs.container.log]]
  match_by = "deployment-name"
  match = ['''deployment*''']
  source = "dummy2"
  service = "dummy2"
  pipeline = "dummy2.p"
  ignore_status = []
  character_encoding = ""
  # multiline_match = '''^\S'''
```

此时该容器能够匹配两个 log，优先使用第二个 `deployment`。

### 日志切割注意事项

使用 pipeline 功能时，如果切割成功，则：

- 取其中的 `time` 字段作为此条数据的产生时间。如果没有 `time` 字段或解析此字段失败，默认使用当前 DataKit 所在机器的系统时间
- 所切割出来日志结果中，必定有一个 `status` 字段。如果切割出来的原始数据中没有该字段，则默认将 `status` 置为 `info`

当前有效的 `status` 字段值如下（三列均不区分大小写）：

| 简写 | 可能的全称                  | 对应值     |
| :--- | ---                         | -------    |
| `f`  | `emerg`                     | `emerg`    |
| `a`  | `alert`                     | `alert`    |
| `c`  | `critical`                  | `critical` |
| `e`  | `error`                     | `error`    |
| `w`  | `warning`                   | `warning`  |
| `n`  | `notice`                    | `notice`   |
| `i`  | `info`                      | `info`     |
| `d`  | `debug`, `trace`, `verbose` | `debug`    |
| `o`  | `s`, `OK`                   | `OK`       |

### 容器日志的特殊字节码过滤

容器日志可能会包含一些不可读的字节码（比如终端输出的颜色等），可以将 `logging_remove_ansi_escape_codes` 设置为 `true` 对其删除过滤。

此配置可能会影响日志的处理性能，基准测试结果如下：

```
goos: linux
goarch: amd64
pkg: gitlab.jiagouyun.com/cloudcare-tools/test
cpu: Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz
BenchmarkRemoveAnsiCodes
BenchmarkRemoveAnsiCodes-8        636033              1616 ns/op
PASS
ok      gitlab.jiagouyun.com/cloudcare-tools/test       1.056s
```

每一条文本的处理耗时增加 `1616 ns` 不等。如果不开启此功能将无额外损耗。

### kubelet 相关采集

在配置文件中打开 `inputs.container.kubelet` 项，填写对应的 `kubelet_url`（默认为 `127.0.0.1:10255`）可以采集 kubelet Pod 相关数据。

kubelet 该端口默认关闭，开启方式请查看[官方文档](https://kubernetes.io/zh/docs/reference/command-line-tools-reference/kubelet/)搜索 `--read-only-port`。

如果 `kubelet_url` 配置的主机端口未监听，或尝试连接失败，则不再采集 kubelet 相关数据。DataKit 只有在启动或重启时才会对该端口进行连接验证，一旦验证失败，直到下次重启前都不会再次连接 kubelet。如果 kubelet 经过修复已经开启对应端口，需重启 DataKit 才能采集。

在连接 kubelet 时，可能会因为 kubelet 认证问题报错：

- 报错一：`/run/secrets/kubernetes.io/serviceaccount/token: no such file or directory`

执行如下两个命令准备对应文件：

```shell
# mkdir -p /run/secrets/kubernetes.io/serviceaccount
# touch /run/secrets/kubernetes.io/serviceaccount/token
```

- 报错二： `error making HTTP request to http://<k8s-host>/stats/summary: dial tcp <k8s-hosst>:10255: connect: connect refused`

按如下步骤调整 k8s 配置：

  1. 编辑所有节点的 `/var/lib/kubelet/config.yaml` 文件，加入`readOnlyPort` 这个参数：`readOnlyPort: 10255`
  1. 重启 kubelet 服务：`systemctl restart kubelet.service`
