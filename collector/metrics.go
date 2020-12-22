package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	deviceLables = []string{"host", "device"}
	podLabels    = []string{"host", "device", "pod", "namespace"}
)

//device metrics
var (
	gpuUtilDesc = prometheus.NewDesc(
		"dana_gpu_util",
		"GPU Utilization",
		deviceLables, nil,
	)
	gpuMemUtilDesc = prometheus.NewDesc(
		"dana_gpu_mem_util",
		"GPU Memory Utilization",
		deviceLables, nil,
	)
	gpuEncUtilDesc = prometheus.NewDesc(
		"dana_gpu_enc_util",
		"GPU Encoder Utilization",
		deviceLables, nil,
	)
	gpuDecUtilDesc = prometheus.NewDesc(
		"dana_gpu_dec_util",
		"GPU Decoder Utilization",
		deviceLables, nil,
	)
	gpuFreeMemDesc = prometheus.NewDesc(
		"dana_gpu_free_mem",
		"GPU Free Memory",
		deviceLables, nil,
	)
	gpuUsedMemDesc = prometheus.NewDesc(
		"dana_gpu_used_mem",
		"GPU Used Memory",
		deviceLables, nil,
	)
)

var (
	gpuPodUtilDesc = prometheus.NewDesc(
		"dana_gpu_pod_util",
		"GPU Pod Utilization",
		podLabels, nil,
	)
	gpuPodMemUtilDesc = prometheus.NewDesc(
		"dana_gpu_pod_mem_util",
		"GPU Pod Memory Utilization",
		podLabels, nil,
	)
	gpuPodEncUtilDesc = prometheus.NewDesc(
		"dana_gpu_pod_enc_util",
		"GPU Pod Encoder Utilization",
		podLabels, nil,
	)
	gpuPodDecUtilDesc = prometheus.NewDesc(
		"dana_gpu_pod_dec_util",
		"GPU Pod Decoder Utilization",
		podLabels, nil,
	)
	gpuPodMemDesc = prometheus.NewDesc(
		"dana_gpu_pod_mem",
		"GPU Pod Memory Used",
		podLabels, nil,
	)
)