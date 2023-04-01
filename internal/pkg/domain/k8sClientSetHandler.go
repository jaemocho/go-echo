package domain

import (
	"backend/config"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClientSetHandler struct {
	clientSet *kubernetes.Clientset
}

func NewK8sClientSetHandler(cfg config.Config) *K8sClientSetHandler {

	clientConfig, _ := clientcmd.NewClientConfigFromBytes([]byte(cfg.ClusterToken))

	restConfig, _ := clientConfig.ClientConfig()

	clientSet, err := kubernetes.NewForConfig(restConfig)

	if err != nil {
		log.Printf("NewForConfig returned error: %v", err)
		return nil
	}

	return &K8sClientSetHandler{
		clientSet: clientSet,
	}

}

func (k *K8sClientSetHandler) GetPodList(namespace, cluster string) ([]*PodInfo, error) {

	pods, err := k.clientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Printf("GetPodList returned error: %v", err)
		return nil, err
	}

	returnPodList := make([]*PodInfo, len(pods.Items))

	for i, v := range pods.Items {

		totalCount := len(v.Status.ContainerStatuses)
		currentCount := 0
		podState := "Running"

		for j := 0; j < totalCount; j++ {
			if *v.Status.ContainerStatuses[j].Started {
				currentCount++
			} else {
				if v.Status.ContainerStatuses[j].State.Waiting != nil {
					podState = v.Status.ContainerStatuses[j].State.Waiting.Reason
				} else {
					podState = v.Status.ContainerStatuses[j].State.Terminated.Reason
				}
			}
		}

		var restartCount int32 = 0
		if len(v.Status.ContainerStatuses) > 0 {
			restartCount = v.Status.ContainerStatuses[0].RestartCount
		}

		returnPodList[i] = &PodInfo{
			Name:         v.Name,
			CurrentCount: currentCount,
			TotalCount:   totalCount,
			PodState:     podState,
			RestartCount: restartCount,
			PodIP:        v.Status.PodIP,
			NodeName:     v.Spec.NodeName,
			Age:          v.Status.StartTime.Time.String(),
		}
	}

	return returnPodList, nil

}

func (k *K8sClientSetHandler) GetPodEvent(namespace, podName string) (*v1.EventList, error) {
	events, err := k.clientSet.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "involvedObject.name=" + podName, TypeMeta: metav1.TypeMeta{Kind: "Pod"}})

	if err != nil {
		log.Printf("GetPodEvent returned error: %v", err)
		return nil, err
	}

	return events, nil
}

// 500 row
func (k *K8sClientSetHandler) GetPodLogs(namespace, podName string, previous bool) (*string, error) {

	count := int64(500)
	podLogOptions := v1.PodLogOptions{
		//Follow:    true,
		TailLines: &count,
		Previous:  previous,
	}

	podLogRequest := k.clientSet.CoreV1().Pods(namespace).GetLogs(podName, &podLogOptions)

	stream, err := podLogRequest.Stream(context.TODO())
	if err != nil {
		log.Printf("GetPodLogs returned error: %v", err)
		return nil, err
	}
	defer stream.Close()

	logs := ""

	for {
		buf := make([]byte, 2000)
		numBytes, err := stream.Read(buf)
		if numBytes == 0 {
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("GetPodLogs stream.Read returned error: %v", err)
			return nil, err
		}

		message := string(buf[:numBytes])
		logs += message
	}

	return &logs, nil
}

