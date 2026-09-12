package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/internal/testdb"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	sqlDB, q := testdb.Open(t)
	return NewHandler(&config.Config{UnderlyingDB: sqlDB, Database: q})
}

func postForm(t *testing.T, h http.Handler, target string, fields map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestCreateCategoryRedirects(t *testing.T) {
	h := newTestHandler(t)
	mux := chi.NewRouter()
	mux.Post("/categories", h.Category.Create)

	w := postForm(t, mux, "/categories", map[string]string{
		"name":            "Groceries",
		"starting_amount": "500.00",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("HX-Redirect"); got != "/" {
		t.Errorf("HX-Redirect = %q, want %q", got, "/")
	}
}

func TestCreateCategoryInvalidAmountToasts(t *testing.T) {
	h := newTestHandler(t)
	mux := chi.NewRouter()
	mux.Post("/categories", h.Category.Create)

	w := postForm(t, mux, "/categories", map[string]string{
		"name":            "Groceries",
		"starting_amount": "not-money",
	})

	if got := w.Header().Get("HX-Retarget"); got != "[data-tui-toaster]" {
		t.Errorf("HX-Retarget = %q, want [data-tui-toaster]", got)
	}
	if got := w.Header().Get("HX-Reswap"); got != "beforeend" {
		t.Errorf("HX-Reswap = %q, want beforeend", got)
	}
	body := w.Body.String()
	if !strings.Contains(body, "data-tui-toast-ssr") {
		t.Errorf("body missing toast stub: %s", body)
	}
	if !strings.Contains(body, "valid dollar amount") {
		t.Errorf("body missing validation message: %s", body)
	}
}

func TestCreateTransactionRendersUpdatedBalance(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()
	if err := h.Category.Categories.Create(ctx, "Groceries", 10000); err != nil {
		t.Fatalf("seed category: %v", err)
	}
	rows, err := h.Category.Categories.List(ctx)
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}

	mux := chi.NewRouter()
	mux.Post("/categories/{id}/transactions", h.Transaction.Create)

	w := postForm(t, mux, "/categories/"+strconv.FormatInt(rows[0].ID, 10)+"/transactions", map[string]string{
		"amount": "25.50",
		"kind":   "expense",
		"note":   "lunch",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "$74.50") {
		t.Errorf("body missing updated balance: %s", w.Body.String())
	}
}
