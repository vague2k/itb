package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/internal/testdb"
)

func newTransactionFixture(t *testing.T) (*sql.DB, *CategoryService, *TransactionService) {
	t.Helper()
	sqlDB, q := testdb.Open(t)
	c := &config.Config{UnderlyingDB: sqlDB, Database: q}
	return sqlDB, NewCategoryService(c), NewTransactionService(c)
}

func TestTransactionUpdatesBalance(t *testing.T) {
	_, categories, transactions := newTransactionFixture(t)
	ctx := context.Background()

	require.NoError(t, categories.Create(ctx, "Groceries", 10000))
	rows, err := categories.List(ctx)
	require.NoError(t, err)

	require.NoError(t, transactions.Create(ctx, rows[0].ID, -2500, "lunch"))

	category, err := categories.Get(ctx, rows[0].ID)
	require.NoError(t, err)
	require.Equal(t, int64(7500), category.BalanceCents)
}

func TestTransactionRequiresAmount(t *testing.T) {
	_, _, transactions := newTransactionFixture(t)

	err := transactions.Create(context.Background(), 1, 0, "")

	var ue UserError
	require.ErrorAs(t, err, &ue)
	require.Equal(t, "Amount is required", ue.Error())
}

func TestClosedDBTransactionReturnsInternalError(t *testing.T) {
	sqlDB, _, transactions := newTransactionFixture(t)
	sqlDB.Close()

	err := transactions.Create(context.Background(), 1, -2500, "lunch")

	var ie InternalError
	require.ErrorAs(t, err, &ie)
}
