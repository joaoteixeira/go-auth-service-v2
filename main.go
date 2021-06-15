package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joaoteixeira/go-auth-service-v2/resource"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	go r.Mount("/users", resource.User{}.Routes())

	r.Route("/api/", func(r chi.Router) {
		r.Post("/login", resource.Login)
	})

	http.ListenAndServe(":8080", r)
}
