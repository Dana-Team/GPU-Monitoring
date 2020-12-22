package collector

import (
	"context"
	"fmt"
	"github.com/Dana-Team/GPU-Monitoring/types"
	"github.com/Dana-Team/gonvml"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

type GpuCollector struct {
	client.Client
	logger  logr.Logger
	Host string
	Devices []*gonvml.Device
}

type PodSample struct {
	pod *corev1.Pod
	putil *gonvml.ProcessUtilizationSample

}

func NewGpuCollector(client client.Client) (*GpuCollector, error) {
	devices, err := getDevices()
	if err != nil {
		return nil, err
	}
	gpuc := &GpuCollector{
		client,
		ctrl.Log.WithName("collector"),
		os.Getenv("NODE_NAME"),
		devices}
	prometheus.WrapRegistererWith(prometheus.Labels{}, metrics.Registry).MustRegister(gpuc)
	return gpuc, nil
}

func (c *GpuCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *GpuCollector) Collect(ch chan<- prometheus.Metric) {
	//get all pods running on the host
	var pods corev1.PodList
	c.List(
		context.Background(),
		&pods,
		client.InNamespace(""),
		client.MatchingFields{types.NodeKey: os.Getenv(types.NodeNameEnv)},)

	for _, device := range c.Devices {
		ds, err := getDeviceStatus(*device)
		if err != nil{
			c.logger.Error(err, "unable to get device: %s status",device.UUID)
			continue
		}
		//device metrics
		ch <- prometheus.MustNewConstMetric(
			gpuUtilDesc,
			prometheus.GaugeValue,
			float64(*ds.Utilization.GPU),
			c.Host,
			device.UUID,
		)
		ch <- prometheus.MustNewConstMetric(
			gpuMemUtilDesc,
			prometheus.GaugeValue,
			float64(*ds.Utilization.Memory),
			c.Host,
			device.UUID,
		)
		ch <- prometheus.MustNewConstMetric(
			gpuEncUtilDesc,
			prometheus.GaugeValue,
			float64(*ds.Utilization.Encoder),
			c.Host,
			device.UUID,
		)
		ch <- prometheus.MustNewConstMetric(
			gpuDecUtilDesc,
			prometheus.GaugeValue,
			float64(*ds.Utilization.Decoder),
			c.Host,
			device.UUID,
		)
		ch <- prometheus.MustNewConstMetric(
			gpuFreeMemDesc,
			prometheus.GaugeValue,
			float64(*ds.Memory.Global.Free),
			c.Host,
			device.UUID,
		)
		ch <- prometheus.MustNewConstMetric(
			gpuUsedMemDesc,
			prometheus.GaugeValue,
			float64(*ds.Memory.Global.Used),
			c.Host,
			device.UUID,
		)
		//pods metrics
		for _, process := range ds.Processes {
			p, err := getPodFromPid(process.PID, pods)
			if err != nil {
				fmt.Println(err)
				//c.logger.Error(err, "unable to get pod from pid: %s", process.PID)
				continue
			}
			n := p.Name
			ns := p.Namespace
			fmt.Printf(n,ns)
			ch <- prometheus.MustNewConstMetric(
				gpuPodUtilDesc,
				prometheus.GaugeValue,
				float64(process.Util.SmUtil),
				c.Host,
				device.UUID,
				n,
				ns,
			)
			ch <- prometheus.MustNewConstMetric(
				gpuPodMemUtilDesc,
				prometheus.GaugeValue,
				float64(process.Util.MemUtil),
				c.Host,
				device.UUID,
				n,
				ns,
			)
			ch <- prometheus.MustNewConstMetric(
				gpuPodEncUtilDesc,
				prometheus.GaugeValue,
				float64(process.Util.EncUtil),
				c.Host,
				device.UUID,
				n,
				ns,
			)
			ch <- prometheus.MustNewConstMetric(
				gpuPodDecUtilDesc,
				prometheus.GaugeValue,
				float64(process.Util.DecUtil),
				c.Host,
				device.UUID,
				n,
				ns,
			)
			ch <- prometheus.MustNewConstMetric(
				gpuPodMemDesc,
				prometheus.GaugeValue,
				float64(process.MemoryUsed),
				c.Host,
				device.UUID,
				n,
				ns,
			)

		}
	}
}

