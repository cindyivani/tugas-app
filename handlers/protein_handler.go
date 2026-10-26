package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trikunc/go_vibe/models"
	"gorm.io/gorm"
)

type ProteinHandler struct {
	db *gorm.DB
}

func NewProteinHandler(db *gorm.DB) *ProteinHandler {
	return &ProteinHandler{db: db}
}

// GetProteins handles GET /api/proteins with search, category filtering, and sorting
func (h *ProteinHandler) GetProteins(c *gin.Context) {
	var proteins []models.Protein
	query := h.db.Model(&models.Protein{})

	// Search filter
	search := strings.TrimSpace(c.Query("search"))
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where(
			"name ILIKE ? OR common_name ILIKE ? OR food_source ILIKE ? OR pdb_id ILIKE ? OR uniprot_id ILIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// Category filter (Hewani / Nabati)
	category := strings.TrimSpace(c.Query("category"))
	if category != "" && category != "Semua" {
		query = query.Where("source_category = ?", category)
	}

	// Sorting
	sortBy := c.DefaultQuery("sort", "name_asc")
	switch sortBy {
	case "name_desc":
		query = query.Order("name DESC")
	case "mw_asc":
		query = query.Order("molecular_weight ASC")
	case "mw_desc":
		query = query.Order("molecular_weight DESC")
	case "score_desc":
		query = query.Order("digestibility_score DESC")
	default:
		query = query.Order("name ASC")
	}

	if err := query.Find(&proteins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal mengambil data katalog protein: " + err.Error(),
		})
		return
	}

	// Parse amino acids for each protein
	for i := range proteins {
		proteins[i].ParseAminoAcids()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"total":   len(proteins),
		"data":    proteins,
	})
}

// GetProteinByID handles GET /api/proteins/:id
func (h *ProteinHandler) GetProteinByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "ID protein harus berupa angka numerik",
		})
		return
	}

	var protein models.Protein
	if err := h.db.First(&protein, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Protein dengan ID tersebut tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Terjadi kesalahan basis data: " + err.Error(),
		})
		return
	}

	protein.ParseAminoAcids()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    protein,
	})
}

// CompareProteins handles GET /api/proteins/compare?ids=1,2,3
func (h *ProteinHandler) CompareProteins(c *gin.Context) {
	idsParam := strings.TrimSpace(c.Query("ids"))
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Parameter 'ids' diperlukan (contoh: ?ids=1,2)",
		})
		return
	}

	idStrings := strings.Split(idsParam, ",")
	var ids []int
	for _, s := range idStrings {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if id, err := strconv.Atoi(s); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Daftar ID tidak valid",
		})
		return
	}

	var proteins []models.Protein
	if err := h.db.Where("id IN ?", ids).Find(&proteins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal memuat data komparasi: " + err.Error(),
		})
		return
	}

	for i := range proteins {
		proteins[i].ParseAminoAcids()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"count":   len(proteins),
		"data":    proteins,
	})
}

// GetStats handles GET /api/stats
func (h *ProteinHandler) GetStats(c *gin.Context) {
	var total int64
	var hewaniCount int64
	var nabatiCount int64

	h.db.Model(&models.Protein{}).Count(&total)
	h.db.Model(&models.Protein{}).Where("source_category = ?", "Hewani").Count(&hewaniCount)
	h.db.Model(&models.Protein{}).Where("source_category = ?", "Nabati").Count(&nabatiCount)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total":        total,
			"hewani_count": hewaniCount,
			"nabati_count": nabatiCount,
		},
	})
}

// HealthCheck handles GET /api/health
func (h *ProteinHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"app":     "Protein Pangan API",
		"version": "1.0.0",
	})
}
