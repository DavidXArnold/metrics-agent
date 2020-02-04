package kubernetes

import "time"

// NodeLabels contains labels for the k8s node
type NodeLabels struct {
	HostID   string `json:"host_id"`
	Hostname string `json:"hostname"`
	Labels   string `json:"labels"`
	Nodename string `json:"nodename"`
}

type node struct {
	Name      string `json:"nodeName"`
	StartTime string `json:"startTime"`
}

// NodeMetrics contains metrics for a k8s node
type NodeMetrics struct {
	CPULimit         []CPULimit        `json:"cpu/limit"`
	CPUUsage         []CPUUsage        `json:"cpu/usage"`
	FilesystemLimit  []FileSystemLimit `json:"filesystem/limit"`
	FilesystemUsage  []FileSystemUsage `json:"filesystem/usage"`
	MemoryCache      []MemoryCache     `json:"memory/cache"`
	MemoryLimit      []MemoryLimit     `json:"memory/limit"`
	MemoryRss        []MemoryRss       `json:"memory/rss"`
	MemoryUsage      []MemoryUsage     `json:"memory/usage"`
	NetworkRx        []NetworkRx       `json:"network/rx"`
	NetworkTx        []NetworkTx       `json:"network/tx"`
	Uptime           []Uptime          `json:"uptime"`
	FilesystemInodes []struct {
		End    time.Time `json:"end"`
		Labels struct {
			ResourceID string `json:"resource_id"`
		} `json:"labels"`
		Start time.Time `json:"start"`
		Value int       `json:"value"`
	} `json:"filesystem/inodes"`
	FilesystemInodesFree []struct {
		End    time.Time `json:"end"`
		Labels struct {
			ResourceID string `json:"resource_id"`
		} `json:"labels"`
		Start time.Time `json:"start"`
		Value int       `json:"value"`
	} `json:"filesystem/inodes_free"`
}

// NodeHeapsterCloneMetric replicates the heapster metrics object
// for a node
type NodeHeapsterCloneMetric struct {
	Labels  NodeLabels  `json:"labels"`
	Metrics NodeMetrics `json:"metrics"`
}

// Name returns the name of the node
func (n NodeHeapsterCloneMetric) Name() string {
	return n.Labels.Nodename
}

// NetworkRx shows the received network traffic
// for a given time window
type NetworkRx struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value uint64    `json:"value"`
}

// NetworkTx records the transmitted network traffic
// for a given time window
type NetworkTx struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value uint64    `json:"value"`
}

type FileSystemLabel struct {
	ResourceID string `json:"resource_id"`
}

type FileSystemUsage struct {
	End    time.Time       `json:"end"`
	Labels FileSystemLabel `json:"labels"`
	Start  time.Time       `json:"start"`
	Value  int             `json:"value"`
}

type FileSystemLimit struct {
	End    time.Time       `json:"end"`
	Labels FileSystemLabel `json:"labels"`
	Start  time.Time       `json:"start"`
	Value  int64           `json:"value"`
}

type Uptime struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value int       `json:"value"`
}

type CPUUsage struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value int64     `json:"value"`
}

type CPULimit struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value string    `json:"value"`
}
type MemoryCache struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value int       `json:"value"`
}
type MemoryLimit struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value string    `json:"value"`
}
type MemoryUsage struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value uint64    `json:"value"`
}

type MemoryRss struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Value int       `json:"value"`
}

type metricSample struct {
	Pods []pod
	Node node
}
type volume struct {
	Name string
}

// k8sEntity is a metrics type that encompasses nodes and pods
type k8sEntity interface {
	Name() string
}

// PodHeapsterCloneMetric replicates the heapster metrics object
// for a pod
type PodHeapsterCloneMetric struct {
	Labels  PodLabels  `json:"labels"`
	Metrics PodMetrics `json:"metrics"`
}

// Name returns the name of the pod
func (p PodHeapsterCloneMetric) Name() string {
	return p.Labels.PodName
}

// PodLabels is the labels on a k8s pod for a given container
type PodLabels struct {
	ContainerName string            `json:"container_name"`
	HostID        string            `json:"host_id"`
	Hostname      string            `json:"hostname"`
	Labels        map[string]string `json:"labels"`
	NamespaceID   string            `json:"namespace_id"`
	Nodename      string            `json:"nodename"`
	PodID         string            `json:"pod_id"`
	PodName       string            `json:"pod_name"`
}

// PodMetrics describes the metrics of a container within a pod
type PodMetrics struct {
	CPULimit        []CPULimit        `json:"cpu/limit"`
	CPUUsage        []CPUUsage        `json:"cpu/usage"`
	FilesystemLimit []FileSystemLimit `json:"filesystem/limit"`
	FilesystemUsage []FileSystemUsage `json:"filesystem/usage"`
	MemoryCache     []MemoryCache     `json:"memory/cache"`
	MemoryLimit     []MemoryLimit     `json:"memory/limit"`
	MemoryRss       []MemoryRss       `json:"memory/rss"`
	MemoryUsage     []MemoryUsage     `json:"memory/usage"`
	NetworkRx       []NetworkRx       `json:"network/rx"`
	NetworkTx       []NetworkTx       `json:"network/tx"`
	Uptime          []Uptime          `json:"uptime"`
}

type pod struct {
	HostID      string            `json:"host_id"`
	Hostname    string            `json:"hostname"`
	Labels      map[string]string `json:"labels"`
	NamespaceID string            `json:"namespace_id"`
	Namespace   string            `json:"namespace"`
	Nodename    string            `json:"nodename"`
	PodID       string            `json:"pod_id"`
	PodName     string            `json:"pod_name"`
	UID         string
	StartTime   string
	ClusterName string
	Containers  map[string]container
	Volumes     []volume
}

type container struct {
	Name         string
	StartTime    string
	Network      containerNetwork
	CPULimits    string
	MemoryLimits string
	ID           string
	RestartCount int32
}

type containerNetwork struct {
	NetworkRx NetworkRx
	NetworkTx NetworkTx
}

type dockerContainerMetrics struct {
	dockerContainers []container
}
