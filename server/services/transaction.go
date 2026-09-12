package services

import (
	"context"
	"strings"

	"itb.ihatedoing.work/config"
	db "itb.ihatedoing.work/server/database/generated"
	"itb.ihatedoing.work/server/database/models"
)

type TransactionService struct {
	config *config.Config
}

func NewTransactionService(c *config.Config) *TransactionService {
	return &TransactionService{config: c}
}

func (s *TransactionService) ListByCategory(ctx context.Context, categoryID int64) ([]models.Transaction, error) {
	transactions, err := s.config.Database.ListTransactionsByCategory(ctx, categoryID)
	if err != nil {
		return nil, NewInternalError(err, "could not list transactions")
	}
	return transactions, nil
}

func (s *TransactionService) Create(ctx context.Context, categoryID, amountCents int64, note string) error {
	if amountCents == 0 {
		return NewUserError("Amount is required")
	}
	if _, err := s.config.Database.CreateTransaction(ctx, db.CreateTransactionParams{
		CategoryID:  categoryID,
		AmountCents: amountCents,
		Note:        strings.TrimSpace(note),
	}); err != nil {
		return NewInternalError(err, "could not add transaction")
	}
	return nil
}
