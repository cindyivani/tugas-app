package models

import (
	"encoding/json"
	"time"
)

// AminoAcid represents an individual amino acid component and its amount
type AminoAcid struct {
	Name        string  `json:"name"`         // e.g. "Leusin", "Lisin"
	Code        string  `json:"code"`         // e.g. "Leu", "Lys"
	Amount      float64 `json:"amount"`       // in g / 100g protein (percentage)
	IsEssential bool    `json:"is_essential"` // true if essential amino acid
}

// Protein represents a food protein entity stored in the database
type Protein struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Name               string     `gorm:"size:120;not null;index" json:"name"`
	CommonName         string     `gorm:"size:120" json:"common_name"`
	SourceCategory     string     `gorm:"size:50;not null;index" json:"source_category"` // "Hewani" atau "Nabati"
	FoodSource         string     `gorm:"size:150;not null" json:"food_source"`
	MolecularWeight    float64    `json:"molecular_weight"`                            // in kDa
	IsoelectricPoint   float64    `json:"isoelectric_point"`                           // pI value
	DigestibilityScore float64    `json:"digestibility_score"`                         // DIAAS / PDCAAS score (0-100)
	BiologicalFunction string     `gorm:"type:text" json:"biological_function"`        // Deskripsi fungsi fisiologis & pangan
	PdbID              string     `gorm:"size:20;index" json:"pdb_id"`                 // RCSB PDB code for 3D viewer (e.g. "1OVA")
	UniprotID          string     `gorm:"size:20;index" json:"uniprot_id"`             // UniProt accession code (e.g. "P01012")
	ReferenceURL       string     `gorm:"size:255" json:"reference_url"`               // Direct scientific citation/database URL
	Description        string     `gorm:"type:text" json:"description"`                // Ringkasan edukasi
	AminoAcidProfile   string     `gorm:"type:text" json:"amino_acid_profile"`         // JSON serialized string of []AminoAcid
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Parsed field (not stored directly in DB column)
	ParsedAminoAcids   []AminoAcid `gorm:"-" json:"parsed_amino_acids,omitempty"`
}

// ParseAminoAcids populates ParsedAminoAcids from AminoAcidProfile JSON string
func (p *Protein) ParseAminoAcids() {
	if p.AminoAcidProfile == "" {
		p.ParsedAminoAcids = []AminoAcid{}
		return
	}
	var list []AminoAcid
	if err := json.Unmarshal([]byte(p.AminoAcidProfile), &list); err == nil {
		p.ParsedAminoAcids = list
	} else {
		p.ParsedAminoAcids = []AminoAcid{}
	}
}
