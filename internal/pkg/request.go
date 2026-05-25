package pkg

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var validate = validator.New()

func DecodeAndValidateBody[T any](w http.ResponseWriter, r *http.Request, logger *zerolog.Logger) (T, int, error, bool) {
	var req T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		logger.Error().Err(err).Msg("decode body request error")
		return req, http.StatusBadRequest, errors.New("invalid body request"), false
	}

	if err := validate.Struct(req); err != nil {
		logger.Error().Err(err).Msg("validation body requst error")
		return req, http.StatusBadRequest, err, false
	}

	return req, 0, nil, true
}
