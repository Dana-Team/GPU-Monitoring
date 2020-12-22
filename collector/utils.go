package collector

import (
	"fmt"
	"github.com/Dana-Team/gonvml"
	v1 "k8s.io/api/core/v1"
	"regexp"
	"strings"
	"os/exec"
	"errors"
)

func getDevices() ([]*gonvml.Device, error){
	var devices []*gonvml.Device
	numOfDevices, err := gonvml.GetDeviceCount()
	if err != nil {
		return nil, err
	}
	for i := 0 ; i < int(numOfDevices) ; i++ {
		device, err := gonvml.NewDevice(uint(i))
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func getDeviceStatus(device gonvml.Device) (*gonvml.DeviceStatus, error){
	deviceStatus, err := device.Status()
	if err != nil{
		return nil, err
	}
	return deviceStatus, nil
}

func getPodFromPid(pid uint, podList v1.PodList) (*v1.Pod, error) {
	cid, err := getContinerIdFromPid(pid)
	if err != nil {
		return nil, err
	}
	pod, err := getPodFromContinarID(cid, podList)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func getContinerIdFromPid(pid uint) (string, error) {
	re := regexp.MustCompile("crio-.*scope")
	cmd := exec.Command("systemctl", "status", fmt.Sprint(pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	cid := re.FindStringSubmatch(string(out))[0]
	cid = strings.TrimPrefix(cid, "crio-")
	cid = strings.TrimSuffix(cid, ".scope")
	return cid, nil
}

func getPodFromContinarID(cid string, podList v1.PodList) (*v1.Pod, error) {
	for _, pod := range podList.Items {
		for _, container := range pod.Status.ContainerStatuses{
			if container.ContainerID[8:] == cid {
				return &pod, nil
			}
		}
	}
	return nil, errors.New("Not Found")
}

