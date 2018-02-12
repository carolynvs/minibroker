package user

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/carolynvs/osb-starter-pack/pkg/broker"
	"github.com/pkg/errors"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"k8s.io/helm/pkg/repo"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.
	return &BusinessLogic{
		instances: make(map[string]*exampleInstance, 10),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Add fields here! These fields are provided purely as an example
	sync.RWMutex
	instances map[string]*exampleInstance
}

var _ broker.BusinessLogic = &BusinessLogic{}

func (b *BusinessLogic) GetCatalog(response http.ResponseWriter, request *http.Request) (*osb.CatalogResponse, error) {
	repoURL := "https://kubernetes-charts.storage.googleapis.com"
	dlrequest, err := http.Get(repoURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not download upstream repo index at %s", repoURL)
	}
	body, err := ioutil.ReadAll(dlrequest.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read %s body", repoURL)
	}

	if dlrequest.StatusCode != http.StatusOK {
		return nil, errors.Errorf("GET %s (%v)\n%s", dlrequest.StatusCode)
	}
	indexFile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, errors.Wrap(err, "Could not create temp file for the index")
	}

	err = ioutil.WriteFile(indexFile.Name(), body, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not write temp file %s", indexFile.Name())
	}
	index, err := repo.LoadIndexFile(indexFile.Name())
	if err != nil {
		return nil, errors.Wrapf(err, "Could not load helm repository index at %s", indexFile.Name())
	}

	catalog := &osb.CatalogResponse{
		Services: make([]osb.Service, 0, len(index.Entries)),
	}

	for name, versions := range index.Entries {
		svc := osb.Service{
			ID:          name,
			Name:        name,
			Description: "Helm Chart for " + name,
			Bindable:    true,
			Plans:       make([]osb.Plan, 0, len(versions)),
		}
		for _, version := range versions {
			planName := fmt.Sprintf("%s@%s", name, version.AppVersion)
			plan := osb.Plan{
				ID:          planName,
				Name:        planName,
				Description: version.Description,
				Free:        boolPtr(true),
			}
			svc.Plans = append(svc.Plans, plan)
		}
		catalog.Services = append(catalog.Services, svc)
	}

	return catalog, nil
}

func (b *BusinessLogic) Provision(pr *osb.ProvisionRequest, w http.ResponseWriter, r *http.Request) (*osb.ProvisionResponse, error) {
	// Your provision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	exampleInstance := &exampleInstance{ID: pr.InstanceID, Params: pr.Parameters}
	b.instances[pr.InstanceID] = exampleInstance

	return &osb.ProvisionResponse{}, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, w http.ResponseWriter, r *http.Request) (*osb.DeprovisionResponse, error) {
	// Your deprovision business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	delete(b.instances, request.InstanceID)

	return &osb.DeprovisionResponse{}, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, w http.ResponseWriter, r *http.Request) (*osb.LastOperationResponse, error) {
	// Your last-operation business logic goes here

	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, w http.ResponseWriter, r *http.Request) (*osb.BindResponse, error) {
	// Your bind business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	return &osb.BindResponse{Credentials: instance.Params}, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, w http.ResponseWriter, r *http.Request) (*osb.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &osb.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, w http.ResponseWriter, r *http.Request) (*osb.UpdateInstanceResponse, error) {
	// Your logic for updating a service goes here.
	return &osb.UpdateInstanceResponse{}, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}

// example types

// exampleInstance is intended as an example of a type that holds information about a service instance
type exampleInstance struct {
	ID     string
	Params map[string]interface{}
}

func boolPtr(value bool) *bool {
	return &value
}
