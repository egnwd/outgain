package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func performRequest(h http.Handler, method, path string) {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
}

func TestRoutes(t *testing.T) {
	path := "/example"
	buffer := new(bytes.Buffer)
	mux := mux.NewRouter()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {})
	logger := ServerLogger(buffer, mux)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		performRequest(logger, method, path)
		assert.Contains(t, buffer.String(), method)
		assert.Contains(t, buffer.String(), strconv.Itoa(http.StatusOK))
		assert.Contains(t, buffer.String(), path)

		buffer.Reset()
	}

	performRequest(logger, http.MethodGet, "/notfound")
	assert.Contains(t, buffer.String(), http.MethodGet)
	assert.Contains(t, buffer.String(), strconv.Itoa(http.StatusNotFound))
	assert.Contains(t, buffer.String(), "/notfound")

}

func TestStatusColours(t *testing.T) {
	assert.Equal(t, colourForStatus(http.StatusOK), green, "200 should be green")

	assert.Equal(t, colourForStatus(http.StatusMovedPermanently), white, "301 should be white")
	assert.Equal(t, colourForStatus(http.StatusTemporaryRedirect), white, "302 should be white")
	assert.Equal(t, colourForStatus(http.StatusNotModified), white, "304 should be white")

	assert.Equal(t, colourForStatus(http.StatusForbidden), yellow, "401 should be yellow")
	assert.Equal(t, colourForStatus(http.StatusNotFound), yellow, "404 should be yellow")

	assert.Equal(t, colourForStatus(http.StatusInternalServerError), red, "500 should be red")
}

func TestMethodColours(t *testing.T) {
	assert.Equal(t, colourForMethod(http.MethodGet), cyan, "GET should be cyan")
	assert.Equal(t, colourForMethod(http.MethodPost), yellow, "POST should be yellow")
	assert.Equal(t, colourForMethod(http.MethodPut), blue, "PUT should be blue")
	assert.Equal(t, colourForMethod(http.MethodDelete), red, "DELETE should be red")
	assert.Equal(t, colourForMethod(http.MethodPatch), green, "PATCH should be green")
	assert.Equal(t, colourForMethod(http.MethodHead), magenta, "HEAD should be magenta")
	assert.Equal(t, colourForMethod(http.MethodOptions), white, "OPTIONS should be white")
	assert.Equal(t, colourForMethod(http.MethodTrace), reset, "TRACE should be reset")
	assert.Equal(t, colourForMethod(http.MethodConnect), reset, "CONNECT should be reset")
}
