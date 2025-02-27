package kubernetes

import (
	"context"
	"fmt"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
	kubev1batchbeta1 "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
)

const kubernetesCronJobName = "kubernetes_cron_jobs"

type cronJob struct {
	client interface {
		getCronJobs() kubev1batchbeta1.CronJobInterface
	}
	tags map[string]string
}

func (c *cronJob) Gather() {
	start := time.Now()
	var pts []*io.Point

	list, err := c.client.getCronJobs().List(context.Background(), metav1ListOption)
	if err != nil {
		l.Errorf("failed of get cronjobs resource: %s", err)
		return
	}

	for _, obj := range list.Items {
		tags := map[string]string{
			"name":          fmt.Sprintf("%v", obj.UID),
			"cron_job_name": obj.Name,
		}
		if obj.ClusterName != "" {
			tags["cluster_name"] = obj.ClusterName
		}
		if obj.Namespace != "" {
			tags["namespace"] = obj.Namespace
		}
		for k, v := range c.tags {
			tags[k] = v
		}

		fields := map[string]interface{}{
			"age":         int64(time.Since(obj.CreationTimestamp.Time).Seconds()),
			"schedule":    obj.Spec.Schedule,
			"active_jobs": len(obj.Status.Active),
		}

		if obj.Spec.Suspend != nil {
			fields["suspend"] = *obj.Spec.Suspend
		} else {
			fields["suspend"] = defaultBoolerValue
		}

		addMapToFields("annotations", obj.Annotations, fields)
		addLabelToFields(obj.Labels, fields)
		addMessageToFields(tags, fields)

		pt, err := io.MakePoint(kubernetesCronJobName, tags, fields, time.Now())
		if err != nil {
			l.Error(err)
		} else {
			pts = append(pts, pt)
		}
	}

	if len(pts) == 0 {
		l.Debug("no points")
		return
	}

	if err := io.Feed(inputName, datakit.Object, pts, &io.Option{CollectCost: time.Since(start)}); err != nil {
		l.Error(err)
	}
}

func (*cronJob) Resource() { /*empty interface*/ }

func (*cronJob) LineProto() (*io.Point, error) { return nil, nil }

//nolint:lll
func (*cronJob) Info() *inputs.MeasurementInfo {
	return &inputs.MeasurementInfo{
		Name: kubernetesCronJobName,
		Desc: "Kubernetes cron job 对象数据",
		Type: "object",
		Tags: map[string]interface{}{
			"name":          inputs.NewTagInfo("UID"),
			"cron_job_name": inputs.NewTagInfo("Name must be unique within a namespace."),
			"cluster_name":  inputs.NewTagInfo("The name of the cluster which the object belongs to."),
			"namespace":     inputs.NewTagInfo("Namespace defines the space within each name must be unique."),
		},
		Fields: map[string]interface{}{
			"age":         &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.DurationSecond, Desc: "age (seconds)"},
			"schedule":    &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron"},
			"active_jobs": &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.NCount, Desc: "The number of pointers to currently running jobs."},
			"suspend":     &inputs.FieldInfo{DataType: inputs.Bool, Unit: inputs.UnknownUnit, Desc: "This flag tells the controller to suspend subsequent executions, it does not apply to already started executions."},
			"annotations": &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "kubernetes annotations"},
			"message":     &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "object details"},
		},
	}
}
