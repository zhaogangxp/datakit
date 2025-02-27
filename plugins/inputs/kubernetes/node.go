package kubernetes

import (
	"context"
	"fmt"
	"time"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/io"
	"gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs"
	kubev1core "k8s.io/client-go/kubernetes/typed/core/v1"
)

const kubernetesNodeName = "kubernetes_nodes"

type node struct {
	client interface {
		getNodes() kubev1core.NodeInterface
	}
	tags map[string]string
}

func (n *node) Gather() {
	start := time.Now()
	var pts []*io.Point

	list, err := n.client.getNodes().List(context.Background(), metav1ListOption)
	if err != nil {
		l.Errorf("failed of get nodes resource: %s", err)
		return
	}

	for _, obj := range list.Items {
		tags := map[string]string{
			"name":      fmt.Sprintf("%v", obj.UID),
			"node_name": obj.Name,
			"status":    fmt.Sprintf("%v", obj.Status.Phase),
		}
		if obj.ClusterName != "" {
			tags["cluster_name"] = obj.ClusterName
		}
		if obj.Namespace != "" {
			tags["namespace"] = obj.Namespace
		}
		if ip := datakit.GetEnv("HOST_IP"); ip != "" {
			tags["node_ip"] = ip
		}
		for k, v := range n.tags {
			tags[k] = v
		}

		fields := map[string]interface{}{
			"age":             int64(time.Since(obj.CreationTimestamp.Time).Seconds()),
			"kubelet_version": obj.Status.NodeInfo.KubeletVersion,
		}

		addMapToFields("annotations", obj.Annotations, fields)
		addLabelToFields(obj.Labels, fields)
		addMessageToFields(tags, fields)

		pt, err := io.MakePoint(kubernetesNodeName, tags, fields, time.Now())
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

func (*node) Resource() { /*empty interface*/ }

func (*node) LineProto() (*io.Point, error) { return nil, nil }

//nolint:lll
func (*node) Info() *inputs.MeasurementInfo {
	return &inputs.MeasurementInfo{
		Name: kubernetesNodeName,
		Desc: "Kubernetes node 对象数据",
		Type: "object",
		Tags: map[string]interface{}{
			"name":         inputs.NewTagInfo("UID"),
			"node_name":    inputs.NewTagInfo("Name must be unique within a namespace."),
			"node_ip":      inputs.NewTagInfo("Node IP"),
			"cluster_name": inputs.NewTagInfo("The name of the cluster which the object belongs to."),
			"namespace":    inputs.NewTagInfo("Namespace defines the space within each name must be unique."),
			"status":       inputs.NewTagInfo("NodePhase is the recently observed lifecycle phase of the node. (Pending/Running/Terminated)"),
		},
		Fields: map[string]interface{}{
			"age":             &inputs.FieldInfo{DataType: inputs.Int, Unit: inputs.DurationSecond, Desc: "age (seconds)"},
			"kubelet_version": &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "Kubelet Version reported by the node."},
			"annotations":     &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "kubernetes annotations"},
			"message":         &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: "object details"},
			// TODO:
			// "schedulability":  &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
			// "role":            &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
			// "taints":                 &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
			// "pods":                   &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
			// "pod_capacity":           &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
			// "pod_usage":              &inputs.FieldInfo{DataType: inputs.String, Unit: inputs.UnknownUnit, Desc: ""},
		},
	}
}
