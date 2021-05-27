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
	"strconv"
	"strings"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

const (
	loadBalancerActiveStatus = "active"
)

const (
	// ServiceAnnotationLoadBalancerID is the ID of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerID = "service.beta.kubernetes.io/ah-loadbalancer-id"

	// ServiceAnnotationLoadBalancerName is the Name of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerName = "service.beta.kubernetes.io/ah-loadbalancer-name"

	// ServiceAnnotationLoadBalancerBalancingAlgorithm is the balancing algorithm of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerBalancingAlgorithm = "service.beta.kubernetes.io/ah-loadbalancer-balancing-algorithm"

	// ServiceAnnotationLoadBalancerEnableHealthCheck enables health check of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerEnableHealthCheck = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-enabled"

	// ServiceAnnotationLoadBalancerHealthCheckType is the health check type of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckType = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-type"

	// ServiceAnnotationLoadBalancerHealthCheckURL is the health check of url the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckURL = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-url"

	// ServiceAnnotationLoadBalancerHealthCheckInterval is the health check interval of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckInterval = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-interval"

	// ServiceAnnotationLoadBalancerHealthCheckTimeout is the health check timeout of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckTimeout = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-timeout"

	// ServiceAnnotationLoadBalancerHealthCheckUnhealthyThreshold is the health check unhealthy threshold of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckUnhealthyThreshold = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-unhealthy-threshold"

	// ServiceAnnotationLoadBalancerHealthCheckHealthyThreshold is the health check healthy threshold of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckHealthyThreshold = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-healthy-threshold"

	// ServiceAnnotationLoadBalancerHealthCheckPort is the health check port of the AH Managed Loadbalancer
	ServiceAnnotationLoadBalancerHealthCheckPort = "service.beta.kubernetes.io/ah-loadbalancer-healthcheck-port"
)

type loadbalancers struct {
	client      *ah.APIClient
	clusterInfo *clusterInfo
}

func newLoadbalancers(client *ah.APIClient, clusterInfo *clusterInfo) *loadbalancers {
	return &loadbalancers{client: client, clusterInfo: clusterInfo}
}

// GetLoadBalancer returns whether the specified load balancer exists, and
// if so, what its status is.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (l *loadbalancers) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (*v1.LoadBalancerStatus, bool, error) {

	lbID := l.loadBalancerID(service)
	if lbID == "" {
		return nil, false, nil
	}
	loadBalancer, err := l.loadBalancerByID(ctx, lbID)

	if err != nil {
		if err == ah.ErrResourceNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return l.loadBalancerStatus(loadBalancer), true, nil
}

// GetLoadBalancerName returns the name of the load balancer. Implementations must treat the
// *v1.Service parameter as read-only and not modify it.
func (l *loadbalancers) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	return l.loadBalancerName(service)
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (l *loadbalancers) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {

	lbID := l.loadBalancerID(service)

	var loadBalancer *ah.LoadBalancer

	loadBalancer, err := l.loadBalancerByID(ctx, lbID)

	switch err {
	case ah.ErrResourceNotFound:
		loadBalancer, err = l.createLoadBalancer(ctx, service, nodes)
		if err != nil {
			return nil, fmt.Errorf("Error creating load balancer: %v", err)
		}
	case nil:

	default:
		return nil, err
	}

	if loadBalancer.State != loadBalancerActiveStatus {
		return nil, fmt.Errorf("Load balancer is not active yet: %s", loadBalancer.State)
	}

	if err = l.updateLoadBalancer(ctx, service, nodes, loadBalancer); err != nil {
		return nil, fmt.Errorf("Error updating load balancer: %v", err)
	}

	return l.loadBalancerStatus(loadBalancer), nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (l *loadbalancers) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	lbID := l.loadBalancerID(service)

	var loadBalancer *ah.LoadBalancer

	loadBalancer, err := l.loadBalancerByID(ctx, lbID)

	if err != nil {
		return err
	}

	if loadBalancer.State != loadBalancerActiveStatus {
		return fmt.Errorf("Load balancer is not active yet: %s", loadBalancer.State)
	}

	return l.updateLoadBalancer(ctx, service, nodes, loadBalancer)

}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (l *loadbalancers) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	lbID := l.loadBalancerID(service)

	var loadBalancer *ah.LoadBalancer

	loadBalancer, err := l.loadBalancerByID(ctx, lbID)

	switch err {
	case ah.ErrResourceNotFound:
		return nil
	case nil:
		break
	default:
		return err
	}

	if loadBalancer.State == "deleting" {
		return fmt.Errorf("Load balancer is already in deletion state")
	}

	if err = l.client.LoadBalancers.Delete(ctx, loadBalancer.ID); err != nil {
		if err == ah.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error deleting load balancer: %v", err)
	}

	return fmt.Errorf("LB deletion has been started")

}

