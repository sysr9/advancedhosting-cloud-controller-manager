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
	"reflect"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/advancedhosting/advancedhosting-cloud-controller-manager/advancedhosting/mocks"
	"github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func testInstanceGetResponse() *ah.Instance {
	return &ah.Instance{
		ID:                         "test-worker-id",
		Name:                       "k8s-worker-test",
		PrimaryInstanceIPAddressID: "test_address_id",
		IPAddresses: []ah.InstanceIPAddress{
			{
				ID:      "test_address_id",
				Address: "1.2.3.4",
			},
		},
		PrivateNetworks: []ah.InstancePrivateNetwork{
			{
				InstancePrivateNetworkInfo: ah.InstancePrivateNetworkInfo{
					IP: "1.0.0.1",
				},
			},
		},
		Image: &ah.InstanceImage{
			Image: &ah.Image{
				Slug: "test-slug",
			},
		},
	}
}

func testInstanceListResponse() []ah.Instance {
	return []ah.Instance{*testInstanceGetResponse()}
}

func TestInstances_NodeAddresses(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedInstancesAPI.EXPECT().List(gomock.Any(), gomock.Any()).Return(testInstanceListResponse(), nil, nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	addresses, err := instances.NodeAddresses(context.TODO(), "k8s-worker-test")

	expectedResult := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "k8s-worker-test",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "1.2.3.4",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "1.0.0.1",
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, addresses) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, addresses)
	}

}

func TestInstances_NodeAddressesByProviderID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedInstancesAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(testInstanceGetResponse(), nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	addresses, err := instances.NodeAddressesByProviderID(context.TODO(), "advancedhosting://test-worker-id")

	expectedResult := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "k8s-worker-test",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "1.2.3.4",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "1.0.0.1",
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, addresses) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, addresses)
	}

}

func TestInstances_InstanceID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedInstancesAPI.EXPECT().List(gomock.Any(), gomock.Any()).Return(testInstanceListResponse(), nil, nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	addresses, err := instances.InstanceID(context.TODO(), "k8s-worker-test")

	expectedResult := "test-worker-id"

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, addresses) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, addresses)
	}

}

func TestInstances_InstanceType(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedInstancesAPI.EXPECT().List(gomock.Any(), gomock.Any()).Return(testInstanceListResponse(), nil, nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	addresses, err := instances.InstanceType(context.TODO(), "k8s-worker-test")

	expectedResult := "test-slug"

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, addresses) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, addresses)
	}

}

func TestInstances_InstanceTypeByProviderID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedInstancesAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(testInstanceGetResponse(), nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	addresses, err := instances.InstanceTypeByProviderID(context.TODO(), "advancedhosting://test-worker-id")

	expectedResult := "test-slug"

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, addresses) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, addresses)
	}

}

func TestInstances_CurrentNodeName(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	addresses, err := instances.CurrentNodeName(context.TODO(), "test-hostname")

	expectedResult := types.NodeName("test-hostname")

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, addresses) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, addresses)
	}

}

func TestInstances_InstanceExistsByProviderID_InstanceExists(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	expectedInstance := &ah.Instance{}

	mockedInstancesAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(expectedInstance, nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	isExist, err := instances.InstanceExistsByProviderID(context.TODO(), "advancedhosting://test-worker-id")

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !isExist {
		t.Errorf("Unexpected result")
	}

}

func TestInstances_InstanceExistsByProviderID_InstanceNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	mockedInstancesAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, ah.ErrResourceNotFound)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	isExist, err := instances.InstanceExistsByProviderID(context.TODO(), "advancedhosting://test-worker-id")

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if isExist {
		t.Errorf("Unexpected result")
	}

}

func TestInstances_InstanceShutdownByProviderID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	expectedInstance := &ah.Instance{
		State: ah.InstanceShutDownStatus,
	}

	mockedInstancesAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(expectedInstance, nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	isShutdown, err := instances.InstanceShutdownByProviderID(context.TODO(), "advancedhosting://test-worker-id")

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !isShutdown {
		t.Errorf("Unexpected result")
	}

}

func TestInstances_InstanceShutdownByProviderID_InstanceInRunningState(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedInstancesAPI := mocks.NewMockInstancesAPI(ctrl)

	expectedInstance := &ah.Instance{
		State: "running",
	}

	mockedInstancesAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(expectedInstance, nil)

	mockedClient := &ah.APIClient{Instances: mockedInstancesAPI}
	instances := newInstances(mockedClient)

	isShutdown, err := instances.InstanceShutdownByProviderID(context.TODO(), "advancedhosting://test-worker-id")

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if isShutdown {
		t.Errorf("Unexpected result")
	}

}
