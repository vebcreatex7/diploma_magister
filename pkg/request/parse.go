package request

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
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

func ParseStringFromQuery(r *http.Request, field string, isRequired bool) (string, error) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("parsing query: %w", err)
	}

	res := q.Get(field)

	if isRequired && len(res) == 0 {
		return "", fmt.Errorf("field '%s' is required", field)
	}

	return res, nil
}

func ParseTimeFromQuery(r *http.Request, field string, isRequired bool) (time.Time, error) {
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing query: %w", err)
	}

	res := q.Get(field)

	if isRequired && len(res) == 0 {
		return time.Time{}, fmt.Errorf("field '%s' is required", field)
	}

	t, err := time.Parse(time.DateOnly, res)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing time: %w", err)
	}

	return t, nil
}