func (l *loadbalancers) loadBalancerID(service *v1.Service) string {
	return service.ObjectMeta.Annotations[ServiceAnnotationLoadBalancerID]
}

func (l *loadbalancers) loadBalancerByID(ctx context.Context, lbID string) (*ah.LoadBalancer, error) {
	if lbID == "" {
		return nil, ah.ErrResourceNotFound
	}
	loadBalancer, err := l.client.LoadBalancers.Get(ctx, lbID)
	if err != nil {
		return nil, err
	}
	return loadBalancer, nil
}

func (l *loadbalancers) loadBalancerStatus(loadBalancer *ah.LoadBalancer) *v1.LoadBalancerStatus {
	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: loadBalancer.IPAddresses[0].Address,
			},
		},
	}
}

func (l *loadbalancers) createLoadBalancer(ctx context.Context, service *v1.Service, nodes []*v1.Node) (*ah.LoadBalancer, error) {

	request, err := l.makeLoadBalancerCreateRequest(ctx, service, nodes)
	if err != nil {
		return nil, fmt.Errorf("Error makeLoadBalancerCreateRequest: %v", err)
	}

	loadBalancer, err := l.client.LoadBalancers.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("API LoadBalancers.Create error: %v", err)
	}

	patcher := newServicePatcher(l.clusterInfo.kclient, service)
	annotateService(service, ServiceAnnotationLoadBalancerID, loadBalancer.ID)
	if err = patcher.Patch(ctx); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (l *loadbalancers) makeLoadBalancerCreateRequest(ctx context.Context, service *v1.Service, nodes []*v1.Node) (*ah.LoadBalancerCreateRequest, error) {
	request := &ah.LoadBalancerCreateRequest{
		Name:                  l.loadBalancerName(service),
		DatacenterID:          l.clusterInfo.DatacenterID,
		CreatePublicIPAddress: true,
		PrivateNetworkIDs:     []string{l.clusterInfo.PrivateNetworkID},
		BalancingAlgorithm:    l.loadBalancerBalancingAlgorithm(service),
		ForwardingRules:       l.loadBalancerForwardingRules(service),
	}

	if l.loadBalancerHealthChecksEnabled(service) {
		healthCheck, err := l.loadBalancerHealthCheckRequest(service)
		if err != nil {
			return nil, err
		}
		request.HealthChecks = []ah.LBHealthCheckCreateRequest{*healthCheck}
	}

	backendNodes, err := l.loadBalancerBackendNodes(nodes)
	if err != nil {
		return nil, err
	}
	request.BackendNodes = backendNodes

	return request, nil
}

func (l *loadbalancers) loadBalancerName(service *v1.Service) string {
	if name, ok := service.Annotations[ServiceAnnotationLoadBalancerName]; ok {
		return name
	}
	return cloudprovider.DefaultLoadBalancerName(service)
}

func (l *loadbalancers) loadBalancerBalancingAlgorithm(service *v1.Service) string {
	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerBalancingAlgorithm]; ok {
		return v
	}
	return "round_robin"
}

func (l *loadbalancers) loadBalancerForwardingRules(service *v1.Service) []ah.LBForwardingRuleCreateRequest {
	requests := make([]ah.LBForwardingRuleCreateRequest, len(service.Spec.Ports))

	for i, port := range service.Spec.Ports {
		requests[i] = l.lbForwardingRuleCreateRequest(port)
	}

	return requests
}

