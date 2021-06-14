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
	"io"
	"os"

	"github.com/advancedhosting/advancedhosting-api-go/ah"

	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

const (
	providerName            = "advancedhosting"
	ahAPIToken              = "AH_API_TOKEN"
	ahAPIBaseURL            = "AH_API_URL"
	ahClusterPrivateNetwork = "AH_CLUSTER_PRIVATE_NETWORK_NUMBER"
	ahClusterDatacenter     = "AH_CLUSTER_DATACENTER"
)

type cloud struct {
	client        *ah.APIClient
	instances     cloudprovider.Instances
	zones         cloudprovider.Zones
	loadbalancers cloudprovider.LoadBalancer
	clusterInfo   *clusterInfo
}

type clusterInfo struct {
	PrivateNetworkID string
	DatacenterID     string
	kclient          kubernetes.Interface
}

func newCloud() (cloudprovider.Interface, error) {

	token := os.Getenv(ahAPIToken)
	if token == "" {
		return nil, fmt.Errorf("AdvancedHosting API token is required")
	}

	baseURL := os.Getenv(ahAPIBaseURL)
	if baseURL == "" {
		baseURL = "https://api.websa.com"
	}

	clientOptions := &ah.ClientOptions{
		Token:   token,
		BaseURL: baseURL,
	}

	client, err := ah.NewAPIClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating Api Client: %s", err)
	}

	clusterInfo, err := newClusterInfo(client)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating clusterInfo: %s", err)
	}

	return &cloud{
		client:        client,
		clusterInfo:   clusterInfo,
		instances:     newInstances(client),
		loadbalancers: newLoadbalancers(client, clusterInfo),
	}, nil
}

func newClusterInfo(client *ah.APIClient) (*clusterInfo, error) {
	pnNumber := os.Getenv(ahClusterPrivateNetwork)
	if pnNumber == "" {
		return nil, fmt.Errorf("private Network Number is required")
	}
	pnID, err := privateNetworkIDbyNumber(pnNumber, client)
	if err != nil {
		return nil, fmt.Errorf("error getting pnID: %v", err)
	}
	datacenterSlug := os.Getenv(ahClusterDatacenter)
	if datacenterSlug == "" {
		return nil, fmt.Errorf("datacenter ID is required")
	}

	datacenterID, err := datacenterIDBySlug(datacenterSlug, client)
	if err != nil {
		return nil, fmt.Errorf("error getting datacenterID: %v", err)
	}

	return &clusterInfo{PrivateNetworkID: pnID, DatacenterID: datacenterID}, nil
}

func privateNetworkIDbyNumber(pnNumber string, client *ah.APIClient) (string, error) {
	options := &ah.ListOptions{
		Filters: []ah.FilterInterface{&ah.EqFilter{Keys: []string{"number"}, Value: pnNumber}},
	}
	privateNetworks, err := client.PrivateNetworks.List(context.Background(), options)
	if err != nil {
		return "", err
	}
	if len(privateNetworks) != 1 {
		return "", ah.ErrResourceNotFound
	}
	return privateNetworks[0].ID, nil
}

func datacenterIDBySlug(datacenterSlug string, client *ah.APIClient) (string, error) {
	datacenters, err := client.Datacenters.List(context.Background(), nil)
	if err != nil {
		return "", err
	}

	for _, datacenter := range datacenters {
		if datacenter.Slug == datacenterSlug {
			return datacenter.ID, nil
		}
	}
	return "", ah.ErrResourceNotFound
}

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(io.Reader) (cloudprovider.Interface, error) {
		return newCloud()
	})
}

func (c *cloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {

	c.clusterInfo.kclient = clientBuilder.ClientOrDie("advancedhosting-cloud-controller-manager")

	klog.Infof("clientset initialized")

}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return c.loadbalancers, true
}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	return c.instances, true
}

func (c *cloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return nil, false
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	return c.zones, false
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (c *cloud) ProviderName() string {
	return providerName
}

func (c *cloud) HasClusterID() bool {
	return false
}
