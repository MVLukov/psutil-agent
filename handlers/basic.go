package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/MVLukov/psutil-dash/metrics"
	"github.com/go-chi/chi/v5"
)

func Basic() http.Handler {
	router := chi.NewRouter()

	router.Get("/basicJSON", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics.GetBasicMetrics())
	})

	router.Get("/basic", func(w http.ResponseWriter, r *http.Request) {
		metrics := metrics.GetBasicMetrics()

		tmpl := template.Must(template.ParseFiles("templates/basic.html"))
		tmpl.Execute(w, metrics)
	})

	return router
}
