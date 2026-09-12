package services

import (
	"context"
	"strings"

	"itb.ihatedoing.work/config"
	db "itb.ihatedoing.work/server/database/generated"
)

type CategoryService struct {
	config *config.Config
}

func NewCategoryService(c *config.Config) *CategoryService {
	return &CategoryService{config: c}
}

func (s *CategoryService) List(ctx context.Context) ([]db.ListCategoriesRow, error) {
	categories, err := s.config.Database.ListCategories(ctx)
	if err != nil {
		return nil, NewInternalError(err, "could not list categories")
	}
	return categories, nil
}

func (s *CategoryService) Get(ctx context.Context, id int64) (db.GetCategoryRow, error) {
	category, err := s.config.Database.GetCategory(ctx, id)
	if err != nil {
		return db.GetCategoryRow{}, NewInternalError(err, "could not get category")
	}
	return category, nil
}

// Create makes a category and, when starting cents are non-zero, records them
// as its first transaction so the starting amount is editable like any other.
func (s *CategoryService) Create(ctx context.Context, name string, startingCents int64) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return NewUserError("Name is required")
	}

	tx, err := s.config.UnderlyingDB.BeginTx(ctx, nil)
	if err != nil {
		return NewInternalError(err, "could not begin transaction")
	}
	defer tx.Rollback()

	q := s.config.Database.WithTx(tx)
	category, err := q.CreateCategory(ctx, name)
	if err != nil {
		return NewInternalError(err, "could not create category")
	}

	if startingCents != 0 {
		if _, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
			CategoryID:  category.ID,
			AmountCents: startingCents,
			Note:        "Starting amount",
		}); err != nil {
			return NewInternalError(err, "could not create starting amount")
		}
	}

	if err := tx.Commit(); err != nil {
		return NewInternalError(err, "could not commit transaction")
	}
	return nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.config.Database.DeleteCategory(ctx, id); err != nil {
		return NewInternalError(err, "could not delete category")
	}
	return nil
}
