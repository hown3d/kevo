package kubernetes

import (
	"context"
	"log"
	"sync"

	"github.com/hown3d/container-image-scanner/pkg/types"
	"github.com/hown3d/container-image-scanner/pkg/util/imageutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k kubernetesFetcher) GetImages(ctx context.Context) ([]types.Image, error) {
	var images []types.Image
	// empty namespace, to fetch from all namespaces
	podsClient := k.client.CoreV1().Pods("")
	pods, err := podsClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, pod := range pods.Items {
		wg.Add(1)
		go func(pod corev1.Pod) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			images = append(images, getImagesFromContainerStatus(pod.Status.ContainerStatuses)...)
			images = append(images, getImagesFromContainerStatus(pod.Status.InitContainerStatuses)...)
			images = append(images, getImagesFromContainerStatus(pod.Status.EphemeralContainerStatuses)...)
		}(pod)
	}
	wg.Wait()

	return images, nil
}

func getImagesFromContainerStatus(status []corev1.ContainerStatus) []types.Image {
	var images []types.Image
	for _, container := range status {
		name, tag := imageutil.SplitImageFromString(container.Image)
		log.Printf("Adding image %v:%v", name, tag)
		images = append(images, types.Image{Name: name, Tag: tag})
	}
	return images
}
