package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"minibox/backend/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret é a chave para assinar os tokens
// No CERNBox, essa chave é configurada via variável de ambiente
var jwtSecret = []byte("minibox-dev-secret-nao-use-em-producao")

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

	// busca o usuário no banco
	user, err := db.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// compara a senha enviada com o hash armazenado
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// gera o token
	token, err := GenerateToken(user.Username, user.Name, user.Email)
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

// RegisterRequest representa o corpo do POST /api/auth/register.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// RegisterHandler processa POST /api/auth/register.
//
// EXERCÍCIO — implemente o cadastro de usuários. Roteiro:
//  1. Decodifique o JSON do corpo (RegisterRequest) — igual ao LoginHandler.
//  2. Valide os campos (username/password/name/email não vazios, senha
//     com tamanho mínimo razoável).
//  3. Gere o hash da senha com bcrypt.GenerateFromPassword(
//     []byte(req.Password), bcrypt.DefaultCost).
//  4. Chame db.CreateUser(r.Context(), username, hash, name, email) —
//     já implementado em backend/db/users.go.
//  5. Se o erro for db.ErrUserExists, responda http.StatusConflict (409).
//     Para outros erros, http.StatusInternalServerError.
//  6. Em caso de sucesso, gere um token com GenerateToken (igual ao
//     LoginHandler) e responda com LoginResponse e status 201 Created.
//
// Depois de implementar, teste com:
//
//	curl -X POST http://localhost:8080/api/auth/register \
//	  -H "Content-Type: application/json" \
//	  -d '{"username":"bob","password":"bobsenha","name":"Bob","email":"bob@cern.ch"}'
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "cadastro não implementado — exercício", http.StatusNotImplemented)
}
