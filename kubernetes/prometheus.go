package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cloudability/metrics-agent/util"
	"github.com/kubernetes/kubernetes/staging/src/k8s.io/client-go/util/retry"
	"github.com/prometheus/client_golang/api"
	pv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

//NewPromClient creates a new prometheus client with a given URL eg: http://demo.robustperception.io:9090"
func NewPromClient(url string, rt http.RoundTripper) (client api.Client, err error) {
	client, err = api.NewClient(api.Config{
		Address:      url,
		RoundTripper: rt,
	})
	if err != nil {
		log.Errorf("Error creating prometheus client: %v\n", err)
		return
	}

	return
}

//APIQuery queries a prometheus instance with a given client and free form query
func APIQuery(c *api.Client, q string) (result model.Value, err error) {

	v1api := pv1.NewAPI(*c)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, q, time.Now())
	if err != nil {
		log.Errorf("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		log.Warnf("Prometheus query warnings: %v\n", warnings)
	}

	log.Debugf("Result:\n%v\n", result)

	return result, err
}

//APIQueryRange queries a prometheus instance with a given client and range query
func APIQueryRange(c *api.Client, q string, r pv1.Range) (result model.Value, err error) {
	// func APIQueryRange(c *api.Client, q string, r pv1.Range) (result model.Value, err error) {

	v1api := pv1.NewAPI(*c)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.QueryRange(ctx, q, r)
	if err != nil {
		log.Errorf("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		log.Warnf("Prometheus query warnings: %v\n", warnings)
	}

	// log.Debugf("Result:\n%v\n", result)

	return result, err
}

//APIQuerySeries returns prometheus series with a given client, label matchers, and time range
func APIQuerySeries(c *api.Client, m []string, startTime time.Time, endTime time.Time) (result []model.LabelSet, err error) {

	v1api := pv1.NewAPI(*c)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// lbls, warnings, err := v1api.Series(ctx, []string{
	// 	"{__name__=~\"scrape_.+\",job=\"node\"}",
	// 	"{__name__=~\"scrape_.+\",job=\"prometheus\"}",
	// }, time.Now().Add(-time.Hour), time.Now())
	lbls, warnings, err := v1api.Series(ctx, m, startTime, endTime)
	if err != nil {
		log.Errorf("Error querying Prometheus: %v\n", err)
		return nil, err
	}
	if len(warnings) > 0 {
		log.Warnf("Prometheus query warnings: %v\n", warnings)
	}
	// log.Debugf("Result:\n%v\n", result)
	// fmt.Println("Result:")
	// for _, lbl := range lbls {
	// 	fmt.Println(lbl)
	// }
	return lbls, err
}

//DownloadPromData downloads prometheus data for a given range of hosts / pods etc..
func DownloadPromData(config KubeAgentConfig, msd string, metricSampleDir *os.File, nodeSource NodeSource) (err error) {

	// c, err := api.NewClient(api.Config{
	// 	Address:      config.PrometheusURL,
	// 	RoundTripper: config.HTTPClient.Transport,
	// })

	c, err := NewPromClient(config.PrometheusURL, config.HTTPClient.Transport)
	if err != nil {
		log.Errorf("Error creating Prometheus client: %v\n", err)
		time.Sleep(time.Minute * 2)
		return err
	}

	var nodes []v1.Node

	// failedNodeList := make(map[string]error)

	err = retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		nodes, err = nodeSource.GetReadyNodes()
		return
	})
	if err != nil {
		return fmt.Errorf("cloudability metric agent is unable to get a list of nodes: %v", err)
	}

	r := pv1.Range{
		Start: time.Now().Add(-time.Duration(3) * time.Minute),
		End:   time.Now(),
		Step:  time.Minute,
	}

	for _, n := range nodes {
		q := "{__name__=~\"kube_node_info|kube_node_status_capacity_memory_bytes|kube_node_status_capacity_cpu_cores|container_cpu_usage_seconds_total|container_memory_rss|container_spec_memory_limit_bytes|container_network_receive_bytes_total|container_network_transmit_bytes_total|container_last_seen\"}"
		if n.Spec.ProviderID == "" {
			log.Error("Provider ID for node does not exist. " +
				"If this condition persists it will cause inconsistent cluster allocation")
		}
		// lm := "{__name__=~\"kube_node_info|kube_node_status_capacity_memory_bytes|kube_node_status_capacity_cpu_cores|container_cpu_usage_seconds_total|container_memory_rss|container_spec_memory_limit_bytes|container_network_receive_bytes_total|container_network_transmit_bytes_total|container_last_seen\"}"
		// q := lm + " and {kubernetes_io_hostname=\"" + n.Name + "\"} or " + lm + " and {instance=\"" + n.Name + "\"}"

		result, err := APIQueryRange(&c, q, r)
		if err != nil {
			log.Errorf("Prometheus query error: %v\n", err)
			return err
		}
		// for _, m := range result.Type().MarshalJSON() {
		// 	log.Infof("host: %v \nresult: %+v", n.Name, m)
		// }
		// r, err := result.Type().MarshalJSON()

		log.Infof("host: %v type: %+v result: %+v", n.Name, result.Type().String, result.String)
		log.Info("Break:\n")
	}

	// q := "{__name__=~\"kube_node_info|kube_node_status_capacity_memory_bytes|kube_node_status_capacity_cpu_cores|container_cpu_usage_seconds_total|container_memory_rss|container_spec_memory_limit_bytes|container_network_receive_bytes_total|container_network_transmit_bytes_total|container_last_seen\"}"

	// result, err := APIQueryRange(&c, q, r)

	// log.Infof("result %+v", result)

	// log.Fatalf("out now bro")
	// os.Exit(2)
	return err
}

func validatePrometheus(config KubeAgentConfig, client rest.HTTPClient) error {

	test, _, err := util.TestHTTPConnection(
		client, config.PrometheusURL, http.MethodGet, "", retryCount, true)
	if err != nil {
		return err
	}
	if !test {
		return fmt.Errorf("Unable to connect to Prometheus URL: %v %v", config.PrometheusURL, err)
	}
	log.Debugf("Connected to Prometheus at: %v", config.PrometheusURL)

	return err
}
