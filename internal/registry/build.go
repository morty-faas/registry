package registry

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/polyxia-org/morty-registry/internal/builder"
	"github.com/polyxia-org/morty-registry/pkg/helpers"
	log "github.com/sirupsen/logrus"
)

var (
	ErrRequiredFunctionName    = errors.New("the function name is required")
	ErrRequiredFunctionRuntime = errors.New("the function runtime is required")
	ErrInvalidFunctionArchive  = errors.New("the function code archive must be a valid zip file")
)

func (s *Server) BuildHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the file, the headers of the file from the multipart request form
	f, _, err := r.FormFile("archive")
	if err != nil {
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
		return
	}
	defer f.Close()

	// Validate DTO
	functionName, functionRuntime := r.Form.Get("name"), r.Form.Get("runtime")
	if functionName == "" {
		s.APIErrorResponse(w, makeAPIError(http.StatusBadRequest, ErrRequiredFunctionName))
		return
	}

	if functionRuntime == "" {
		s.APIErrorResponse(w, makeAPIError(http.StatusBadRequest, ErrRequiredFunctionRuntime))
		return
	}

	// Build the function image with the given options
	opts := &builder.BuildOptions{
		Id:      functionName + "-" + uuid.NewString(),
		Runtime: functionRuntime,
		Archive: f,
	}

	image, err := s.builder.ImageBuild(r.Context(), opts)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, builder.ErrInvalidRuntime) {
			status = http.StatusBadRequest
		}

		s.APIErrorResponse(w, makeAPIError(status, err))
		return
	}

	// Compress the rootfs using the LZ4 compression algorithm.
	// We use LZ4 because it is one of the faster compression algorithm.
	// See: https://github.com/lz4/lz4#benchmarks
	archive := fmt.Sprintf("%s.lz4", image)
	if err := helpers.CompressLZ4(image, archive); err != nil {
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
		return
	}

	// Upload the file to the remote storage before returning a response to the user
	f, _ = os.Open(archive)
	if err := s.storage.PutFile(functionName, f); err != nil {
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
		return
	}

	log.Infof("build/%s: function build successful", opts.Id)

	w.WriteHeader(http.StatusNoContent)
}

func makeAPIError(status int, err error) *APIError {
	return &APIError{
		StatusCode: status,
		Message:    err.Error(),
	}
}
