/*
Copyright 2021 Advanced Hosting

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

package ah

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
)

type servicePatcher struct {
	kclinet kubernetes.Interface
	origin  *v1.Service
	updated *v1.Service
}

func newServicePatcher(kclient kubernetes.Interface, origin *v1.Service) servicePatcher {
	return servicePatcher{
		kclinet: kclient,
		origin:  origin.DeepCopy(),
		updated: origin,
	}
}

func (sp *servicePatcher) Patch(ctx context.Context) error {

	originJSON, err := json.Marshal(sp.origin)
	if err != nil {
		return fmt.Errorf("failed to serialize current original object: %s", err)
	}

	updatedJSON, err := json.Marshal(sp.updated)
	if err != nil {
		return fmt.Errorf("failed to serialize modified updated object: %s", err)
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(originJSON, updatedJSON, v1.Service{})
	if err != nil {
		return fmt.Errorf("failed to create 2-way merge patch: %s", err)
	}

	if len(patch) == 0 || string(patch) == "{}" {
		return nil
	}

	_, err = sp.kclinet.CoreV1().Services(sp.origin.Namespace).Patch(ctx, sp.origin.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch service object %s/%s: %s", sp.origin.Namespace, sp.origin.Name, err)
	}

	return nil
}
