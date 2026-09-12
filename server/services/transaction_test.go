package services

import (
	"context"
	"errors"
	"testing"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/internal/testdb"
)

func TestTransactionUpdatesBalance(t *testing.T) {
	sqlDB, q := testdb.Open(t)
	cfg := &config.Config{UnderlyingDB: sqlDB, Database: q}
	categories := NewCategoryService(cfg)
	transactions := NewTransactionService(cfg)
	ctx := context.Background()

	if err := categories.Create(ctx, "Groceries", 10000); err != nil {
		t.Fatalf("Create category: %v", err)
	}
	rows, err := categories.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if err := transactions.Create(ctx, rows[0].ID, -2500, "lunch"); err != nil {
		t.Fatalf("Create transaction: %v", err)
	}

	category, err := categories.Get(ctx, rows[0].ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if category.BalanceCents != 7500 {
		t.Errorf("balance = %d, want 7500", category.BalanceCents)
	}
}

func TestTransactionRequiresAmount(t *testing.T) {
	sqlDB, q := testdb.Open(t)
	transactions := NewTransactionService(&config.Config{UnderlyingDB: sqlDB, Database: q})

	err := transactions.Create(context.Background(), 1, 0, "")
	var ue UserError
	if !errors.As(err, &ue) {
		t.Fatalf("want UserError, got %v", err)
	}
}
