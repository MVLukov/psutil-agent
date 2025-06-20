package main

import (
	"html/template"
	"net/http"
	"strings"

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

	r.Mount("/api", handlers.MetricsHandler())
	FileServer(r, "/static", http.Dir("./static"))

	http.ListenAndServe("0.0.0.0:3000", r)

	// almost every return value is a struct

	// convert to JSON. String() is also implemented
	// fmt.Println(v)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	r.Get(path+"/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
