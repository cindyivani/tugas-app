package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/trikunc/go_vibe/config"
	"github.com/trikunc/go_vibe/handlers"
	"github.com/trikunc/go_vibe/models"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := handlers.NewProteinHandler(nil)
	r.GET("/api/health", handler.HealthCheck)

	req, _ := http.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse json response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Fatalf("Expected status 'healthy', got '%v'", resp["status"])
	}
}

func TestProteinParseAminoAcids(t *testing.T) {
	p := models.Protein{
		Name:             "Test Protein",
		AminoAcidProfile: `[{"name":"Leusin","code":"Leu","amount":8.8,"is_essential":true}]`,
	}

	p.ParseAminoAcids()

	if len(p.ParsedAminoAcids) != 1 {
		t.Fatalf("Expected 1 amino acid, got %d", len(p.ParsedAminoAcids))
	}

	if p.ParsedAminoAcids[0].Code != "Leu" {
		t.Fatalf("Expected code 'Leu', got '%s'", p.ParsedAminoAcids[0].Code)
	}

	if !p.ParsedAminoAcids[0].IsEssential {
		t.Fatalf("Expected IsEssential to be true")
	}
}

func TestGetProteinsIntegration(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := handlers.NewProteinHandler(db)
	r.GET("/api/proteins", handler.GetProteins)

	req, _ := http.NewRequest(http.MethodGet, "/api/proteins", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool             `json:"success"`
		Total   int              `json:"total"`
		Data    []models.Protein `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("Expected success true")
	}
}

func TestGetProteinByIDIntegration(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := handlers.NewProteinHandler(db)
	r.GET("/api/proteins/:id", handler.GetProteinByID)

	req, _ := http.NewRequest(http.MethodGet, "/api/proteins/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool           `json:"success"`
		Data    models.Protein `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !resp.Success || resp.Data.ID != 1 {
		t.Fatalf("Expected protein with ID 1, got %+v", resp.Data)
	}
}

func TestCompareProteinsIntegration(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := handlers.NewProteinHandler(db)
	r.GET("/api/proteins/compare", handler.CompareProteins)

	req, _ := http.NewRequest(http.MethodGet, "/api/proteins/compare?ids=1,2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool             `json:"success"`
		Count   int              `json:"count"`
		Data    []models.Protein `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !resp.Success || resp.Count < 2 {
		t.Fatalf("Expected at least 2 proteins in compare, got %d", resp.Count)
	}
}

func TestGetStatsIntegration(t *testing.T) {
	db, err := config.InitDB()
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
		return
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := handlers.NewProteinHandler(db)
	r.GET("/api/stats", handler.GetStats)

	req, _ := http.NewRequest(http.MethodGet, "/api/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Total       int `json:"total"`
			HewaniCount int `json:"hewani_count"`
			NabatiCount int `json:"nabati_count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !resp.Success || resp.Data.Total == 0 {
		t.Fatalf("Expected non-zero stats total")
	}
}