func (l *loadbalancers) lbForwardingRuleCreateRequest(port v1.ServicePort) ah.LBForwardingRuleCreateRequest {

	protocol := strings.ToLower(string(port.Protocol))
	return ah.LBForwardingRuleCreateRequest{
		RequestProtocol:       protocol,
		RequestPort:           int(port.Port),
		CommunicationProtocol: protocol,
		CommunicationPort:     int(port.NodePort),
	}
}

func (l *loadbalancers) loadBalancerHealthChecksEnabled(service *v1.Service) bool {
	v, ok := service.Annotations[ServiceAnnotationLoadBalancerEnableHealthCheck]
	if !ok {
		return false
	}

	res, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return res
}

func (l *loadbalancers) loadBalancerHealthCheckRequest(service *v1.Service) (*ah.LBHealthCheckCreateRequest, error) {
	var request ah.LBHealthCheckCreateRequest

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckType]; ok {
		request.Type = v
	}

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckURL]; ok {
		request.URL = v
	}

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckInterval]; ok {
		interval, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Invalid health check interval: %v", err)
		}
		request.Interval = interval
	}

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckTimeout]; ok {
		timeout, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Invalid health check timeout: %v", err)
		}
		request.Timeout = timeout
	}

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckUnhealthyThreshold]; ok {
		unhealthyThreshold, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Invalid health check unhealthy threshold: %v", err)
		}
		request.UnhealthyThreshold = unhealthyThreshold
	}

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckHealthyThreshold]; ok {
		healthyThreshold, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Invalid health check healthy threshold: %v", err)
		}
		request.HealthyThreshold = healthyThreshold
	}

	if v, ok := service.Annotations[ServiceAnnotationLoadBalancerHealthCheckPort]; ok {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("Invalid health check port: %v", err)
		}
		request.Port = port
	}

	return &request, nil
}

