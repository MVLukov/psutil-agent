package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/MVLukov/psutil-dash/metrics"
	"github.com/go-chi/chi/v5"
)

func MetricsHandler() http.Handler {
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

	router.Get("/disksJSON", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(metrics.GetDisksMetrics())

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics.GetDisksMetrics())
	})

	router.Get("/disks", func(w http.ResponseWriter, r *http.Request) {
		metrics := metrics.GetDisksMetrics()

		fmt.Println(metrics)

		tmpl := template.Must(template.ParseFiles("templates/disks.html"))
		tmpl.Execute(w, metrics)
	})

	return router
}
