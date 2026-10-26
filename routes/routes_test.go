package routes_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/trikunc/go_vibe/routes"
)

func TestRoutesEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Change working directory to project root if needed
	_ = os.Chdir("..")
	defer func() {
		_ = os.Chdir("routes")
	}()

	r := routes.SetupRouter(nil)

	endpoints := []string{"/", "/detail", "/compare"}
	for _, ep := range endpoints {
		req, _ := http.NewRequest(http.MethodGet, ep, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for %s, got %d", ep, w.Code)
		}
	}
}
