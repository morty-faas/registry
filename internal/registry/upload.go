package registry

import (
	"encoding/json"
	"net/http"
)

func (s *Server) UploadHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&healthcheck{Status: up})
}
