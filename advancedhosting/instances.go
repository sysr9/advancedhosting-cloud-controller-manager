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
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/advancedhosting/advancedhosting-api-go/ah"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

const ahProviderPrefix = "advancedhosting://"

var providerIDRegexp = regexp.MustCompile(fmt.Sprintf("%s(?P<instanceID>.*)", ahProviderPrefix))

type instances struct {
	client *ah.APIClient
}

func newInstances(client *ah.APIClient) *instances {
	return &instances{client: client}
}

// NodeAddresses returns the addresses of the specified instance.
func (i *instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	instance, err := i.instanceByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return i.instanceAddresses(instance)
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (i *instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	instance, err := i.instanceByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return i.instanceAddresses(instance)
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist, we must return ("", cloudprovider.InstanceNotFound)
// cloudprovider.InstanceNotFound should NOT be returned for instances that exist but are stopped/sleeping
func (i *instances) InstanceID(ctx context.Context, name types.NodeName) (string, error) {
	instance, err := i.instanceByName(ctx, name)
	if err != nil {
		return "", err
	}
	return instance.ID, nil
}

// InstanceType returns the type of the specified instance.
func (i *instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	instance, err := i.instanceByName(ctx, name)
	if err != nil {
		return "", err
	}
	return instance.Image.Slug, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (i *instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	instance, err := i.instanceByProviderID(ctx, providerID)
	if err != nil {
		return "", err
	}
	return instance.Image.Slug, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (i *instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (i *instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider exists.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
// This method should still return true for instances that exist but are stopped/sleeping.
func (i *instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	if _, err := i.instanceByProviderID(ctx, providerID); err != nil {
		if err == ah.ErrResourceNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (i *instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	instance, err := i.instanceByProviderID(ctx, providerID)
	if err != nil {
		return false, err
	}
	return instance.State == ah.InstanceShutDownStatus, nil
}

func (i *instances) instanceByName(ctx context.Context, nodeName types.NodeName) (*ah.Instance, error) {
	options := &ah.ListOptions{
		Filters: []ah.FilterInterface{
			&ah.EqFilter{
				Keys:  []string{"name"},
				Value: string(nodeName),
			},
		},
	}

	instances, _, err := i.client.Instances.List(ctx, options)
	if err != nil {
		return nil, err
	}

	if len(instances) != 1 {
		return nil, cloudprovider.InstanceNotFound
	}

	return &instances[0], nil
}

func (i *instances) instanceAddresses(instance *ah.Instance) ([]v1.NodeAddress, error) {
	var addresses []v1.NodeAddress
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: strings.ToLower(instance.Name)})

	publicIP, err := instance.PrimaryIPAddr()
	if err != nil {
		return nil, fmt.Errorf("Could not get public ip: %v", err)
	}
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeExternalIP, Address: publicIP.Address})

	privateIP := instance.PrivateNetworks[0].IP
	addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: privateIP})

	return addresses, nil
}

func (i *instances) instanceByProviderID(ctx context.Context, providerID string) (*ah.Instance, error) {
	instanceID, err := instanceIDByProviderID(providerID)
	if err != nil {
		return nil, err
	}

	instance, err := i.client.Instances.Get(ctx, instanceID)

	if err != nil {
		return nil, err
	}

	return instance, nil

}

func instanceIDByProviderID(providerID string) (string, error) {
	if providerID == "" {
		return "", errors.New("Empty ProviderID")
	}

	if !providerIDRegexp.MatchString(providerID) {
		return "", errors.New("Invalid ProviderID")
	}

	match := providerIDRegexp.FindStringSubmatch(providerID)

	return match[1], nil
}
