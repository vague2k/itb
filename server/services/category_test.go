package services

import (
	"context"
	"errors"
	"testing"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/internal/testdb"
)

func newCategoryService(t *testing.T) (*config.Config, *CategoryService) {
	t.Helper()
	sqlDB, q := testdb.Open(t)
	c := &config.Config{UnderlyingDB: sqlDB, Database: q}
	return c, NewCategoryService(c)
}

func TestCreateCategoryWithStartingAmount(t *testing.T) {
	_, svc := newCategoryService(t)
	ctx := context.Background()

	if err := svc.Create(ctx, "Groceries", 50000); err != nil {
		t.Fatalf("Create: %v", err)
	}

	rows, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("want 1 category, got %d", len(rows))
	}
	if rows[0].BalanceCents != 50000 {
		t.Errorf("balance = %d, want 50000", rows[0].BalanceCents)
	}
}

func TestCreateCategoryRequiresName(t *testing.T) {
	_, svc := newCategoryService(t)

	err := svc.Create(context.Background(), "   ", 0)
	var ue UserError
	if !errors.As(err, &ue) {
		t.Fatalf("want UserError, got %v", err)
	}
}

func TestDeleteCategoryCascadesTransactions(t *testing.T) {
	cfg, svc := newCategoryService(t)
	ctx := context.Background()

	if err := svc.Create(ctx, "Groceries", 10000); err != nil {
		t.Fatalf("Create: %v", err)
	}
	rows, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if err := svc.Delete(ctx, rows[0].ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	var count int
	if err := cfg.UnderlyingDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count); err != nil {
		t.Fatalf("count transactions: %v", err)
	}
	if count != 0 {
		t.Errorf("want 0 transactions after cascade delete, got %d", count)
	}
}
