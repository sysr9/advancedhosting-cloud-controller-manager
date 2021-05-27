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
	"fmt"
	"reflect"
	"testing"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	"github.com/advancedhosting/advancedhosting-cloud-controller-manager/advancedhosting/mocks"
	"github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func testAnnotaions() map[string]string {
	return map[string]string{
		ServiceAnnotationLoadBalancerID:                            "test-lb-id",
		ServiceAnnotationLoadBalancerName:                          "test-lb-name",
		ServiceAnnotationLoadBalancerBalancingAlgorithm:            "round_robin",
		ServiceAnnotationLoadBalancerEnableHealthCheck:             "true",
		ServiceAnnotationLoadBalancerHealthCheckType:               "tcp",
		ServiceAnnotationLoadBalancerHealthCheckURL:                "",
		ServiceAnnotationLoadBalancerHealthCheckInterval:           "5",
		ServiceAnnotationLoadBalancerHealthCheckTimeout:            "2",
		ServiceAnnotationLoadBalancerHealthCheckUnhealthyThreshold: "5",
		ServiceAnnotationLoadBalancerHealthCheckHealthyThreshold:   "5",
		ServiceAnnotationLoadBalancerHealthCheckPort:               "8080",
	}
}

func testPorts() []v1.ServicePort {
	return []v1.ServicePort{
		{
			Name:     "test-port",
			Protocol: "tcp",
			Port:     int32(80),
			NodePort: int32(30000),
		},
	}
}

func testService(kclient kubernetes.Interface, annotations map[string]string, ports []v1.ServicePort) *v1.Service {

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "default",
			Name:        "test-service",
			Annotations: annotations,
		},
		Spec: v1.ServiceSpec{
			Ports: ports,
		},
	}
	if _, err := kclient.CoreV1().Services(svc.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil {
		panic(err)
	}
	return svc
}

func testNodes() []*v1.Node {
	return []*v1.Node{
		{
			Spec: v1.NodeSpec{
				ProviderID: "advancedhosting://test-cloud-server-1",
			},
		},
		{
			Spec: v1.NodeSpec{
				ProviderID: "advancedhosting://test-cloud-server-2",
			},
		},
	}
}

func testLBGetResponse() *ah.LoadBalancer {
	return &ah.LoadBalancer{
		ID:   "test-lb-id",
		Name: "test-lb-name",
		IPAddresses: []ah.LBIPAddress{
			{
				Address: "1.2.3.4",
			},
		},
		State:              "active",
		BalancingAlgorithm: "round_robin",
		HealthChecks: []ah.LBHealthCheck{
			{
				Type:               "tcp",
				URL:                "",
				Interval:           5,
				Timeout:            2,
				UnhealthyThreshold: 5,
				HealthyThreshold:   5,
				Port:               8080,
			},
		},
		ForwardingRules: []ah.LBForwardingRule{
			{
				RequestProtocol:       "tcp",
				RequestPort:           80,
				CommunicationProtocol: "tcp",
				CommunicationPort:     30000,
			},
		},
		BackendNodes: []ah.LBBackendNode{
			{
				ID:            "test-id-1",
				CloudServerID: "test-cloud-server-1",
			},
			{
				ID:            "test-id-2",
				CloudServerID: "test-cloud-server-2",
			},
		},
	}
}

func TestLoadBalancers_GetLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)

	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Eq("test-lb-id")).Return(testLBGetResponse(), nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}

	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerID] = "test-lb-id"

	svc := testService(clusterInfo.kclient, anno, testPorts())

	status, exists, err := loadBalancers.GetLoadBalancer(context.TODO(), "test-sluster-name", svc)

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !exists {
		t.Errorf("Unexpected value: %v", exists)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_GetLoadBalancerName(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)

	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	name := loadBalancers.GetLoadBalancerName(context.TODO(), "test-sluster-name", testService(clusterInfo.kclient, testAnnotaions(), testPorts()))

	expectedResult := "test-lb-name"

	if !reflect.DeepEqual(expectedResult, name) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, name)
	}

}

