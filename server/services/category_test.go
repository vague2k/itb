package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/internal/testdb"
)

func newCategoryService(t *testing.T) (*sql.DB, *CategoryService) {
	t.Helper()
	sqlDB, q := testdb.Open(t)
	c := &config.Config{UnderlyingDB: sqlDB, Database: q}
	return sqlDB, NewCategoryService(c)
}

func TestCreateCategoryWithStartingAmount(t *testing.T) {
	_, svc := newCategoryService(t)
	ctx := context.Background()

	require.NoError(t, svc.Create(ctx, "Groceries", 50000))

	rows, err := svc.List(ctx)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "Groceries", rows[0].Name)
	require.Equal(t, int64(50000), rows[0].BalanceCents)
}

func TestCreateCategoryRequiresName(t *testing.T) {
	_, svc := newCategoryService(t)

	err := svc.Create(context.Background(), "   ", 0)

	var ue UserError
	require.ErrorAs(t, err, &ue)
	require.Equal(t, "Name is required", ue.Error())
}

func TestDeleteCategoryCascadesTransactions(t *testing.T) {
	sqlDB, svc := newCategoryService(t)
	ctx := context.Background()

	require.NoError(t, svc.Create(ctx, "Groceries", 10000))
	rows, err := svc.List(ctx)
	require.NoError(t, err)

	require.NoError(t, svc.Delete(ctx, rows[0].ID))

	var count int
	require.NoError(t, sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count))
	require.Zero(t, count)
}

func TestClosedDBCategoryReturnsInternalError(t *testing.T) {
	sqlDB, svc := newCategoryService(t)
	sqlDB.Close()

	err := svc.Create(context.Background(), "Groceries", 0)

	var ie InternalError
	require.ErrorAs(t, err, &ie)
}
