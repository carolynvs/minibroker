package user

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/carolynvs/osb-starter-pack/pkg/broker"
	"github.com/pkg/errors"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
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
	stableURL := "https://kubernetes-charts.storage.googleapis.com"

	home := helmpath.Home(environment.DefaultHelmHome)
	f, err := repo.LoadRepositoriesFile(home.RepositoryFile())
	if err != nil {
		return nil, err
	}

	cif := home.CacheIndex("stable")
	c := repo.Entry{
		Name:  "stable",
		Cache: cif,
		URL:   stableURL,
	}

	var settings environment.EnvSettings
	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return nil, err
	}

	if err := r.DownloadIndexFile(home.Cache()); err != nil {
		return nil, errors.Wrapf(err, "Looks like %q is not a valid chart repository or cannot be reached", stableURL)
	}

	f.Update(&c)
	f.WriteFile(home.RepositoryFile(), 0644)

	// Load the repositories.yaml
	rf, err := repo.LoadRepositoriesFile(home.RepositoryFile())
	if err != nil {
		return nil, err
	}

	catalog := &osb.CatalogResponse{}
	for _, re := range rf.Repositories {
		n := re.Name
		f := home.CacheIndex(n)
		index, err := repo.LoadIndexFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not load helm repository index at %s", f)
		}

		for name, chartVersions := range index.Entries {
			svc := osb.Service{
				ID:          name,
				Name:        name,
				Description: "Helm Chart for " + name,
				Bindable:    true,
				Plans:       make([]osb.Plan, 0, len(chartVersions)),
			}
			appVersions := map[string]*repo.ChartVersion{}
			for _, chartVersion := range chartVersions {
				if chartVersion.AppVersion == "" {
					continue
				}

				curV, err := semver.NewVersion(chartVersion.Version)
				if err != nil {
					fmt.Printf("Skipping %s@%s because %s is not a valid semver", name, chartVersion.AppVersion, chartVersion.Version)
					continue
				}

				currentMax, ok := appVersions[chartVersion.AppVersion]
				if !ok {
					appVersions[chartVersion.AppVersion] = chartVersion
				} else {
					maxV, _ := semver.NewVersion(currentMax.Version)
					if curV.GreaterThan(maxV) {
						appVersions[chartVersion.AppVersion] = chartVersion
					} else {
						fmt.Printf("Skipping %s@%s because %s<%s", name, chartVersion.AppVersion, curV, maxV)
						continue
					}
				}
			}

			for _, chartVersion := range appVersions {
				planToken := fmt.Sprintf("%s@%s", name, chartVersion.AppVersion)
				cleaner := regexp.MustCompile(`[^a-z0-9]`)
				planName := cleaner.ReplaceAllString(strings.ToLower(planToken), "-")
				plan := osb.Plan{
					ID:          planName,
					Name:        planName,
					Description: fmt.Sprintf("%s - %s", planToken, chartVersion.Description),
					Free:        boolPtr(true),
				}
				svc.Plans = append(svc.Plans, plan)
			}

			if len(svc.Plans) == 0 {
				continue
			}
			catalog.Services = append(catalog.Services, svc)
		}
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
