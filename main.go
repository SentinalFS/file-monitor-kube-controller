package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// FileMonitorCRD represents the structure of our custom resource
type FileMonitorCRD struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Path      string `json:"path"`
		Namespace string `json:"namespace"`
	} `json:"spec"`
	Status struct {
		Files []FileInfo `json:"files,omitempty"`
	} `json:"status,omitempty"`
}

// FileInfo represents file information to be stored in CRD
type FileInfo struct {
	Name    string    `json:"name"`
	Inode   uint64    `json:"inode"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"isDir"`
}

func main() {
	// Initialize Kubernetes client
	_, dynamicClient, err := initKubernetesClients()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes clients: %v", err)
	}

	// Define the CRD GroupVersionResource
	crdGVR := schema.GroupVersionResource{
		Group:    "sentinalfs.io",
		Version:  "v1",
		Resource: "filemonitors",
	}

	ctx := context.Background()

	// In while true watch for changes in the crds
	for {
		log.Println("Querying CRDs...")

		// Query all CRDs in all namespaces
		if err := queryCRDs(ctx, dynamicClient, crdGVR); err != nil {
			log.Printf("Error querying CRDs: %v", err)
		}

		// Query CRDs in specific namespace
		namespace := "default"
		if err := queryCRDsInNamespace(ctx, dynamicClient, crdGVR, namespace); err != nil {
			log.Printf("Error querying CRDs in namespace %s: %v", namespace, err)
		}

		// append data into crds accordingly to the namespace it is in
		if err := updateCRDWithFileInfo(ctx, dynamicClient, crdGVR, namespace); err != nil {
			log.Printf("Error updating CRD with file info: %v", err)
		}

		// Wait before next iteration
		time.Sleep(30 * time.Second)
	}
}

// initKubernetesClients initializes both regular and dynamic Kubernetes clients
func initKubernetesClients() (kubernetes.Interface, dynamic.Interface, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first (when running inside a pod)
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig (for local development)
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to build config: %v", err)
		}
	}

	// Create regular client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	// Create dynamic client for CRDs
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create dynamic client: %v", err)
	}

	return clientset, dynamicClient, nil
}

// queryCRDs queries all CRDs across all namespaces
func queryCRDs(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource) error {
	log.Println("Querying CRDs in all namespaces...")

	list, err := client.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("CRD not found - it may not be installed yet")
			return nil
		}
		return fmt.Errorf("failed to list CRDs: %v", err)
	}

	log.Printf("Found %d CRDs across all namespaces", len(list.Items))

	for _, item := range list.Items {
		name := item.GetName()
		namespace := item.GetNamespace()
		log.Printf("CRD: %s in namespace: %s", name, namespace)

		// Print spec if available
		if spec, found, err := unstructured.NestedMap(item.Object, "spec"); err == nil && found {
			log.Printf("  Spec: %+v", spec)
		}

		// Print status if available
		if status, found, err := unstructured.NestedMap(item.Object, "status"); err == nil && found {
			log.Printf("  Status: %+v", status)
		}
	}

	return nil
}

// queryCRDsInNamespace queries CRDs in a specific namespace
func queryCRDsInNamespace(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource, namespace string) error {
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

// updateCRDWithFileInfo updates CRDs with file information such as inode, file name etc
func updateCRDWithFileInfo(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource, namespace string) error {
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

		// Example: Add file information to the status
		// In a real implementation, you would scan the actual filesystem
		fileInfo := []interface{}{
			map[string]interface{}{
				"name":    "example.txt",
				"inode":   12345,
				"size":    1024,
				"modTime": time.Now().Format(time.RFC3339),
				"path":    "/tmp/example.txt",
				"isDir":   false,
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
