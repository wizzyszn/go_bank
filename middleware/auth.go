package middleware

import (
	"context"
	"net/http"
	"strings"
	"github.com/wizzyszn/go_bank/models"
	"github.com/wizzyszn/go_bank/service"
	"github.com/wizzyszn/go_bank/utils"
)

type contextKey string

const (
	ContextKeyAccount contextKey = "account"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.WriteUnAuthorized(w, "Missing authorization")
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteUnAuthorized(w, "Invalid authorization header format")
			return
		}

		sessionID := parts[1]

		if sessionID == "" {
			utils.WriteUnAuthorized(w, "Missing session token")
			return
		}

		account, err := m.authService.ValidateSession(sessionID)

		if err != nil {
			utils.WriteUnAuthorized(w, "Invalid or expired session")
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyAccount, account)
		r = r.WithContext(ctx)
		next(w, r)

	}

}

func GetAccountFromContext(ctx context.Context) (*models.Account, bool) {
	account, ok := ctx.Value(ContextKeyAccount).(*models.Account)

	return account, ok
}
func RequireAccount(w http.ResponseWriter, r *http.Request) *models.Account {
	account, ok := GetAccountFromContext(r.Context())
	if !ok {
		utils.WriteUnAuthorized(w, "Authentication required")
		return nil
	}
	return account
}
