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

package mocks

import (
	context "context"
	reflect "reflect"

	ah "github.com/advancedhosting/advancedhosting-api-go/ah"
	gomock "github.com/golang/mock/gomock"
)

// MockLoadBalancersAPI is a mock of LoadBalancersAPI interface.
type MockLoadBalancersAPI struct {
	ctrl     *gomock.Controller
	recorder *MockLoadBalancersAPIMockRecorder
}

// MockLoadBalancersAPIMockRecorder is the mock recorder for MockLoadBalancersAPI.
type MockLoadBalancersAPIMockRecorder struct {
	mock *MockLoadBalancersAPI
}

// NewMockLoadBalancersAPI creates a new mock instance.
func NewMockLoadBalancersAPI(ctrl *gomock.Controller) *MockLoadBalancersAPI {
	mock := &MockLoadBalancersAPI{ctrl: ctrl}
	mock.recorder = &MockLoadBalancersAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoadBalancersAPI) EXPECT() *MockLoadBalancersAPIMockRecorder {
	return m.recorder
}

// AddBackendNodes mocks base method.
func (m *MockLoadBalancersAPI) AddBackendNodes(arg0 context.Context, arg1 string, arg2 []string) ([]ah.LBBackendNode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBackendNodes", arg0, arg1, arg2)
	ret0, _ := ret[0].([]ah.LBBackendNode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBackendNodes indicates an expected call of AddBackendNodes.
func (mr *MockLoadBalancersAPIMockRecorder) AddBackendNodes(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBackendNodes", reflect.TypeOf((*MockLoadBalancersAPI)(nil).AddBackendNodes), arg0, arg1, arg2)
}

// AssignIPAddresses mocks base method.
func (m *MockLoadBalancersAPI) AssignIPAddresses(arg0 context.Context, arg1 string, arg2 []string) ([]ah.LBIPAddress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignIPAddresses", arg0, arg1, arg2)
	ret0, _ := ret[0].([]ah.LBIPAddress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignIPAddresses indicates an expected call of AssignIPAddresses.
func (mr *MockLoadBalancersAPIMockRecorder) AssignIPAddresses(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignIPAddresses", reflect.TypeOf((*MockLoadBalancersAPI)(nil).AssignIPAddresses), arg0, arg1, arg2)
}

// ConnectPrivateNetworks mocks base method.
func (m *MockLoadBalancersAPI) ConnectPrivateNetworks(arg0 context.Context, arg1 string, arg2 []string) ([]ah.LBPrivateNetwork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectPrivateNetworks", arg0, arg1, arg2)
	ret0, _ := ret[0].([]ah.LBPrivateNetwork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectPrivateNetworks indicates an expected call of ConnectPrivateNetworks.
func (mr *MockLoadBalancersAPIMockRecorder) ConnectPrivateNetworks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectPrivateNetworks", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ConnectPrivateNetworks), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockLoadBalancersAPI) Create(arg0 context.Context, arg1 *ah.LoadBalancerCreateRequest) (*ah.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*ah.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockLoadBalancersAPIMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLoadBalancersAPI)(nil).Create), arg0, arg1)
}

// CreateForwardingRule mocks base method.
func (m *MockLoadBalancersAPI) CreateForwardingRule(arg0 context.Context, arg1 string, arg2 *ah.LBForwardingRuleCreateRequest) (*ah.LBForwardingRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateForwardingRule", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBForwardingRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateForwardingRule indicates an expected call of CreateForwardingRule.
func (mr *MockLoadBalancersAPIMockRecorder) CreateForwardingRule(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateForwardingRule", reflect.TypeOf((*MockLoadBalancersAPI)(nil).CreateForwardingRule), arg0, arg1, arg2)
}

// CreateHealthCheck mocks base method.
func (m *MockLoadBalancersAPI) CreateHealthCheck(arg0 context.Context, arg1 string, arg2 *ah.LBHealthCheckCreateRequest) (*ah.LBHealthCheck, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHealthCheck", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBHealthCheck)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHealthCheck indicates an expected call of CreateHealthCheck.
func (mr *MockLoadBalancersAPIMockRecorder) CreateHealthCheck(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHealthCheck", reflect.TypeOf((*MockLoadBalancersAPI)(nil).CreateHealthCheck), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockLoadBalancersAPI) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockLoadBalancersAPIMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLoadBalancersAPI)(nil).Delete), arg0, arg1)
}

// DeleteBackendNode mocks base method.
func (m *MockLoadBalancersAPI) DeleteBackendNode(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBackendNode", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBackendNode indicates an expected call of DeleteBackendNode.
func (mr *MockLoadBalancersAPIMockRecorder) DeleteBackendNode(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBackendNode", reflect.TypeOf((*MockLoadBalancersAPI)(nil).DeleteBackendNode), arg0, arg1, arg2)
}

// DeleteForwardingRule mocks base method.
func (m *MockLoadBalancersAPI) DeleteForwardingRule(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteForwardingRule", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteForwardingRule indicates an expected call of DeleteForwardingRule.
func (mr *MockLoadBalancersAPIMockRecorder) DeleteForwardingRule(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteForwardingRule", reflect.TypeOf((*MockLoadBalancersAPI)(nil).DeleteForwardingRule), arg0, arg1, arg2)
}

// DeleteHealthCheck mocks base method.
func (m *MockLoadBalancersAPI) DeleteHealthCheck(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteHealthCheck", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHealthCheck indicates an expected call of DeleteHealthCheck.
func (mr *MockLoadBalancersAPIMockRecorder) DeleteHealthCheck(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHealthCheck", reflect.TypeOf((*MockLoadBalancersAPI)(nil).DeleteHealthCheck), arg0, arg1, arg2)
}

// DisconnectPrivateNetwork mocks base method.
func (m *MockLoadBalancersAPI) DisconnectPrivateNetwork(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisconnectPrivateNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisconnectPrivateNetwork indicates an expected call of DisconnectPrivateNetwork.
func (mr *MockLoadBalancersAPIMockRecorder) DisconnectPrivateNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisconnectPrivateNetwork", reflect.TypeOf((*MockLoadBalancersAPI)(nil).DisconnectPrivateNetwork), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockLoadBalancersAPI) Get(arg0 context.Context, arg1 string) (*ah.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*ah.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockLoadBalancersAPIMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLoadBalancersAPI)(nil).Get), arg0, arg1)
}

// GetBackendNode mocks base method.
func (m *MockLoadBalancersAPI) GetBackendNode(arg0 context.Context, arg1, arg2 string) (*ah.LBBackendNode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackendNode", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBBackendNode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBackendNode indicates an expected call of GetBackendNode.
func (mr *MockLoadBalancersAPIMockRecorder) GetBackendNode(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackendNode", reflect.TypeOf((*MockLoadBalancersAPI)(nil).GetBackendNode), arg0, arg1, arg2)
}

// GetForwardingRule mocks base method.
func (m *MockLoadBalancersAPI) GetForwardingRule(arg0 context.Context, arg1, arg2 string) (*ah.LBForwardingRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForwardingRule", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBForwardingRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForwardingRule indicates an expected call of GetForwardingRule.
func (mr *MockLoadBalancersAPIMockRecorder) GetForwardingRule(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForwardingRule", reflect.TypeOf((*MockLoadBalancersAPI)(nil).GetForwardingRule), arg0, arg1, arg2)
}

// GetHealthCheck mocks base method.
func (m *MockLoadBalancersAPI) GetHealthCheck(arg0 context.Context, arg1, arg2 string) (*ah.LBHealthCheck, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHealthCheck", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBHealthCheck)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHealthCheck indicates an expected call of GetHealthCheck.
func (mr *MockLoadBalancersAPIMockRecorder) GetHealthCheck(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHealthCheck", reflect.TypeOf((*MockLoadBalancersAPI)(nil).GetHealthCheck), arg0, arg1, arg2)
}

// GetIPAddress mocks base method.
func (m *MockLoadBalancersAPI) GetIPAddress(arg0 context.Context, arg1, arg2 string) (*ah.LBIPAddress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIPAddress", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBIPAddress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIPAddress indicates an expected call of GetIPAddress.
func (mr *MockLoadBalancersAPIMockRecorder) GetIPAddress(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIPAddress", reflect.TypeOf((*MockLoadBalancersAPI)(nil).GetIPAddress), arg0, arg1, arg2)
}

// GetPrivateNetwork mocks base method.
func (m *MockLoadBalancersAPI) GetPrivateNetwork(arg0 context.Context, arg1, arg2 string) (*ah.LBPrivateNetwork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivateNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(*ah.LBPrivateNetwork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrivateNetwork indicates an expected call of GetPrivateNetwork.
func (mr *MockLoadBalancersAPIMockRecorder) GetPrivateNetwork(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivateNetwork", reflect.TypeOf((*MockLoadBalancersAPI)(nil).GetPrivateNetwork), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockLoadBalancersAPI) List(arg0 context.Context) ([]ah.LoadBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]ah.LoadBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockLoadBalancersAPIMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockLoadBalancersAPI)(nil).List), arg0)
}

// ListBackendNodes mocks base method.
func (m *MockLoadBalancersAPI) ListBackendNodes(arg0 context.Context, arg1 string) ([]ah.LBBackendNode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBackendNodes", arg0, arg1)
	ret0, _ := ret[0].([]ah.LBBackendNode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBackendNodes indicates an expected call of ListBackendNodes.
func (mr *MockLoadBalancersAPIMockRecorder) ListBackendNodes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBackendNodes", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ListBackendNodes), arg0, arg1)
}

// ListForwardingRules mocks base method.
func (m *MockLoadBalancersAPI) ListForwardingRules(arg0 context.Context, arg1 string) ([]ah.LBForwardingRule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListForwardingRules", arg0, arg1)
	ret0, _ := ret[0].([]ah.LBForwardingRule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListForwardingRules indicates an expected call of ListForwardingRules.
func (mr *MockLoadBalancersAPIMockRecorder) ListForwardingRules(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListForwardingRules", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ListForwardingRules), arg0, arg1)
}

// ListHealthChecks mocks base method.
func (m *MockLoadBalancersAPI) ListHealthChecks(arg0 context.Context, arg1 string) ([]ah.LBHealthCheck, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListHealthChecks", arg0, arg1)
	ret0, _ := ret[0].([]ah.LBHealthCheck)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListHealthChecks indicates an expected call of ListHealthChecks.
func (mr *MockLoadBalancersAPIMockRecorder) ListHealthChecks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListHealthChecks", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ListHealthChecks), arg0, arg1)
}

// ListIPAddresses mocks base method.
func (m *MockLoadBalancersAPI) ListIPAddresses(arg0 context.Context, arg1 string) ([]ah.LBIPAddress, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIPAddresses", arg0, arg1)
	ret0, _ := ret[0].([]ah.LBIPAddress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIPAddresses indicates an expected call of ListIPAddresses.
func (mr *MockLoadBalancersAPIMockRecorder) ListIPAddresses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIPAddresses", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ListIPAddresses), arg0, arg1)
}

// ListPrivateNetworks mocks base method.
func (m *MockLoadBalancersAPI) ListPrivateNetworks(arg0 context.Context, arg1 string) ([]ah.LBPrivateNetwork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPrivateNetworks", arg0, arg1)
	ret0, _ := ret[0].([]ah.LBPrivateNetwork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPrivateNetworks indicates an expected call of ListPrivateNetworks.
func (mr *MockLoadBalancersAPIMockRecorder) ListPrivateNetworks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPrivateNetworks", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ListPrivateNetworks), arg0, arg1)
}

// ReleaseIPAddress mocks base method.
func (m *MockLoadBalancersAPI) ReleaseIPAddress(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseIPAddress", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReleaseIPAddress indicates an expected call of ReleaseIPAddress.
func (mr *MockLoadBalancersAPIMockRecorder) ReleaseIPAddress(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseIPAddress", reflect.TypeOf((*MockLoadBalancersAPI)(nil).ReleaseIPAddress), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockLoadBalancersAPI) Update(arg0 context.Context, arg1 string, arg2 *ah.LoadBalancerUpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockLoadBalancersAPIMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLoadBalancersAPI)(nil).Update), arg0, arg1, arg2)
}

// UpdateHealthCheck mocks base method.
func (m *MockLoadBalancersAPI) UpdateHealthCheck(arg0 context.Context, arg1, arg2 string, arg3 *ah.LBHealthCheckUpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateHealthCheck", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateHealthCheck indicates an expected call of UpdateHealthCheck.
func (mr *MockLoadBalancersAPIMockRecorder) UpdateHealthCheck(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateHealthCheck", reflect.TypeOf((*MockLoadBalancersAPI)(nil).UpdateHealthCheck), arg0, arg1, arg2, arg3)
}
