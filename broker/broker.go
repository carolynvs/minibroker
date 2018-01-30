package broker

import (
	"context"

	"github.com/pivotal-cf/brokerapi"
)

type InClusterBroker struct{}

func (InClusterBroker) Services(ctx context.Context) []brokerapi.Service {
	panic("implement me")
}

func (InClusterBroker) Provision(ctx context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	panic("implement me")
}

func (InClusterBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	panic("implement me")
}

func (InClusterBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	panic("implement me")
}

func (InClusterBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	panic("implement me")
}

func (InClusterBroker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	panic("implement me")
}

func (InClusterBroker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	panic("implement me")
}
