package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var validate = validator.New()

// maxRequestBodyBytes caps JSON request bodies. All call sites decode small
// DTOs (project/deployment params, refresh tokens), so 1 MiB leaves generous
// headroom while bounding decode memory per request.
const maxRequestBodyBytes = 1 << 20

func DecodeAndValidateBody[T any](w http.ResponseWriter, r *http.Request, logger *zerolog.Logger) (T, int, error, bool) {
	var req T

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			logger.Warn().Msg("request body too large")
			return req, http.StatusRequestEntityTooLarge, errors.New("request body too large"), false
		}
		// Client-caused 400s are Warn, not Error: they must not pollute
		// error dashboards. io.EOF specifically means an empty body,
		// which is routine (curl without -d, buggy client, probe).
		if errors.Is(err, io.EOF) {
			logger.Warn().Msg("empty request body")
		} else {
			logger.Warn().Err(err).Msg("decode body request error")
		}
		return req, http.StatusBadRequest, errors.New("invalid body request"), false
	}

	if err := validate.Struct(req); err != nil {
		logger.Warn().Err(err).Msg("validation body request error")
		return req, http.StatusBadRequest, err, false
	}

	return req, 0, nil, true
}