func (l *loadbalancers) loadBalancerBackendNodes(nodes []*v1.Node) ([]ah.LBBackendNodeCreateRequest, error) {
	var requests []ah.LBBackendNodeCreateRequest
	for _, node := range nodes {

		id, err := instanceIDByProviderID(node.Spec.ProviderID)

		if err != nil {
			return nil, err
		}

		request := ah.LBBackendNodeCreateRequest{
			CloudServerID: id,
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (l *loadbalancers) updateLoadBalancer(ctx context.Context, service *v1.Service, nodes []*v1.Node, lb *ah.LoadBalancer) error {
	if err := l.updateLoadBalancerInfo(ctx, service, lb); err != nil {
		return err
	}

	if err := l.updateForwardingRules(ctx, service, lb); err != nil {
		return err
	}

	if l.loadBalancerHealthChecksEnabled(service) {
		if err := l.updateHealthChecks(ctx, service, lb); err != nil {
			return err
		}
	} else {
		if err := l.deleteHealthChecks(ctx, lb); err != nil {
			return err
		}
	}

	if err := l.updateBackendNodes(ctx, nodes, lb); err != nil {
		return err
	}

	return nil
}

func (l *loadbalancers) updateLoadBalancerInfo(ctx context.Context, service *v1.Service, lb *ah.LoadBalancer) error {
	var request ah.LoadBalancerUpdateRequest
	var updated bool
	if l.loadBalancerName(service) != lb.Name {
		request.Name = l.loadBalancerName(service)
		updated = true
	}

	if l.loadBalancerBalancingAlgorithm(service) != lb.BalancingAlgorithm {
		request.BalancingAlgorithm = l.loadBalancerBalancingAlgorithm(service)
		updated = true
	}

	if !updated {
		return nil
	}

	if err := l.client.LoadBalancers.Update(ctx, lb.ID, &request); err != nil {
		return err
	}

	stateFunc := func(ctx context.Context) (state string, err error) {
		lb, err := l.client.LoadBalancers.Get(ctx, lb.ID)
		if err != nil {
			return "", err
		}
		return lb.State, nil
	}

	if err := waitForState(ctx, stateFunc, "active"); err != nil {
		return err
	}

	return nil

}

func (l *loadbalancers) updateForwardingRules(ctx context.Context, service *v1.Service, lb *ah.LoadBalancer) error {
	frsToDelete := make(map[int]ah.LBForwardingRule, len(lb.ForwardingRules))

	for _, fr := range lb.ForwardingRules {
		frsToDelete[fr.RequestPort] = fr
	}

	for _, port := range service.Spec.Ports {
		fr, ok := frsToDelete[int(port.Port)]
		if !ok {
			if err := l.addForwardingRule(ctx, lb.ID, &port); err != nil {
				return err
			}
		} else {
			if err := l.updateForwardingRule(ctx, &port, &fr, lb.ID); err != nil {
				return err
			}
			delete(frsToDelete, int(port.Port))
		}

	}

	for _, fr := range frsToDelete {
		if err := l.removeForwardingRule(ctx, lb.ID, fr.ID); err != nil {
			return err
		}
	}

	return nil
}

func (l *loadbalancers) addForwardingRule(ctx context.Context, lbID string, port *v1.ServicePort) error {
	request := l.lbForwardingRuleCreateRequest(*port)
	fr, err := l.client.LoadBalancers.CreateForwardingRule(ctx, lbID, &request)
	if err != nil {
		return err
	}

	stateFunc := func(ctx context.Context) (state string, err error) {
		fr, err := l.client.LoadBalancers.GetForwardingRule(ctx, lbID, fr.ID)
		if err != nil {
			return "", err
		}
		return fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "active"); err != nil {
		return err
	}

	return nil
}

func (l *loadbalancers) removeForwardingRule(ctx context.Context, lbID, frID string) error {
	if err := l.client.LoadBalancers.DeleteForwardingRule(ctx, lbID, frID); err != nil {
		return err
	}

	stateFunc := func(ctx context.Context) (state string, err error) {
		fr, err := l.client.LoadBalancers.GetForwardingRule(ctx, lbID, frID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return "deleted", nil
			}
			return "", err
		}
		return fr.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleted"); err != nil {
		return err
	}

	return nil
}

func (l *loadbalancers) updateForwardingRule(ctx context.Context, port *v1.ServicePort, fr *ah.LBForwardingRule, lbID string) error {
	protocol := strings.ToLower(string(port.Protocol))
	if fr.RequestProtocol != protocol ||
		fr.RequestPort != int(port.Port) ||
		fr.CommunicationProtocol != protocol ||
		fr.CommunicationPort != int(port.NodePort) {

		if err := l.removeForwardingRule(ctx, lbID, fr.ID); err != nil {
			return err
		}

		if err := l.addForwardingRule(ctx, lbID, port); err != nil {
			return err
		}
	}
	return nil
}

func (l *loadbalancers) updateHealthChecks(ctx context.Context, service *v1.Service, lb *ah.LoadBalancer) error {
	hc, err := l.loadBalancerHealthCheckRequest(service)
	if err != nil {
		return err
	}

	if len(lb.HealthChecks) == 0 {
		healthCheck, err := l.client.LoadBalancers.CreateHealthCheck(ctx, lb.ID, hc)

		if err != nil {
			return err
		}

		stateFunc := func(ctx context.Context) (state string, err error) {
			hc, err := l.client.LoadBalancers.GetHealthCheck(ctx, lb.ID, healthCheck.ID)
			if err != nil {
				return "", err
			}
			return hc.State, nil
		}

		if err := waitForState(ctx, stateFunc, "active"); err != nil {
			return err
		}

		return nil

	}

	origHC := lb.HealthChecks[0]

	if origHC.Type != hc.Type ||
		origHC.URL != hc.URL ||
		origHC.Interval != hc.Interval ||
		origHC.Timeout != hc.Timeout ||
		origHC.UnhealthyThreshold != hc.UnhealthyThreshold ||
		origHC.HealthyThreshold != hc.HealthyThreshold ||
		origHC.Port != hc.Port {

		request := &ah.LBHealthCheckUpdateRequest{
			Type:               hc.Type,
			URL:                hc.URL,
			Interval:           hc.Interval,
			Timeout:            hc.Timeout,
			UnhealthyThreshold: hc.UnhealthyThreshold,
			HealthyThreshold:   hc.HealthyThreshold,
			Port:               hc.Port,
		}

		if err = l.client.LoadBalancers.UpdateHealthCheck(ctx, lb.ID, origHC.ID, request); err != nil {
			return err
		}

		stateFunc := func(ctx context.Context) (state string, err error) {
			hc, err := l.client.LoadBalancers.GetHealthCheck(ctx, lb.ID, origHC.ID)
			if err != nil {
				return "", err
			}
			return hc.State, nil
		}

		if err := waitForState(ctx, stateFunc, "active"); err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (l *loadbalancers) deleteHealthChecks(ctx context.Context, lb *ah.LoadBalancer) error {
	if len(lb.HealthChecks) == 0 {
		return nil
	}

	origHC := lb.HealthChecks[0]

	if err := l.client.LoadBalancers.DeleteHealthCheck(ctx, lb.ID, origHC.ID); err != nil {
		return err
	}

	stateFunc := func(ctx context.Context) (state string, err error) {
		hc, err := l.client.LoadBalancers.GetHealthCheck(ctx, lb.ID, origHC.ID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return "deleted", nil
			}
			return "", err
		}
		return hc.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleted"); err != nil {
		return err
	}

	return nil

}

func (l *loadbalancers) updateBackendNodes(ctx context.Context, nodes []*v1.Node, lb *ah.LoadBalancer) error {
	bnsToDelete := make(map[string]ah.LBBackendNode, len(lb.BackendNodes))
	var bnsToAdd []string

	for _, bn := range lb.BackendNodes {
		bnsToDelete[bn.CloudServerID] = bn
	}

	for _, node := range nodes {

		id, err := instanceIDByProviderID(node.Spec.ProviderID)

		if err != nil {
			return err
		}

		if _, ok := bnsToDelete[id]; !ok {
			bnsToAdd = append(bnsToAdd, id)
		} else {
			delete(bnsToDelete, id)
		}

	}

	if len(bnsToAdd) > 0 {
		if err := l.addBackendNodes(ctx, lb.ID, bnsToAdd); err != nil {
			return err
		}
	}

	for _, bn := range bnsToDelete {
		if err := l.removeBackendNode(ctx, lb.ID, bn.ID); err != nil {
			return err
		}
	}

	return nil
}

func (l *loadbalancers) addBackendNodes(ctx context.Context, lbID string, bns []string) error {
	backendNodes, err := l.client.LoadBalancers.AddBackendNodes(ctx, lbID, bns)
	if err != nil {
		return err
	}

	stateFunc := func(ctx context.Context) (state string, err error) {
		bns, err := l.client.LoadBalancers.ListBackendNodes(ctx, lbID)
		if err != nil {
			return "", err
		}

		backendMap := make(map[string]string, len(bns))
		for _, bn := range bns {
			backendMap[bn.ID] = bn.State
		}
		for _, bn := range backendNodes {
			curState := backendMap[bn.ID]
			if curState != "active" {
				return curState, nil
			}
		}

		return "active", nil
	}

	if err := waitForState(ctx, stateFunc, "active"); err != nil {
		return err
	}

	return nil

}

func (l *loadbalancers) removeBackendNode(ctx context.Context, lbID, bnID string) error {
	if err := l.client.LoadBalancers.DeleteBackendNode(ctx, lbID, bnID); err != nil {
		return err
	}

	stateFunc := func(ctx context.Context) (state string, err error) {
		bn, err := l.client.LoadBalancers.GetBackendNode(ctx, lbID, bnID)
		if err != nil {
			if err == ah.ErrResourceNotFound {
				return "deleted", nil
			}
			return "", err
		}
		return bn.State, nil
	}

	if err := waitForState(ctx, stateFunc, "deleted"); err != nil {
		return err
	}

	return nil

}
