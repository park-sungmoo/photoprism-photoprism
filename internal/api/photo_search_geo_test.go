package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchPhotosGeo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app, router, _ := NewApiTest()

		SearchPhotosGeo(router)

		result := PerformRequest(app, "GET", "/api/v1/geo")
		assert.Equal(t, http.StatusOK, result.Code)
	})
}
