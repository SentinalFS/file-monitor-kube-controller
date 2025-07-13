package models

type FileInfo struct {
	Path   string `json:"path"`
	CGroupID string `json:"cGroupID"`
	Inode uint64 `json:"inode"`
}

type PodSelector struct {
	Labels map[string]string `json:"labels"`
}

type Dir struct {
	Path string `json:"path"`
	Depth int    `json:"depth"`
}

type FileMonitorCRD struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		PodSelector PodSelector `json:"podSelector"`
		Dirs        []Dir       `json:"dirs"`
		Files       []FileInfo  `json:"files"`
	} `json:"spec"`
}