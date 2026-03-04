package handlers

import (
	"net/http"

	"github.com/wizzyszn/go_bank/middleware"
	"github.com/wizzyszn/go_bank/models"
	"github.com/wizzyszn/go_bank/service"
	"github.com/wizzyszn/go_bank/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthService(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAccountRequest

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteBadRequest(w, "Invalid request body:"+err.Error())
		return
	}

	account, err := h.authService.Register(&req)

	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteBadRequest(w, validationErr.Error())
			return
		}
		utils.WriteBadRequest(w, err.Error())
		return
	}
	utils.WriteCreated(w, account)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginAccountRequest

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteBadRequest(w, "Invalid request body:"+err.Error())

		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		utils.WriteUnAuthorized(w, err.Error())
		return
	}
	utils.WriteSuccess(w, res)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	account := middleware.RequireAccount(w, r)

	if account == nil {
		return
	}

	authHeader := r.Header.Get("Authorization")
	sessionID := ""

	if len(authHeader) > 7 {
		sessionID = authHeader[7:]
	}

	if sessionID == "" {
		utils.WriteBadRequest(w, "Missing Session ID")
		return
	}

	if err := h.authService.Logout(sessionID); err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}
	utils.WriteSuccess(w, map[string]string{
		"message": "Logged out successfully",
	})
}
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}

	utils.WriteSuccess(w, account.ToResponse())
}
