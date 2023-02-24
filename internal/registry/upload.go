package registry

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UploadLinkResponse struct {
	HttpMethod string `json:"httpMethod"`
	UploadLink string `json:"uploadLink"`
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	method, link, err := s.storage.GetUploadLink(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&UploadLinkResponse{HttpMethod: method, UploadLink: link})
}
