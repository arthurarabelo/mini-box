package auth

import (
	"context"
	"net/http"
	"strings"
)

// contextKey é um tipo privado para evitar colisão de chaves no contexto
// O contexto HTTP (r.Context()) é como passar dados entre middlewares —
// equivalente ao res.locals no Express
type contextKey string

const ClaimsKey contextKey = "claims"

// RequireAuth é o middleware que protege rotas
// Lê o header "Authorization: Bearer <token>", valida, e injeta os claims no contexto
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Token ausente", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Injeta os claims no contexto da requisição
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// GetClaims extrai os claims do contexto — usado nos handlers protegidos
func GetClaims(r *http.Request) *Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*Claims)
	return claims
}
