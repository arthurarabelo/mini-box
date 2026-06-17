package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret é a chave para assinar os tokens
// No CERNBox, essa chave é configurada via variável de ambiente
var jwtSecret = []byte("minibox-dev-secret-nao-use-em-producao")

// Usuários hardcoded — no CERNBox vêm do LDAP/Keycloak
var users = map[string]struct {
	Password string
	Name     string
	Email    string
}{
	"admin": {Password: "admin", Name: "Admin", Email: "admin@cern.ch"},
	"alice": {Password: "alice123", Name: "Alice Smith", Email: "alice@cern.ch"},
}

// LoginRequest representa o corpo do POST /api/auth/login
// A tag `json:"username"` define o nome no JSON — como Zod no TypeScript
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse é o que devolvemos após login bem-sucedido
type LoginResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Claims são os dados embutidos no JWT
// No CERNBox o equivalente é o reva token com claims do OIDC
type Claims struct {
	Username             string `json:"username"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // struct embutida: exp, iat, etc.
}

// cria um token JWT para o usuário
func GenerateToken(username, name, email string) (string, error) {
	claims := Claims{
		Username: username,
		Name:     name,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}

// verifica o token e retorna os claims
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// processa POST /api/auth/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// decodifica o JSON do corpo da requisição e escreve em req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// verifica credenciais
	user, exists := users[req.Username]
	if !exists || user.Password != req.Password {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// gera o token
	token, err := GenerateToken(req.Username, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		Name:  user.Name,
		Email: user.Email,
	})
}
