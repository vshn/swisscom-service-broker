package custom

import (
	"github.com/vshn/crossplane-service-broker/pkg/reqcontext"
)

// APISpec describes the service broker endpoints not defined by the open service broker API spec.
type APISpec interface {
	// Endpoints lists service endpoints
	// GET /custom/service_instances/{service_instance_id}/endpoint
	Endpoints(rctx *reqcontext.ReqContext, instanceID string) ([]Endpoint, error)
	// ServiceUsage returns service usage
	// GET /custom/service_instances/{service_instance_id}/usage
	ServiceUsage(rctx *reqcontext.ReqContext, instanceID string) (*ServiceUsage, error)
	// CreateUpdateServiceDefinition is not implemented
	// POST /custom/admin/service-definition
	CreateUpdateServiceDefinition(rctx *reqcontext.ReqContext, sd *ServiceDefinitionRequest) error
	// DeleteServiceDefinition is not implemented
	// DELETE /custom/admin/service-definition/{id}
	DeleteServiceDefinition(rctx *reqcontext.ReqContext, id string) error
	// CreateBackup is not implemented
	// POST /custom/service_instances/{service_instance_id}/backups
	CreateBackup(rctx *reqcontext.ReqContext, instanceID string, b *BackupRequest) (*Backup, error)
	// DeleteBackup is not implemented
	// DELETE /custom/service_instances/{service_instance_id}/backups/{backup_id}
	DeleteBackup(rctx *reqcontext.ReqContext, instanceID, backupID string) (string, error)
	// Backup is not implemented
	// GET /custom/service_instances/{service_instance_id}/backups/{backup_id}
	Backup(rctx *reqcontext.ReqContext, instanceID, backupID string) (*Backup, error)
	// ListBackups is not implemented
	// GET /custom/service_instances/{service_instance_id}/backups
	ListBackups(rctx *reqcontext.ReqContext, instanceID string) ([]Backup, error)
	// RestoreBackup is not implemented
	// POST /custom/service_instances/{service_instance_id}/backups/{backup_id}/restores
	RestoreBackup(rctx *reqcontext.ReqContext, instanceID, backupID string, r *RestoreRequest) (*Restore, error)
	// RestoreStatus is not implemented
	// GET /custom/service_instances/{service_instance_id}/backups/{backup_id}/restores/{restore_id}
	RestoreStatus(rctx *reqcontext.ReqContext, instanceID, backupID, restoreID string) (*Restore, error)
	// APIDocs is not implemented
	// GET /custom/service_instances/{service_instance_id}/api-docs
	APIDocs(rctx *reqcontext.ReqContext, instanceID string) (string, error)
}

// Endpoint describes available service endpoints.
type Endpoint struct {
	Destination string `json:"destination"`
	Ports       string `json:"ports"`
	Protocol    string `json:"protocol"`
}

// ServiceUsage is a placeholder
type ServiceUsage struct{}

// Backup is a placeholder
type Backup struct{}

// Restore is a placeholder
type Restore struct{}

// ServiceDefinitionRequest is a placeholder
type ServiceDefinitionRequest struct{}

// BackupRequest is a placeholder
type BackupRequest struct{}

// RestoreRequest is a placeholder
type RestoreRequest struct{}
