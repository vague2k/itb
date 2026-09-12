package handlers

import (
	"net/http"
	"strings"

	"itb.ihatedoing.work/internal/money"
	"itb.ihatedoing.work/server/services"
	"itb.ihatedoing.work/views/components"
	"itb.ihatedoing.work/views/pages"
)

type CategoryHandler struct {
	Categories   *services.CategoryService
	Transactions *services.TransactionService
	Log          *services.LoggingService
}

func NewCategoryHandler(categories *services.CategoryService, transactions *services.TransactionService, log *services.LoggingService) *CategoryHandler {
	return &CategoryHandler{Categories: categories, Transactions: transactions, Log: log}
}

func (h *CategoryHandler) Index(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Categories.List(r.Context())
	if err != nil {
		handleError(h.Log, w, r, err)
		return
	}
	pages.Categories(rows).Render(r.Context(), w)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	starting, err := parseStartingAmount(r.FormValue("starting_amount"))
	if err != nil {
		handleError(h.Log, w, r, err)
		return
	}

	if err := h.Categories.Create(r.Context(), r.FormValue("name"), starting); err != nil {
		handleError(h.Log, w, r, err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := categoryID(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	category, err := h.Categories.Get(r.Context(), id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	transactions, err := h.Transactions.ListByCategory(r.Context(), id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pages.CategoryDetail(category, transactions).Render(r.Context(), w)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := categoryID(r)
	if err != nil {
		handleError(h.Log, w, r, services.NewUserError("Invalid category"))
		return
	}

	if err := h.Categories.Delete(r.Context(), id); err != nil {
		handleError(h.Log, w, r, err)
		return
	}

	rows, err := h.Categories.List(r.Context())
	if err != nil {
		handleError(h.Log, w, r, err)
		return
	}
	components.CategoriesSection(rows).Render(r.Context(), w)
}

func parseStartingAmount(raw string) (int64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, nil
	}
	cents, err := money.Parse(raw)
	if err != nil {
		return 0, services.NewUserError("Starting amount is not a valid dollar amount")
	}
	return cents, nil
}
