package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"itb.ihatedoing.work/config"
	"itb.ihatedoing.work/views/assets"
	components "itb.ihatedoing.work/views/shadcn"
)

type Server struct {
	port     string
	deployed bool
	router   *chi.Mux
}

func NewServer(c *config.Config) *Server {
	s := &Server{
		port:     ":" + c.Port,
		deployed: c.Deployed(),
		router:   chi.NewRouter(),
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
		s.noStoreWhenLocal(w)
		if s.deployed && r.URL.Query().Get("v") != "" {
			// Content-versioned assets never change under the same URL.
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		fs := http.FileServer(http.FS(assets.Assets))
		if !s.deployed {
			fs = http.FileServer(http.Dir("./views/assets"))
		}

		fs.ServeHTTP(w, r)
	})
	s.router.Handle("GET /assets/*", http.StripPrefix("/assets/", assetHandler))

	s.router.Handle("GET /components/{bundle}", components.ScriptsHandler())
}

// noStoreWhenLocal disables caching in dev so asset edits show up on reload;
// deployed builds keep normal browser caching.
func (s *Server) noStoreWhenLocal(w http.ResponseWriter) {
	if !s.deployed {
		w.Header().Set("Cache-Control", "no-store")
	}
}
