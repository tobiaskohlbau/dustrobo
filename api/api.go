package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Handler struct {
	mux *chi.Mux
}

func NewHandler() *Handler {
	h := &Handler{
		mux: chi.NewRouter(),
	}

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		MaxAge:         300,
	})
	h.mux.Use(cors.Handler)

	mh := &musicHandler{}

	h.mux.Route("/music", func(r chi.Router) {
		r.Post("/play", mh.Play)
		r.Get("/pause", mh.Pause)
		r.Get("/stop", mh.Stop)
		r.Post("/volume", mh.Volume)
	})

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
