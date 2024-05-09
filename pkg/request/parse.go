package request

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func ParseUIDFromPath(r *http.Request, isRequired bool) (string, error) {
	uid, err := ParseStringFromPath(r, "uid", true)
	if err != nil {
		return "", err
	}

	if _, err := uuid.Parse(uid); err != nil {
		return "", fmt.Errorf("parsing uid: %w", err)
	}

	return uid, nil
}

func ParseStringFromPath(r *http.Request, field string, isRequired bool) (string, error) {
	res := chi.URLParam(r, field)

	if isRequired && len(res) == 0 {
		return "", fmt.Errorf("field '%s' is required", field)
	}

	return res, nil
}
