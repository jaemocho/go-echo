package domain

type PodInfo struct {
	Name         string `json:"name"`
	CurrentCount int    `json:"currentCount"`
	TotalCount   int    `json:"totalCount"`
	RestartCount int32  `json:"restartCount"`
	PodState     string `json:"podState"`
	PodIP        string `json:"podIP"`
	NodeName     string `json:"nodeName"`
	Age          string `json:"age"`
}
