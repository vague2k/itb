package handlers

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/internal/testdb"
)

func newTestHandler(t *testing.T) (*Handler, *sql.DB) {
	t.Helper()
	sqlDB, q := testdb.Open(t)
	c := &config.Config{UnderlyingDB: sqlDB, Database: q}
	return NewHandler(c), sqlDB
}

func categoriesRouter(h *Handler) http.Handler {
	mux := chi.NewRouter()
	mux.Post("/categories", h.Category.Create)
	mux.Delete("/categories/{id}", h.Category.Delete)
	mux.Post("/categories/{id}/transactions", h.Transaction.Create)
	return mux
}

func serveRequest(h http.Handler, method, target string, body io.Reader, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, body)
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func urlEncoded(fields map[string]string) (io.Reader, string) {
	vals := url.Values{}
	for k, v := range fields {
		vals.Set(k, v)
	}
	return strings.NewReader(vals.Encode()), "application/x-www-form-urlencoded"
}

func TestCreateCategorySuccess(t *testing.T) {
	h, _ := newTestHandler(t)
	router := categoriesRouter(h)
	ctx := context.Background()

	body, ct := urlEncoded(map[string]string{"name": "Groceries", "starting_amount": "500.00"})
	w := serveRequest(router, http.MethodPost, "/categories", body, ct)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "/", w.Header().Get("HX-Redirect"))

	list, err := h.Category.Categories.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Groceries", list[0].Name)
	require.Equal(t, int64(50000), list[0].BalanceCents)
}

func TestCreateCategoryMalformedAmount(t *testing.T) {
	h, _ := newTestHandler(t)

	body, ct := urlEncoded(map[string]string{"name": "Groceries", "starting_amount": "not-money"})
	w := serveRequest(categoriesRouter(h), http.MethodPost, "/categories", body, ct)

	require.Equal(t, "[data-tui-toaster]", w.Header().Get("HX-Retarget"))
	require.Equal(t, "beforeend", w.Header().Get("HX-Reswap"))
	require.Empty(t, w.Header().Get("HX-Redirect"), "malformed amount must not redirect")
	require.Contains(t, w.Body.String(), "data-tui-toast-ssr")
	require.Contains(t, w.Body.String(), "valid dollar amount")
}

func TestDeleteCategoryViaHandler(t *testing.T) {
	h, sqlDB := newTestHandler(t)
	router := categoriesRouter(h)
	ctx := context.Background()

	require.NoError(t, h.Category.Categories.Create(ctx, "Groceries", 10000))
	rows, err := h.Category.Categories.List(ctx)
	require.NoError(t, err)

	target := "/categories/" + strconv.FormatInt(rows[0].ID, 10)
	w := serveRequest(router, http.MethodDelete, target, nil, "")

	require.Equal(t, http.StatusOK, w.Code)

	var count int
	require.NoError(t, sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM categories").Scan(&count))
	require.Zero(t, count)
}

func TestCreateTransactionViaHandler(t *testing.T) {
	h, _ := newTestHandler(t)
	router := categoriesRouter(h)
	ctx := context.Background()

	require.NoError(t, h.Category.Categories.Create(ctx, "Groceries", 10000))
	rows, err := h.Category.Categories.List(ctx)
	require.NoError(t, err)

	body, ct := urlEncoded(map[string]string{"amount": "25.50", "kind": "expense", "note": "lunch"})
	target := "/categories/" + strconv.FormatInt(rows[0].ID, 10) + "/transactions"
	w := serveRequest(router, http.MethodPost, target, body, ct)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "$74.50")
}
