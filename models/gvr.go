package models

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	QueryCRD = schema.GroupVersionResource{
		Group:    "sentinelfs.io",
		Version:  "v2",
		Resource: "filemonitors",
	}
)
