package kubernetes

import (
	"context"
	"fmt"

	"github.com/hown3d/kevo/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *kubernetesFetcher) getImagePullSecret(ctx context.Context, image *types.Image, namespace string, secretRefs []corev1.LocalObjectReference) error {
	client := k.client.CoreV1().Secrets(namespace)
	for _, secret := range secretRefs {
		var auth types.RegistryAuth
		secret, err := client.Get(ctx, secret.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		var data []byte
		switch secret.Type {
		// always has secret data .dockerconfigjson
		case corev1.SecretTypeDockerConfigJson:
			data = secret.Data[".dockerconfigjson"]
		default:
			return fmt.Errorf("Secret %v is not of type dockerconfigjson: %w", secret.Name, err)
		}

		err = auth.UnmarshalRegistryAuthJSON(data)
		if err != nil {
			return err
		}

		imageDomain, err := image.RegistryDomain()
		if err != nil {
			return err
		}
		// check if the registry image is
		if imageDomain == auth.Domain {
			image.Auth = auth
		}
	}
	return fmt.Errorf("No secret from %v matched the image domain", secretRefs)
}
