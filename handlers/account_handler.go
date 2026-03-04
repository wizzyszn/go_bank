package handlers

import (
	"net/http"

	"github.com/wizzyszn/go_bank/middleware"
	"github.com/wizzyszn/go_bank/models"
	"github.com/wizzyszn/go_bank/service"
	"github.com/wizzyszn/go_bank/utils"
)

type AccountHandler struct {
	authService        *service.AuthService
	transactionService *service.TransactionService
}

func NewAccountHandler(authService *service.AuthService, transactionService *service.TransactionService) *AccountHandler {

	return &AccountHandler{
		authService:        authService,
		transactionService: transactionService,
	}
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {

	account := middleware.RequireAccount(w, r)

	if account == nil {
		return
	}

	accountData, err := h.authService.GetAccount(account.ID)

	if err != nil {
		utils.WriteNotFound(w, "account not found")
		return
	}

	utils.WriteSuccess(w, accountData)

}

func (h *AccountHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}

	balance, err := h.transactionService.GetBalance(account.ID)

	if err != nil {

		utils.WriteNotFound(w, "Account not found")
		return
	}

	utils.WriteSuccess(w, balance)
}

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}
	var req models.UpdateAccountRequest

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteBadRequest(w, "Invalid request body: "+err.Error())
	}

	updated, err := h.authService.UpdateAccount(account.ID, &req)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteBadRequest(w, validationErr.Error())
			return
		}
		utils.WriteNotFound(w, err.Error())
		return
	}

	utils.WriteSuccess(w, updated)
}
