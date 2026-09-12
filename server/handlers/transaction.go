package handlers

import (
	"net/http"

	"itb.ihatedoing.work/internal/money"
	"itb.ihatedoing.work/server/services"
	"itb.ihatedoing.work/views/components"
)

type TransactionHandler struct {
	Categories   *services.CategoryService
	Transactions *services.TransactionService
	Log          *services.LoggingService
}

func NewTransactionHandler(categories *services.CategoryService, transactions *services.TransactionService, log *services.LoggingService) *TransactionHandler {
	return &TransactionHandler{Categories: categories, Transactions: transactions, Log: log}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	id, err := categoryID(r)
	if err != nil {
		handleError(h.Log, w, r, services.NewUserError("Invalid category"))
		return
	}

	amount, err := money.Parse(r.FormValue("amount"))
	if err != nil {
		handleError(h.Log, w, r, services.NewUserError("Amount is not a valid dollar amount"))
		return
	}
	if r.FormValue("kind") == "expense" {
		amount = -amount
	}

	if err := h.Transactions.Create(r.Context(), id, amount, r.FormValue("note")); err != nil {
		handleError(h.Log, w, r, err)
		return
	}

	category, err := h.Categories.Get(r.Context(), id)
	if err != nil {
		handleError(h.Log, w, r, err)
		return
	}
	transactions, err := h.Transactions.ListByCategory(r.Context(), id)
	if err != nil {
		handleError(h.Log, w, r, err)
		return
	}

	components.CategoryDetailFragment(category, transactions).Render(r.Context(), w)
}
