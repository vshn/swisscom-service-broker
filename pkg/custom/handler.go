package custom

import (
	"errors"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi/v7/domain/apiresponses"
	"github.com/vshn/crossplane-service-broker/pkg/crossplane"
	"github.com/vshn/crossplane-service-broker/pkg/reqcontext"

	xrv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

var errNotImplemented = apiresponses.NewFailureResponseBuilder(
	errors.New("not implemented"),
	http.StatusNotImplemented,
	"not-implemented").
	WithErrorKey("NotImplemented").
	Build()

// APIHandler handles the actual implementations and implements APISpec
type APIHandler struct {
	c      *crossplane.Crossplane
	logger lager.Logger
}

// NewAPIHandler sets up a new instance.
func NewAPIHandler(c *crossplane.Crossplane, logger lager.Logger) *APIHandler {
	return &APIHandler{c, logger}
}

// Endpoints retrieves the endpoints using the service binder.
func (h APIHandler) Endpoints(rctx *reqcontext.ReqContext, instanceID string) ([]Endpoint, error) {
	instance, _, exists, err := h.c.FindInstanceWithoutPlan(rctx, instanceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apiresponses.ErrInstanceDoesNotExist
	}

	// Get connection details of the actual Galera cluster
	if instance.Labels.ServiceName == crossplane.MariaDBDatabaseService {
		instance, err = h.getGaleraClusterFromDB(rctx, instance)
		if err != nil {
			return nil, err
		}
	}

	connectionDetails, err := h.c.GetConnectionDetails(rctx.Context, instance.Composite)
	if err != nil {
		return nil, err
	}

	dest := string(connectionDetails.Data[xrv1.ResourceCredentialsSecretEndpointKey])
	port := string(connectionDetails.Data[xrv1.ResourceCredentialsSecretPortKey])

	if len(dest) == 0 || len(port) == 0 {
		return nil, fmt.Errorf("instance %q is not yet ready", instanceID)
	}

	endpoints := []Endpoint{
		{
			Destination: dest,
			Ports:       port,
			Protocol:    "tcp",
		},
	}
	if instance.Labels.ServiceName == crossplane.RedisService {
		endpoints = append(endpoints, Endpoint{
			Destination: dest,
			Ports:       string(connectionDetails.Data[crossplane.SentinelPortKey]),
			Protocol:    "tcp",
		})
	}

	if p := string(connectionDetails.Data[crossplane.MetricsPortKey]); p != "" {
		endpoints = append(endpoints, Endpoint{
			Destination: dest,
			Ports:       p,
			Protocol:    "tcp",
		})
	}

	return endpoints, nil
}

func (h APIHandler) getGaleraClusterFromDB(rctx *reqcontext.ReqContext, db *crossplane.Instance) (*crossplane.Instance, error) {
	pRef, err := db.ParentReference()
	if err != nil {
		return nil, err
	}
	c, _, ok, err := h.c.FindInstanceWithoutPlan(rctx, pRef)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apiresponses.ErrInstanceDoesNotExist
	}
	return c, nil
}

// ServiceUsage is not implemented
func (h APIHandler) ServiceUsage(rctx *reqcontext.ReqContext, instanceID string) (*ServiceUsage, error) {
	return nil, errNotImplemented
}

// CreateUpdateServiceDefinition is not implemented
func (h APIHandler) CreateUpdateServiceDefinition(rctx *reqcontext.ReqContext, sd *ServiceDefinitionRequest) error {
	return errNotImplemented
}

// DeleteServiceDefinition is not implemented
func (h APIHandler) DeleteServiceDefinition(rctx *reqcontext.ReqContext, id string) error {
	return errNotImplemented
}

// CreateBackup is not implemented
func (h APIHandler) CreateBackup(rctx *reqcontext.ReqContext, instanceID string, b *BackupRequest) (*Backup, error) {
	return nil, errNotImplemented
}

// DeleteBackup is not implemented
func (h APIHandler) DeleteBackup(rctx *reqcontext.ReqContext, instanceID, backupID string) (string, error) {
	return "", errNotImplemented
}

// Backup is not implemented
func (h APIHandler) Backup(rctx *reqcontext.ReqContext, instanceID, backupID string) (*Backup, error) {
	return nil, errNotImplemented
}

// ListBackups is not implemented
func (h APIHandler) ListBackups(rctx *reqcontext.ReqContext, instanceID string) ([]Backup, error) {
	return nil, errNotImplemented
}

// RestoreBackup is not implemented
func (h APIHandler) RestoreBackup(rctx *reqcontext.ReqContext, instanceID, backupID string, r *RestoreRequest) (*Restore, error) {
	return nil, errNotImplemented
}

// RestoreStatus is not implemented
func (h APIHandler) RestoreStatus(rctx *reqcontext.ReqContext, instanceID, backupID, restoreID string) (*Restore, error) {
	return nil, errNotImplemented
}

// APIDocs is not implemented
func (h APIHandler) APIDocs(rctx *reqcontext.ReqContext, instanceID string) (string, error) {
	return "", errNotImplemented
}
