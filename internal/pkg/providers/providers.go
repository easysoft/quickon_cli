// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package providers

import (
	"fmt"
	"sync"

	"github.com/easysoft/qcadmin/internal/pkg/types"
)

// Factory is a function that returns a Provider.Interface.
type Factory func() (Provider, error)

var (
	providersMutex sync.Mutex
	providers      = make(map[string]Factory)
)

type Provider interface {
	GetProviderName() string
	// Get command usage example.
	GetUsageExample(action string) string
	CreateCluster() error
	JoinNode() error
	InitQucheng() error
	GetCreateFlags() []types.Flag
	GetJoinFlags() []types.Flag
	CreateCheck(skip bool) error
	PreSystemInit() error
	Show()
}

// RegisterProvider registers a provider.Factory by name.
func RegisterProvider(name string, p Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := providers[name]; !found {
		// log.Flog.Infof("registered provider %s", name)
		providers[name] = p
	}
}

// GetProvider creates an instance of the named provider, or nil if
// the name is unknown.  The error return is only used if the named provider
// was known but failed to initialize.
func GetProvider(name string) (Provider, error) {
	if name == "" {
		name = "native"
	}
	providersMutex.Lock()
	defer providersMutex.Unlock()
	f, found := providers[name]
	if !found {
		return nil, fmt.Errorf("provider %s is not registered", name)
	}
	return f()
}

// ListProviders returns current supported providers.
func ListProviders() []types.Provider {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	list := make([]types.Provider, 0)
	for p := range providers {
		list = append(list, types.Provider{
			Name: p,
		})
	}
	return list
}
