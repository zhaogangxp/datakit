package demo

import (
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
)

type demoMetric2 struct {
	name   string
	tags   map[string]string
	fields map[string]interface{}
	ts     time.Time
}

func (m *demoMetric2) LineProto() (*io.Point, error) {
	return io.MakePoint(m.name, m.tags, m.fields, m.ts)
}

func (m *demoMetric2) Info() *inputs.MeasurementInfo {
	return &inputs.MeasurementInfo{
		Name: "demo-metric2",
		Type: "metric",
		Desc: "这是一个指标集的 demo(**务必加上每个指标集的描述**)",
		Tags: map[string]interface{}{
			"tag_a": &inputs.TagInfo{Desc: "示例 tag A"},
			"tag_b": &inputs.TagInfo{Desc: "示例 tag B"},
		},
		Fields: map[string]interface{}{
			"usage": &inputs.FieldInfo{
				DataType: inputs.Float,
				Type:     inputs.Gauge,
				Unit:     inputs.Percent,
				Desc:     "this is CPU usage",
			},
			"disk_size": &inputs.FieldInfo{
				DataType: inputs.Int,
				Type:     inputs.Gauge,
				Unit:     inputs.SizeByte,
				Desc:     "this is disk size",
			},
			"mem_size": &inputs.FieldInfo{
				DataType: inputs.Int,
				Type:     inputs.Gauge,
				Unit:     inputs.SizeIByte,
				Desc:     "this is memory size",
			},
			"some_string": &inputs.FieldInfo{
				DataType: inputs.String,
				Type:     inputs.Gauge,
				Unit:     inputs.UnknownUnit,
				Desc:     "some string field",
			},
			"ok": &inputs.FieldInfo{
				DataType: inputs.Bool,
				Type:     inputs.Gauge,
				Unit:     inputs.UnknownUnit,
				Desc:     "some boolean field",
			},
		},
	}
}
