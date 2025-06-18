package main

import (
	"html/template"
	"net/http"

	"github.com/MVLukov/psutil-dash/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})

	r.Mount("/api", handlers.Basic())

	http.ListenAndServe(":3000", r)

	// almost every return value is a struct

	// convert to JSON. String() is also implemented
	// fmt.Println(v)
}
