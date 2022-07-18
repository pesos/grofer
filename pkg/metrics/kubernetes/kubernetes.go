package kubernetes_metrics

// TODO import client-go
import (
	"github.com/pesos/grofer/pkg/core"
)

// TODO define k8s-specific metrics structs
/*
type containerMetrics struct {
	client      *client.Client
	all         bool
	refreshRate uint64
	sink        core.Sink // defaults to TUI.
	metricBus   chan container.OverallMetrics
}
*/

type Pod struct {
	name string
}

type KubernetesMetricsScraper struct {
	// TODO reference to initialized client-go goes here.
	Clientset   string
	RefreshRate uint64
	Sink        core.Sink
	MetricBus   chan KubernetesMetrics
}

// TODO add more data here!
type KubernetesMetrics struct {
	pods []string
}
