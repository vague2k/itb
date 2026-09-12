package server

import (
	"itb.ihatedoing.work/server/handlers"
)

func (s *Server) RegisterRoutes(h *handlers.Handler) {
	s.router.Get("/", h.Category.Index)
	s.router.Post("/categories", h.Category.Create)
	s.router.Get("/categories/{id}", h.Category.Show)
	s.router.Delete("/categories/{id}", h.Category.Delete)
	s.router.Post("/categories/{id}/transactions", h.Transaction.Create)
}
