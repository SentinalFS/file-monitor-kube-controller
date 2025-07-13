package main

import (
	"context"
	k8s "fileMonitorKubeController/kubernetes"
	models "fileMonitorKubeController/models"
	"log"
	"os"
	"time"
)



func main() {
	_, dynamicClient, err := k8s.InitKubernetesClients()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes clients: %v", err)
	}

	crdGVR := models.QueryCRD

	ctx := context.Background()

	for {
		log.Println("Querying CRDs...")

		namespace := os.Getenv("NAMESPACE")
		if err := k8s.QueryCRDsInNamespace(ctx, dynamicClient, crdGVR, namespace); err != nil {
			log.Printf("Error querying CRDs in namespace %s: %v", namespace, err)
		}

		time.Sleep(30 * time.Second)
	}
}
