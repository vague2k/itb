package handlers

import (
	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/server/services"
)

type Handler struct {
	Category    *CategoryHandler
	Transaction *TransactionHandler
}

func NewHandler(c *config.Config) *Handler {
	log := services.NewLoggingService()
	categories := services.NewCategoryService(c)
	transactions := services.NewTransactionService(c)

	return &Handler{
		Category:    NewCategoryHandler(categories, transactions, log),
		Transaction: NewTransactionHandler(categories, transactions, log),
	}
}
