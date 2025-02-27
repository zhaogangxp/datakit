package zipkin

import (
	"strconv"

	"github.com/openzipkin/zipkin-go/model"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/internal/trace"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs/zipkin/corev1"
)

type zipkinThriftV1SpansFilter func([]*corev1.Span) []*corev1.Span

func zpkThriftV1Sample(zspans []*corev1.Span) []*corev1.Span {
	var rootSpan *corev1.Span
	for _, span := range zspans {
		if span.GetParentID() == 0 {
			rootSpan = span
		}
	}
	if rootSpan == nil {
		log.Debugf("can not find root span in zipkin thrift v1 spans %v", zspans)

		return zspans
	}

	sconf := trace.SampleConfMatcher(sampleConfs, convertZipkinBinaryAnnotationToTags(rootSpan.BinaryAnnotations))
	if sconf == nil {
		log.Debug("can not find sample config")

		return zspans
	}
	log.Debug(*sconf)

	for _, span := range zspans {
		tags := convertZipkinBinaryAnnotationToTags(span.BinaryAnnotations)
		// bypass error
		if _, ok := tags["error"]; ok {
			return zspans
		}
		// bypass ignore keys
		if trace.SampleIgnoreKeys(tags, sconf.IgnoreTagsList) {
			return zspans
		}
	}
	// do sample
	if trace.DefSampleFunc(uint64(rootSpan.TraceID), sconf.Rate, sconf.Scope) {
		return zspans
	} else {
		return nil
	}
}

func convertZipkinBinaryAnnotationToTags(annos []*corev1.BinaryAnnotation) map[string]string {
	tags := map[string]string{}
	for _, anno := range annos {
		tags[anno.Key] = string(anno.Value)
	}

	return tags
}

type zipkinJSONV1SpansFilter func([]*ZipkinSpanV1) []*ZipkinSpanV1

func zpkJSONV1Sample(zspans []*ZipkinSpanV1) []*ZipkinSpanV1 {
	var rootSpan *ZipkinSpanV1
	for _, span := range zspans {
		if span.ParentID == "" {
			rootSpan = span
		}
	}
	if rootSpan == nil {
		log.Debugf("can not find root span in zipkin json v1 spans %v", zspans)

		return zspans
	}

	sconf := trace.SampleConfMatcher(sampleConfs, convertBinaryAnnotationToTags(rootSpan.BinaryAnnotations))
	if sconf == nil {
		log.Debug("can not find sample config")

		return zspans
	}

	for _, span := range zspans {
		tags := convertBinaryAnnotationToTags(span.BinaryAnnotations)
		// bypass error
		if _, ok := tags["error"]; ok {
			return zspans
		}
		// bypass ignore keys
		if trace.SampleIgnoreKeys(tags, sconf.IgnoreTagsList) {
			return zspans
		}
	}
	// do sample
	traceid, err := strconv.ParseInt(rootSpan.TraceID, 10, 64)
	if err != nil {
		log.Debug(err)

		return zspans
	}
	if trace.DefSampleFunc(uint64(traceid), sconf.Rate, sconf.Scope) {
		return zspans
	} else {
		return nil
	}
}

func convertBinaryAnnotationToTags(annos []*BinaryAnnotation) map[string]string {
	tags := map[string]string{}
	for _, anno := range annos {
		tags[anno.Key] = anno.Value
	}

	return tags
}

type zipkinProtoBufV2SpansFilter func([]*model.SpanModel) []*model.SpanModel

func zpkProtoBufV2Sample(zspans []*model.SpanModel) []*model.SpanModel {
	var rootSpan *model.SpanModel
	for _, span := range zspans {
		if span.ParentID == nil {
			rootSpan = span
		}
	}
	if rootSpan == nil {
		log.Debugf("can not find root span in zipkin protobuf v2 spans %v", zspans)

		return zspans
	}

	sconf := trace.SampleConfMatcher(sampleConfs, rootSpan.Tags)
	if sconf == nil {
		log.Debug("can not find sample config")

		return zspans
	}

	for _, span := range zspans {
		// bypass error
		if span.Err != nil {
			return zspans
		}
		// bypass ignore keys
		if trace.SampleIgnoreKeys(span.Tags, sconf.IgnoreTagsList) {
			return zspans
		}
	}
	// do sample
	if trace.DefSampleFunc(uint64(trace.GetTraceID(int64(rootSpan.TraceID.High),
		int64(rootSpan.TraceID.Low))),
		sconf.Rate, sconf.Scope) {
		return zspans
	} else {
		return nil
	}
}

type zipkinJSONV2SpansFilter func([]*model.SpanModel) []*model.SpanModel

func zpkJSONV2Sample(zspans []*model.SpanModel) []*model.SpanModel {
	var rootSpan *model.SpanModel
	for _, span := range zspans {
		if span.ParentID == nil {
			rootSpan = span
		}
	}
	if rootSpan == nil {
		log.Debugf("can not find root span in zipkin json v2 spans %v", zspans)

		return zspans
	}

	sconf := trace.SampleConfMatcher(sampleConfs, rootSpan.Tags)
	if sconf == nil {
		log.Debug("can not find sample config")

		return zspans
	}

	for _, span := range zspans {
		// bypass error
		if span.Err != nil {
			return zspans
		}
		// bypass ignore keys
		if trace.SampleIgnoreKeys(span.Tags, sconf.IgnoreTagsList) {
			return zspans
		}
	}
	// do sample
	if trace.DefSampleFunc(uint64(trace.GetTraceID(int64(rootSpan.TraceID.High),
		int64(rootSpan.TraceID.Low))),
		sconf.Rate, sconf.Scope) {
		return zspans
	} else {
		return nil
	}
}
