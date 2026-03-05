package cache

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"strings"
	schedulingv1beta1 "volcano.sh/apis/pkg/apis/scheduling/v1beta1"
)

func ShouldAcceptPod(podObj interface{}, schedulerSelector map[string]bool) bool {
	switch v := podObj.(type) {
	case *v1.Pod:
		return isPodSchedulerAtSelector(v, schedulerSelector)
	case cache.DeletedFinalStateUnknown:
		if p, ok := v.Obj.(*v1.Pod); !ok {
			klog.Errorf("Cannot convert object %T to *v1.Pod", v.Obj)
			return false
		} else {
			return isPodSchedulerAtSelector(p, schedulerSelector)
		}
	default:
		return false
	}
}

func isPodSchedulerAtSelector(pod *v1.Pod, schedulerSelector map[string]bool) bool {
	if _, ok := schedulerSelector[pod.Spec.SchedulerName]; ok {
		return true
	}
	// 非volcano的系统pod也要被处理，否则node资源会计算不准
	if !strings.Contains(pod.Spec.SchedulerName, "volcano") {
		return true
	}
	return false
}

func ShouldAcceptQueue(queueObj interface{}, queueSelector map[string]bool) bool {
	switch v := queueObj.(type) {
	case *schedulingv1beta1.Queue:
		return isQueueAtScheduler(v, queueSelector)
	case cache.DeletedFinalStateUnknown:
		if q, ok := v.Obj.(*schedulingv1beta1.Queue); !ok {
			return false
		} else {
			return isQueueAtScheduler(q, queueSelector)
		}
	default:
		klog.V(3).Infof("Unknow queue type %v", v)
		return false
	}
}

func isQueueAtScheduler(queue *schedulingv1beta1.Queue, queueSelector map[string]bool) bool {
	if _, ok := queueSelector[queue.Name]; ok {
		return true
	}
	return false
}

func ShouldAcceptNode(nodeObj interface{}, nodeSelectors map[string]sets.Empty) bool {
	switch t := nodeObj.(type) {
	case *v1.Node:
		if isNodeAtSelector(t, nodeSelectors) {
			klog.V(3).Infof("Node <%s> will be accept by current scheduler.", t.Name)
			return true
		}
		klog.V(3).Infof("Node <%s> should not be accept by current scheduler.", t.Name)
		return false
	case cache.DeletedFinalStateUnknown:
		n, ok := t.Obj.(*v1.Node)
		if !ok {
			klog.Errorf("Cannot convert to *v1.Node: %v", t.Obj)
			return false
		}

		if isNodeAtSelector(n, nodeSelectors) {
			klog.V(3).Infof("Node <%s> will be accept by current scheduler.", n.Name)
			return true
		}
		klog.V(3).Infof("Node <%s> should not be accept by current scheduler.", n.Name)
		return false
	default:
		klog.V(3).Infof("Node <%s> should not be accept by current scheduler.", nodeObj)
		return false
	}
}

func isNodeAtSelector(node *v1.Node, nodeSelectors map[string]sets.Empty) bool {
	if len(nodeSelectors) == 0 {
		return false
	}
	for labelName, labelValue := range node.Labels {
		key := labelName + ":" + labelValue
		if _, ok := nodeSelectors[key]; ok {
			return true
		}
	}
	return false
}
