package kubernetes

import (
	"context"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)


func QueryCRDsInNamespace(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource, namespace string) error {
	log.Printf("Querying CRDs in namespace: %s", namespace)

	list, err := client.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("No CRDs found in namespace %s", namespace)
			return nil
		}
		return fmt.Errorf("failed to list CRDs in namespace %s: %v", namespace, err)
	}

	log.Printf("Found %d CRDs in namespace: %s", len(list.Items), namespace)

	for _, item := range list.Items {
		name := item.GetName()
		log.Printf("CRD: %s", name)
	}

	return nil
}



func UpdateCRDWithFileInfo(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource, namespace string) error {
	log.Printf("Updating CRDs with file information in namespace: %s", namespace)

	// Get existing CRDs
	list, err := client.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("No CRDs found to update in namespace %s", namespace)
			return nil
		}
		return fmt.Errorf("failed to list CRDs for update: %v", err)
	}

	for _, item := range list.Items {
		name := item.GetName()
		log.Printf("Updating CRD: %s", name)

		fileInfo := []interface{}{
			map[string]interface{}{
				"inode":   12345,
				"cGroupId": "12232",
				"path":    "/tmp/example.txt",
			},
		}

		// Update the status with file information
		if err := unstructured.SetNestedSlice(item.Object, fileInfo, "status", "files"); err != nil {
			log.Printf("Failed to set file info for CRD %s: %v", name, err)
			continue
		}

		// Update the CRD
		_, err := client.Resource(gvr).Namespace(namespace).UpdateStatus(ctx, &item, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Failed to update CRD %s: %v", name, err)
			continue
		}

		log.Printf("Successfully updated CRD: %s", name)
	}

	return nil
}