/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/pkg/internal/scheme"
)

const (
	deploymentKind          = "Deployment"
	controllerContainerName = "manager"
)

// InspectImages identifies the container images required to install the objects defined in the objs.
// NB. The implemented approach is specific for the provider components YAML & for the cert-manager manifest; it is not
// intended to cover all the possible objects used to deploy containers existing in Kubernetes.
func InspectImages(objs []unstructured.Unstructured) ([]string, error) {
	images := []string{}

	for i := range objs {
		o := objs[i]
		if o.GetKind() == deploymentKind {
			d := &appsv1.Deployment{}
			if err := scheme.Scheme.Convert(&o, d, nil); err != nil {
				return nil, err
			}

			for _, c := range d.Spec.Template.Spec.Containers {
				images = append(images, c.Image)
			}

			for _, c := range d.Spec.Template.Spec.InitContainers {
				images = append(images, c.Image)
			}
		}
	}

	return images, nil
}
