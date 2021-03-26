package custom

import (
	"errors"
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
		parentRef, err := instance.ParentReference()
		if err != nil {
			return nil, err
		}
		instance, _, exists, err = h.c.FindInstanceWithoutPlan(rctx, parentRef)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, apiresponses.ErrInstanceDoesNotExist
		}
	}

	connectionDetails, err := h.c.GetConnectionDetails(rctx.Context, instance.Composite)
	if err != nil {
		return nil, err
	}

	dest := string(connectionDetails.Data[xrv1.ResourceCredentialsSecretEndpointKey])

	endpoints := []Endpoint{
		{
			Destination: dest,
			Ports:       string(connectionDetails.Data[xrv1.ResourceCredentialsSecretPortKey]),
			Protocol:    "tcp",
		},
	}
	if instance.Labels.ServiceName == crossplane.RedisService {
		endpoints = append(endpoints, Endpoint{
			Destination: dest,
			Ports:       string(connectionDetails.Data["sentinelPort"]),
			Protocol:    "tcp",
		})
	}
	return endpoints, nil
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
