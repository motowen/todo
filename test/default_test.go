package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/router"
)

func TestHealth(t *testing.T) {
	t.Run("health 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.Router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestVersion(t *testing.T) {
	t.Run("version 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/version", nil)
		router.Router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
