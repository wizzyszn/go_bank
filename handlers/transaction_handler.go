package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/wizzyszn/go_bank/middleware"
	"github.com/wizzyszn/go_bank/models"
	"github.com/wizzyszn/go_bank/service"
	"github.com/wizzyszn/go_bank/utils"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {

	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}
	var req models.DepositRequest

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	transaction, err := h.transactionService.Deposit(account.ID, &req)

	if err != nil {
		if Validation, ok := err.(*utils.ValidationError); ok {
			utils.WriteBadRequest(w, Validation.Error())
			return
		}
		utils.WriteBadRequest(w, err.Error())
		return
	}
	utils.WriteCreated(w, transaction)

}

func (h *TransactionHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req models.WitdrawRequest

	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteBadRequest(w, "Invalid request body"+err.Error())
		return
	}

	transaction, err := h.transactionService.WithDraw(account.ID, &req)
	if err != nil {
		if ValidationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteBadRequest(w, ValidationErr.Error())
			return
		}
		utils.WriteBadRequest(w, err.Error())
		return
	}

	utils.WriteCreated(w, transaction)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req models.TransferRequest

	account := middleware.RequireAccount(w, r)

	if account == nil {
		return
	}

	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteBadRequest(w, "Invalid request body: "+err.Error())
		return
	}

	transaction, err := h.transactionService.Transfer(account.ID, &req)
	if err != nil {
		if ValidationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteBadRequest(w, ValidationErr.Error())
			return
		}
		utils.WriteBadRequest(w, err.Error())
		return
	}

	utils.WriteCreated(w, transaction)
}

func (h *TransactionHandler) GetTransations(w http.ResponseWriter, r *http.Request) {
	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}

	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	transactions, err := h.transactionService.GetTransactions(account.ID, page, limit)

	if err != nil {
		utils.WriteInternalError(w, "")
		return
	}
	utils.WriteSuccess(w, transactions)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	account := middleware.RequireAccount(w, r)
	if account == nil {
		return
	}

	path := r.URL.Path

	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 3 {
		utils.WriteBadRequest(w, "Invalid transaction ID")
		return
	}

	transactionID, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.WriteBadRequest(w, "Invalid transaction ID")
		return
	}

	transaction, err := h.transactionService.GetTransaction(account.ID, transactionID)

	if err != nil {
		utils.WriteNotFound(w, "Transaction not found")
		return
	}

	utils.WriteSuccess(w, transaction)
}
