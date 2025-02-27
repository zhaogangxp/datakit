{{.CSS}}

- DataKit 版本：{{.Version}}
- 文档发布日期：{{.ReleaseDate}}
- 操作系统支持：`{{.AvailableArchs}}`

# {{.InputName}}

该采集器是网络拨测结果数据采集，所有拨测产生的数据，上报观测云。

## 私有拨测节点部署

私有拨测节点部署，需在 [观测云页面创建私有拨测节点](https://www.yuque.com/dataflux/doc/phmtep)。创建完成后，将页面上相关信息填入 `conf.d/{{.Catalog}}/{{.InputName}}.conf` 即可：

```toml
#  中心任务存储的服务地址
server = "https://dflux-dial.guance.com"

# require，节点惟一标识ID
region_id = "reg_c2jlokxxxxxxxxxxx"

# 若server配为中心任务服务地址时，需要配置相应的ak或者sk
ak = "ZYxxxxxxxxxxxx"
sk = "BNFxxxxxxxxxxxxxxxxxxxxxxxxxxx"

[inputs.dialtesting.tags]
  # some_tag = "some_value"
  # more_tag = "some_other_value"
  # ...
```

## 拨测部署图

![](https://zhuyun-static-files-production.oss-cn-hangzhou.aliyuncs.com/images/datakit/dialtesting-net-arch.png)

<!--
## 浏览器拨测（Headless 拨测）

浏览器测需在 DataKit 上安装 Chrome 浏览器，以 Ubuntu 为例：

```shell
sudo apt-get install libxss1 libappindicator1 libindicator7
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i google-chrome-stable_current_amd64.deb
sudo apt-get install -f
```

浏览器拨测无需修改 DataKit 配置，只需在观测云中[设置浏览器拨测任务](https://www.yuque.com/dataflux/doc/qnfc4a#UkJNb)即可。
-->

## 配置

进入 DataKit 安装目录下的 `conf.d/{{.Catalog}}` 目录，复制 `{{.InputName}}.conf.sample` 并命名为 `{{.InputName}}.conf`。示例如下：

```toml
{{.InputSample}}
```

配置好后，重启 DataKit 即可。

## 指标集

以下所有数据采集，默认会追加名为 `host` 的全局 tag（tag 值为 DataKit 所在主机名），也可以在配置中通过 `[[inputs.{{.InputName}}.tags]]` 另择 host 来命名。

{{ range $i, $m := .Measurements }}

### `{{$m.Name}}`

-  标签

{{$m.TagsMarkdownTable}}

- 指标列表

{{$m.FieldsMarkdownTable}}

{{ end }}
