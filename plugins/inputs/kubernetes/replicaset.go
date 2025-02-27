package kubernetes

import (
	"context"
	"fmt"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
	kubev1apps "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const kubernetesReplicaSetName = "kubernetes_replica_sets"

type replicaSet struct {
	client interface {
		getReplicaSets() kubev1apps.ReplicaSetInterface
	}
	tags map[string]string
}

func (r *replicaSet) Gather() {
	start := time.Now()
	var pts []*io.Point

	list, err := r.client.getReplicaSets().List(context.Background(), metav1ListOption)
	if err != nil {
		l.Errorf("failed of get replicaSet resource: %s", err)
		return
	}

	for _, obj := range list.Items {
		tags := map[string]string{
			"name":             fmt.Sprintf("%v", obj.UID),
			"replica_set_name": obj.Name,
		}
		if obj.ClusterName != "" {
			tags["cluster_name"] = obj.ClusterName
		}
		if obj.Namespace != "" {
			tags["namespace"] = obj.Namespace
		}
		for k, v := range r.tags {
			tags[k] = v
		}

		fields := map[string]interface{}{
			"age":       int64(time.Since(obj.CreationTimestamp.Time).Seconds()),
			"ready":     obj.Status.ReadyReplicas,
			"available": obj.Status.AvailableReplicas,
		}

		// addMapToFields("selectors", obj.Spec.Selector, fields)
		addMapToFields("annotations", obj.Annotations, fields)
		addLabelToFields(obj.Labels, fields)
		addMessageToFields(tags, fields)

		pt, err := io.MakePoint(kubernetesReplicaSetName, tags, fields, time.Now())
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

//nolint:unused
func (*replicaSet) resource() { /*empty interface*/ }

func (*replicaSet) LineProto() (*io.Point, error) { return nil, nil }

//nolint:lll
func (*replicaSet) Info() *inputs.MeasurementInfo {
	return &inputs.MeasurementInfo{
		Name: kubernetesReplicaSetName,
		Desc: "Kubernetes replicaSet 对象数据",
		Type: "object",
		Tags: map[string]interface{}{
			"name":             inputs.NewTagInfo("UID"),
			"replica_set_name": inputs.NewTagInfo("Name must be unique within a namespace."),
			"cluster_name":     inputs.NewTagInfo("The name of the cluster which the object belongs to."),
			"namespace":        inputs.NewTagInfo("Namespace defines the space within each name must be unique."),
		},
		Fields: map[string]interface{}{
			"age":         &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.DurationSecond, Desc: "age (seconds)"},
			"ready":       &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.UnknownUnit, Desc: "The number of ready replicas for this replica set."},
			"available":   &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.UnknownUnit, Desc: "The number of available replicas (ready for at least minReadySeconds) for this replica set."},
			"annotations": &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "kubernetes annotations"},
			"message":     &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "object details"},
			// TODO:
			// "selectors": &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
			// "current/desired":        &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
		},
	}
}
