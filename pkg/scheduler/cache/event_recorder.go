package cache

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"volcano.sh/volcano/pkg/scheduler/metrics"
)

type NoopEventRecorderWrapper struct {
}

func NewEventRecorderWrapper(config *rest.Config, schedulerNames []string) *NoopEventRecorderWrapper {
	return &NoopEventRecorderWrapper{}

}

func (r *NoopEventRecorderWrapper) Event(object runtime.Object, eventtype, reason, message string) {
	metrics.IncEventCount()
	//r.recorder.Event(object, eventtype, reason, message)
}

func (r *NoopEventRecorderWrapper) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	metrics.IncEventCount()
	//r.recorder.Eventf(object, eventtype, reason, messageFmt, args...)
}

func (r *NoopEventRecorderWrapper) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	metrics.IncEventCount()
	//r.recorder.AnnotatedEventf(object, annotations, eventtype, reason, messageFmt, args...)
}