func TestLoadBalancers_EnsureLoadBalancerCreateNewLB(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, ah.ErrResourceNotFound)

	testLB := testLBGetResponse()
	testLB.State = "creating"

	mockedLBAPI.EXPECT().Create(gomock.Any(), gomock.Any()).Return(testLB, nil)

	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}
	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	svc := testService(clusterInfo.kclient, testAnnotaions(), testPorts())

	_, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	if err.Error() != fmt.Sprintf("Load balancer is not active yet: %s", testLB.State) {
		t.Errorf("Unexpected Error: %v", err)
	}

	updatedService, err := clusterInfo.kclient.CoreV1().Services(svc.Namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Error getting service: %v", err)
	}

	if lbID := updatedService.Annotations[ServiceAnnotationLoadBalancerID]; testLB.ID != lbID {
		t.Errorf("Unexpected result, expected %v. got: %v", testLB.ID, lbID)
	}

}

func TestLoadBalancers_EnsureNotReadyLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)

	testLB := testLBGetResponse()
	testLB.State = "creating"
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(testLB, nil)

	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}
	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	svc := testService(clusterInfo.kclient, testAnnotaions(), testPorts())
	_, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	if err.Error() != fmt.Sprintf("Load balancer is not active yet: %s", testLB.State) {
		t.Errorf("Unexpected Error: %v", err)
	}

}

