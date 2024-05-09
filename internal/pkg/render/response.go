package render

import (
	"bytes"
	"html/template"
	"net/http"
)

func ErrResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		return
	}
}

func HTMLResponse(
	w http.ResponseWriter,
	t *template.Template,
	name string,
	data any,
	headers map[string]string,
	statusCode int,
) {
	var buf bytes.Buffer

	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		ErrResponse(w, err)
		return
	}

	w.Header().Set("content-type", "text/html")

	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(statusCode)

	if _, err := w.Write(buf.Bytes()); err != nil {
		ErrResponse(w, err)
		return
	}
}
