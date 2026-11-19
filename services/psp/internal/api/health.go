package api

import (
	"encoding/json"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Instance string `json:"instance"`
}

func (h HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	// Instance (container hostname) makes load balancing between replicas visible.
	instance, _ := os.Hostname()

	response := healthResponse{Status: "ok", Database: "ok", Instance: instance}
	statusCode := http.StatusOK

	if err := h.pingDatabase(); err != nil {
		response.Status = "unavailable"
		response.Database = "unreachable"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h HealthHandler) pingDatabase() error {
	sqlDB, err := h.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
