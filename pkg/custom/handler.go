package custom

import (
	"errors"
	"net/http"
	"strconv"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi/v7/domain/apiresponses"
	"github.com/vshn/crossplane-service-broker/pkg/crossplane"
	"github.com/vshn/crossplane-service-broker/pkg/reqcontext"
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

	sb, err := crossplane.ServiceBinderFactory(h.c, instance.Labels.ServiceName, instance.ID(), instance.ResourceRefs(), instance.Parameters(), h.logger)
	if err != nil {
		return nil, err
	}

	creds, err := sb.GetBinding(rctx.Context, instanceID)
	if err != nil {
		return nil, err
	}

	endpoints := []Endpoint{
		{
			Destination: creds["host"].(string),
			Ports:       strconv.Itoa(creds["port"].(int)),
			Protocol:    "tcp",
		},
	}
	if instance.Labels.ServiceName == crossplane.RedisService {
		sentinels := creds["sentinels"].([]crossplane.Credentials)
		for _, v := range sentinels {
			endpoints = append(endpoints, Endpoint{
				Destination: v["host"].(string),
				Ports:       strconv.Itoa(v["port"].(int)),
				Protocol:    "tcp",
			})
		}

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
