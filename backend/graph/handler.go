package graph

import (
	"encoding/json"
	"net/http"

	"minibox/backend/auth"
)

// MeResponse espelha a estrutura do /graph/v1.0/me do CERNBox
type MeResponse struct {
	ID                       string `json:"id"`
	DisplayName              string `json:"displayName"`
	Mail                     string `json:"mail"`
	OnPremisesSamAccountName string `json:"onPremisesSamAccountName"`
}

// MeHandler processa GET /api/graph/v1.0/me
func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r) // lê os claims injetados pelo middleware

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MeResponse{
		ID:                       claims.Username,
		DisplayName:              claims.Name,
		Mail:                     claims.Email,
		OnPremisesSamAccountName: claims.Username,
	})
}
