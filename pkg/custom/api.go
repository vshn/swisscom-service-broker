package custom

import (
	"encoding/json"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/pivotal-cf/brokerapi/v7/auth"
	"github.com/pivotal-cf/brokerapi/v7/domain/apiresponses"
	"github.com/pivotal-cf/brokerapi/v7/middlewares"
	"github.com/vshn/crossplane-service-broker/pkg/api"
	"github.com/vshn/crossplane-service-broker/pkg/reqcontext"
)

// API exposes the custom api handlers specific for this broker.
type API struct {
	handler APISpec
	logger  lager.Logger
}

// NewAPI registers the routes and middlewares.
func NewAPI(router *mux.Router, handler APISpec, username, password string, logger lager.Logger) *API {
	a := API{
		handler: handler,
		logger:  logger,
	}

	attachRoutes(router, a)

	authMiddleware := auth.NewWrapper(username, password).Wrap

	router.Use(middlewares.AddCorrelationIDToContext)
	router.Use(authMiddleware)
	router.Use(middlewares.AddOriginatingIdentityToContext)
	router.Use(middlewares.AddInfoLocationToContext)
	router.Use(api.LoggerMiddleware(logger))

	return &a
}

func attachRoutes(router *mux.Router, api API) {
	router.HandleFunc("/custom/service_instances/{service_instance_id}/endpoint", api.Endpoints).Methods("GET")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/usage", api.ServiceUsage).Methods("GET")
	router.HandleFunc("/custom/admin/service-definition", api.CreateUpdateServiceDefinition).Methods("POST")
	router.HandleFunc("/custom/admin/service-definition/{id}", api.DeleteServiceDefinition).Methods("DELETE")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/backups", api.CreateBackup).Methods("POST")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/backups/{backup_id}", api.DeleteBackup).Methods("DELETE")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/backups/{backup_id}", api.Backup).Methods("GET")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/backups", api.ListBackups).Methods("GET")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/backups/{backup_id}/restores", api.RestoreBackup).Methods("POST")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/backups/{backup_id}/restores/{restore_id}", api.Endpoints).Methods("GET")
	router.HandleFunc("/custom/service_instances/{service_instance_id}/api-docs", api.APIDocs).Methods("GET")
}

func (a API) respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if response == nil {
		return
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		a.logger.Error("encoding response", err, lager.Data{"status": status, "response": response})
	}
}

func (a API) handleAPIError(rctx *reqcontext.ReqContext, w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *apiresponses.FailureResponse:
		rctx.Logger.Error(err.LoggerAction(), err)
		a.respond(w, err.ValidatedStatusCode(a.logger), err.ErrorResponse())
	default:
		rctx.Logger.Error("unknown-error", err)
		a.respond(w, http.StatusInternalServerError, apiresponses.ErrorResponse{
			Description: err.Error(),
		})
	}
}

// Endpoints lists service endpoints
func (a API) Endpoints(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
	})
	rctx.Logger.Info("endpoints")

	r, err := a.handler.Endpoints(rctx, instanceID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// ServiceUsage returns service usage
func (a API) ServiceUsage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
	})
	rctx.Logger.Info("service-usage")

	r, err := a.handler.ServiceUsage(rctx, instanceID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// CreateUpdateServiceDefinition is not implemented
func (a API) CreateUpdateServiceDefinition(w http.ResponseWriter, req *http.Request) {
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, nil)
	rctx.Logger.Info("endpoints")

	var sd ServiceDefinitionRequest
	err := json.NewDecoder(req.Body).Decode(&sd)
	if err != nil {
		a.handleAPIError(rctx, w, apiresponses.NewFailureResponse(err, http.StatusBadRequest, "json-unmarshal"))
		return
	}
	defer req.Body.Close()

	err = a.handler.CreateUpdateServiceDefinition(rctx, &sd)
	if err != nil {
		a.handleAPIError(rctx, w, err)
	}
	a.respond(w, http.StatusNoContent, nil)
}

// DeleteServiceDefinition is not implemented
func (a API) DeleteServiceDefinition(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"id": id,
	})
	rctx.Logger.Info("delete-service-definition")

	err := a.handler.DeleteServiceDefinition(rctx, id)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusNoContent, nil)
}

// CreateBackup is not implemented
func (a API) CreateBackup(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
	})
	rctx.Logger.Info("create-backup")

	var br BackupRequest
	err := json.NewDecoder(req.Body).Decode(&br)
	if err != nil {
		a.handleAPIError(rctx, w, apiresponses.NewFailureResponse(err, http.StatusBadRequest, "json-unmarshal"))
		return
	}
	defer req.Body.Close()

	b, err := a.handler.CreateBackup(rctx, instanceID, &br)
	if err != nil {
		a.handleAPIError(rctx, w, err)
	}
	a.respond(w, http.StatusCreated, b)
}

// DeleteBackup is not implemented
func (a API) DeleteBackup(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	backupID := vars["backup_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
		"backup-id":   backupID,
	})
	rctx.Logger.Info("delete-backup")

	r, err := a.handler.DeleteBackup(rctx, instanceID, backupID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// Backup is not implemented
func (a API) Backup(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	backupID := vars["backup_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
		"backup-id":   backupID,
	})
	rctx.Logger.Info("backup")

	r, err := a.handler.Backup(rctx, instanceID, backupID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// ListBackups is not implemented
func (a API) ListBackups(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
	})
	rctx.Logger.Info("list-backups")

	r, err := a.handler.ListBackups(rctx, instanceID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// RestoreBackup is not implemented
func (a API) RestoreBackup(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	backupID := vars["backup_id"]
	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
		"backup-id":   backupID,
	})
	rctx.Logger.Info("restore-backup")

	var restore RestoreRequest
	err := json.NewDecoder(req.Body).Decode(&restore)
	if err != nil {
		a.handleAPIError(rctx, w, apiresponses.NewFailureResponse(err, http.StatusBadRequest, "json-unmarshal"))
		return
	}
	defer req.Body.Close()

	r, err := a.handler.RestoreBackup(rctx, instanceID, backupID, &restore)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// RestoreStatus is not implemented
func (a API) RestoreStatus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]
	backupID := vars["backup_id"]
	restoreID := vars["restore_id"]

	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
		"backup-id":   backupID,
		"restore-id":  restoreID,
	})
	rctx.Logger.Info("restore-status")

	r, err := a.handler.RestoreStatus(rctx, instanceID, backupID, restoreID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}

// APIDocs is not implemented
func (a API) APIDocs(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["service_instance_id"]

	rctx := reqcontext.NewReqContext(req.Context(), a.logger, lager.Data{
		"instance-id": instanceID,
	})
	rctx.Logger.Info("api-docs")

	r, err := a.handler.APIDocs(rctx, instanceID)
	if err != nil {
		a.handleAPIError(rctx, w, err)
		return
	}
	a.respond(w, http.StatusOK, r)
}