func TestLoadBalancers_EnsureActiveLoadBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)

	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Return(testLBGetResponse(), nil)

	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}
	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	svc := testService(clusterInfo.kclient, testAnnotaions(), testPorts())
	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_UpdateBalancingAlgorithm(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(2).Return(testLBGetResponse(), nil)
	mockedLBAPI.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Eq(&ah.LoadBalancerUpdateRequest{BalancingAlgorithm: "least_requests"})).Return(nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerBalancingAlgorithm] = "least_requests"
	svc := testService(clusterInfo.kclient, anno, testPorts())

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_UpdateLBName(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(2).Return(testLBGetResponse(), nil)
	mockedLBAPI.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Eq(&ah.LoadBalancerUpdateRequest{Name: "test2"})).Return(nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerName] = "test2"
	svc := testService(clusterInfo.kclient, anno, testPorts())

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_UpdateHC(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLBGetResponse(), nil)
	updateRequest := &ah.LBHealthCheckUpdateRequest{
		Type:               "http",
		URL:                "/",
		Interval:           6,
		Timeout:            3,
		UnhealthyThreshold: 4,
		HealthyThreshold:   3,
		Port:               9090,
	}
	mockedLBAPI.EXPECT().UpdateHealthCheck(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(updateRequest)).Return(nil)
	mockedLBAPI.EXPECT().GetHealthCheck(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&ah.LBHealthCheck{State: "updating"}, nil)
	mockedLBAPI.EXPECT().GetHealthCheck(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&ah.LBHealthCheck{State: "active"}, nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerHealthCheckType] = "http"
	anno[ServiceAnnotationLoadBalancerHealthCheckURL] = "/"
	anno[ServiceAnnotationLoadBalancerHealthCheckInterval] = "6"
	anno[ServiceAnnotationLoadBalancerHealthCheckTimeout] = "3"
	anno[ServiceAnnotationLoadBalancerHealthCheckUnhealthyThreshold] = "4"
	anno[ServiceAnnotationLoadBalancerHealthCheckHealthyThreshold] = "3"
	anno[ServiceAnnotationLoadBalancerHealthCheckPort] = "9090"
	svc := testService(clusterInfo.kclient, anno, testPorts())

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_DisableHC(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLBGetResponse(), nil)
	mockedLBAPI.EXPECT().DeleteHealthCheck(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockedLBAPI.EXPECT().GetHealthCheck(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&ah.LBHealthCheck{State: "deleting"}, nil)
	mockedLBAPI.EXPECT().GetHealthCheck(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, ah.ErrResourceNotFound)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerEnableHealthCheck] = "false"
	svc := testService(clusterInfo.kclient, anno, testPorts())

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_EnableHC(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	testLB := testLBGetResponse()
	testLB.HealthChecks = []ah.LBHealthCheck{}
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLB, nil)
	createRequest := &ah.LBHealthCheckCreateRequest{
		Type:               "tcp",
		URL:                "",
		Interval:           5,
		Timeout:            2,
		UnhealthyThreshold: 5,
		HealthyThreshold:   5,
		Port:               8080,
	}
	mockedLBAPI.EXPECT().CreateHealthCheck(gomock.Any(), gomock.Any(), gomock.Eq(createRequest)).Return(&ah.LBHealthCheck{ID: "test-id"}, nil)
	mockedLBAPI.EXPECT().GetHealthCheck(gomock.Any(), gomock.Any(), gomock.Eq("test-id")).Times(1).Return(&ah.LBHealthCheck{State: "creating"}, nil)
	mockedLBAPI.EXPECT().GetHealthCheck(gomock.Any(), gomock.Any(), gomock.Eq("test-id")).Times(1).Return(&ah.LBHealthCheck{State: "active"}, nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	svc := testService(clusterInfo.kclient, testAnnotaions(), testPorts())

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_UpdateForwardingRules(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)

	testLB := testLBGetResponse()
	testLB.ForwardingRules = []ah.LBForwardingRule{
		{
			ID:                    "fr-to-update-id",
			RequestProtocol:       "http",
			RequestPort:           80,
			CommunicationProtocol: "tcp",
			CommunicationPort:     30000,
		},
		{
			ID:                    "fr-to-delete-id",
			RequestProtocol:       "tcp",
			RequestPort:           70,
			CommunicationProtocol: "tcp",
			CommunicationPort:     30070,
		},
	}
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLB, nil)

	createRequest := &ah.LBForwardingRuleCreateRequest{
		RequestProtocol:       "tcp",
		RequestPort:           90,
		CommunicationProtocol: "tcp",
		CommunicationPort:     30003,
	}

	// update fr
	mockedLBAPI.EXPECT().DeleteForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("fr-to-update-id")).Return(nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("fr-to-update-id")).Times(1).Return(&ah.LBForwardingRule{State: "deleting"}, nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("fr-to-update-id")).Times(1).Return(nil, ah.ErrResourceNotFound)

	updateRequest := &ah.LBForwardingRuleCreateRequest{
		RequestProtocol:       "tcp",
		RequestPort:           80,
		CommunicationProtocol: "tcp",
		CommunicationPort:     30002,
	}
	mockedLBAPI.EXPECT().CreateForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq(updateRequest)).Return(&ah.LBForwardingRule{ID: "test-updated-id"}, nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("test-updated-id")).Times(1).Return(&ah.LBForwardingRule{State: "creating"}, nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("test-updated-id")).Times(1).Return(&ah.LBForwardingRule{State: "active"}, nil)

	// create fr
	mockedLBAPI.EXPECT().CreateForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq(createRequest)).Return(&ah.LBForwardingRule{ID: "test-new-id"}, nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("test-new-id")).Times(1).Return(&ah.LBForwardingRule{State: "creating"}, nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("test-new-id")).Times(1).Return(&ah.LBForwardingRule{State: "active"}, nil)

	// delete unused fr
	mockedLBAPI.EXPECT().DeleteForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("fr-to-delete-id")).Return(nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("fr-to-delete-id")).Times(1).Return(&ah.LBForwardingRule{State: "deleting"}, nil)
	mockedLBAPI.EXPECT().GetForwardingRule(gomock.Any(), gomock.Any(), gomock.Eq("fr-to-delete-id")).Times(1).Return(nil, ah.ErrResourceNotFound)

	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	ports := []v1.ServicePort{
		{
			Name:     "updated-port",
			Protocol: "tcp",
			Port:     int32(80),
			NodePort: int32(30002),
		},
		{
			Name:     "new-port",
			Protocol: "tcp",
			Port:     int32(90),
			NodePort: int32(30003),
		},
	}

	svc := testService(clusterInfo.kclient, testAnnotaions(), ports)

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, testNodes())

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_UpdateBackendNodes(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)

	testLB := testLBGetResponse()
	testLB.BackendNodes = []ah.LBBackendNode{
		{
			ID:            "test-backend-node-id-1",
			CloudServerID: "test-cloud-server-1",
		},
		{
			ID:            "test-backend-node-id-2",
			CloudServerID: "test-cloud-server-2",
		},
	}

	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLB, nil)

	// add bns
	addRequest := []string{"test-cloud-server-3", "test-cloud-server-4"}
	addResponse := []ah.LBBackendNode{
		{
			ID:            "test-backend-node-id-3",
			CloudServerID: "test-cloud-server-3",
			State:         "creating",
		},
		{
			ID:            "test-backend-node-id-4",
			CloudServerID: "test-cloud-server-4",
			State:         "creating",
		},
	}
	mockedLBAPI.EXPECT().AddBackendNodes(gomock.Any(), gomock.Any(), gomock.Eq(addRequest)).Return(addResponse, nil)
	mockedLBAPI.EXPECT().ListBackendNodes(gomock.Any(), gomock.Any()).Times(1).Return(addResponse, nil)

	mockedLBAPI.EXPECT().ListBackendNodes(gomock.Any(), gomock.Any()).Times(1).Return([]ah.LBBackendNode{{ID: "test-backend-node-id-3", State: "active"}, {ID: "test-backend-node-id-4", State: "creating"}}, nil)

	mockedLBAPI.EXPECT().ListBackendNodes(gomock.Any(), gomock.Any()).Times(1).Return([]ah.LBBackendNode{{ID: "test-backend-node-id-3", State: "active"}, {ID: "test-backend-node-id-4", State: "active"}}, nil)

	// delete unused bn
	mockedLBAPI.EXPECT().DeleteBackendNode(gomock.Any(), gomock.Any(), gomock.Eq("test-backend-node-id-2")).Return(nil)
	mockedLBAPI.EXPECT().GetBackendNode(gomock.Any(), gomock.Any(), gomock.Eq("test-backend-node-id-2")).Times(1).Return(&ah.LBBackendNode{State: "deleting"}, nil)
	mockedLBAPI.EXPECT().GetBackendNode(gomock.Any(), gomock.Any(), gomock.Eq("test-backend-node-id-2")).Times(1).Return(nil, ah.ErrResourceNotFound)

	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	svc := testService(clusterInfo.kclient, testAnnotaions(), testPorts())

	nodes := []*v1.Node{
		{
			Spec: v1.NodeSpec{
				ProviderID: "advancedhosting://test-cloud-server-1",
			},
		},
		{
			Spec: v1.NodeSpec{
				ProviderID: "advancedhosting://test-cloud-server-3",
			},
		},
		{
			Spec: v1.NodeSpec{
				ProviderID: "advancedhosting://test-cloud-server-4",
			},
		},
	}

	status, err := loadBalancers.EnsureLoadBalancer(context.TODO(), "test-sluster-name", svc, nodes)

	expectedResult := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: "1.2.3.4",
			},
		},
	}

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if !reflect.DeepEqual(expectedResult, status) {
		t.Errorf("Unexpected result, expected %v. got: %v", expectedResult, status)
	}

}