func (k *K8sClientSetHandler) GetPodDesc(namespace, podName string) (string, error) {

	result, err := k.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("GetPodDesc returned error: %v", err)
		return "", err
	}

	var podDesc bytes.Buffer

	podDesc.WriteString("Name : " + result.GetObjectMeta().GetName() + "\n")
	podDesc.WriteString("Namespace: " + result.GetObjectMeta().GetNamespace() + "\n")
	//podDesc.WriteString( "Priority: " + &result.Spec.Priority
	podDesc.WriteString("Node: " + result.Spec.NodeName + "\n")
	podDesc.WriteString("Start Time: " + result.Status.StartTime.String() + "\n")

	labelMap := result.GetObjectMeta().GetLabels()
	podDesc.WriteString("Labels: " + "\n")

	for key, v := range labelMap {
		podDesc.WriteString("             " + key + "=" + v + "\n")
	}
	annotationMap := result.ObjectMeta.GetAnnotations()
	podDesc.WriteString("Annotations: " + "\n")
	for key, v := range annotationMap {
		podDesc.WriteString("             " + key + ": " + v + "\n")
	}

	podDesc.WriteString("Status: " + string(result.Status.Phase) + "\n")
	podDesc.WriteString("IP: " + result.Status.PodIP + "\n")

	podDesc.WriteString("Containers: " + "\n")
	for _, v := range result.Status.ContainerStatuses {
		podDesc.WriteString("  " + v.Name + "\n")
		podDesc.WriteString("    Container ID: " + v.ContainerID + "\n")
		podDesc.WriteString("    Image: " + v.Image + "\n")
		podDesc.WriteString("    Image ID: " + v.ImageID + "\n")

		if *v.Started {
			podDesc.WriteString("    State: Running" + "\n")
			podDesc.WriteString("      Started: " + v.State.Running.StartedAt.String() + "\n")
		} else {
			if v.State.Waiting != nil {
				podDesc.WriteString("    State: " + v.State.Waiting.Reason + "\n")
			} else {
				podDesc.WriteString("    State: " + v.State.Terminated.Reason + "\n")
			}
		}
		podDesc.WriteString("    Ready: " + strconv.FormatBool(v.Ready) + "\n")
		podDesc.WriteString("    Restart Count: " + strconv.Itoa(int(v.RestartCount)) + "\n")

	}

	for _, v := range result.Spec.Containers {
		podDesc.WriteString("  " + v.Name + "\n")
		//podDesc.WriteString( "    Image: " + val.Image + "\n")
		for _, portVal := range v.Ports {
			podDesc.WriteString("    Port: " + strconv.Itoa(int(portVal.ContainerPort)) + "/" + string(portVal.Protocol) + "\n")
			podDesc.WriteString("    Host Port: " + strconv.Itoa(int(portVal.HostPort)) + "/" + string(portVal.Protocol) + "\n")

		}

		podDesc.WriteString("    Limits: " + "\n")
		podDesc.WriteString("      cpu: " + v.Resources.Limits.Cpu().String() + "\n")
		podDesc.WriteString("      memory: " + v.Resources.Limits.Memory().String() + "\n")

		podDesc.WriteString("    Requests: " + "\n")
		podDesc.WriteString("      cpu: " + v.Resources.Requests.Cpu().String() + "\n")
		podDesc.WriteString("      memory: " + v.Resources.Requests.Memory().String() + "\n")

		podDesc.WriteString("    VolumeMounts: " + "\n")
		for _, volumeMountsVal := range v.VolumeMounts {
			podDesc.WriteString("      " + volumeMountsVal.Name + "\n")
			podDesc.WriteString("        ReadOnly: " + strconv.FormatBool(volumeMountsVal.ReadOnly) + "\n")
			podDesc.WriteString("        MountPath: " + volumeMountsVal.MountPath + "\n")
			podDesc.WriteString("        SubPath: " + volumeMountsVal.SubPath + "\n")
		}
	}

	podDesc.WriteString("Volumes: " + "\n")
	for _, v := range result.Spec.Volumes {
		podDesc.WriteString("  " + v.Name + "\n")

		v := reflect.ValueOf(v.VolumeSource)
		typeOfField := v.Type()

		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsNil() {
				podDesc.WriteString("    type: " + typeOfField.Field(i).Name + "\n")
				podDesc.WriteString("    name: " + fmt.Sprintf("%v", v.Field(i).Interface()) + "\n")
			}
		}

	}

	podDesc.WriteString("PodStatus: " + "\n")
	for i, v := range result.Status.Conditions {
		podDesc.WriteString("  [" + strconv.Itoa(i) + "]\n")
		podDesc.WriteString("  Type: " + string(v.Type) + "\n")
		podDesc.WriteString("  Type: " + string(v.Type) + "\n")
		podDesc.WriteString("  Status: " + string(v.Status) + "\n")
		podDesc.WriteString("  [" + strconv.Itoa(i) + "]\n")
		podDesc.WriteString("  Type: " + string(v.Type) + "\n")
		podDesc.WriteString("  Status: " + string(v.Status) + "\n")
		podDesc.WriteString("  LastProbeTime: " + v.LastProbeTime.String() + "\n")
		podDesc.WriteString("  LastTransitionTime: " + v.LastTransitionTime.String() + "\n")
		podDesc.WriteString("  Reason: " + string(v.Reason) + "\n")
		podDesc.WriteString("  Messae: " + string(v.Message) + "\n")

	}

	return podDesc.String(), nil

}
