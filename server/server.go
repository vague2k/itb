package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"itb.ihatedoing.work/config"
	components "itb.ihatedoing.work/views/shadcn"
)

type Server struct {
	port   string
	router *chi.Mux
}

func NewServer(c *config.Config) *Server {
	s := &Server{
		port:   c.Port,
		router: chi.NewRouter(),
	}
	s.router.Use(middleware.Logger, middleware.Recoverer)
	return s
}

func (s *Server) Run() error {
	server := http.Server{
		Addr:    s.port,
		Handler: s.router,
	}
	return server.ListenAndServe()
}

func (s *Server) SetupAssetsRoutes() {
	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("GO_ENV") == "production" {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		} else {
			w.Header().Set("Cache-Control", "no-store")
		}
		http.FileServer(http.Dir("./views/assets")).ServeHTTP(w, r)
	})
	s.router.Handle("GET /assets/*", http.StripPrefix("/assets/", assetHandler))

	s.router.Handle("GET /components/{bundle}", components.ScriptsHandler())
}