func TestLoadBalancers_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLBGetResponse(), nil)
	mockedLBAPI.EXPECT().Delete(gomock.Any(), gomock.Eq("test-lb-id")).Times(1).Return(nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerID] = "test-lb-id"

	svc := testService(clusterInfo.kclient, anno, testPorts())

	err := loadBalancers.EnsureLoadBalancerDeleted(context.TODO(), "test-sluster-name", svc)

	if err.Error() != "LB deletion has been started" {
		t.Errorf("Unexpected Error: %v", err)
	}
}

func TestLoadBalancers_EnsureDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	testLB := testLBGetResponse()
	testLB.State = "deleting"
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testLB, nil)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerID] = "test-lb-id"

	svc := testService(clusterInfo.kclient, anno, testPorts())

	err := loadBalancers.EnsureLoadBalancerDeleted(context.TODO(), "test-sluster-name", svc)

	if err.Error() != "Load balancer is already in deletion state" {
		t.Errorf("Unexpected Error: %v", err)
	}
}

func TestLoadBalancers_EnsureAlreadyDeletedLB(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedLBAPI := mocks.NewMockLoadBalancersAPI(ctrl)
	mockedLBAPI.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(nil, ah.ErrResourceNotFound)
	mockedClient := &ah.APIClient{LoadBalancers: mockedLBAPI}

	clusterInfo := &clusterInfo{kclient: fake.NewSimpleClientset()}
	loadBalancers := newLoadbalancers(mockedClient, clusterInfo)

	anno := testAnnotaions()
	anno[ServiceAnnotationLoadBalancerID] = "test-lb-id"

	svc := testService(clusterInfo.kclient, anno, testPorts())

	err := loadBalancers.EnsureLoadBalancerDeleted(context.TODO(), "test-sluster-name", svc)

	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

}
