package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/MVLukov/psutil-dash/metrics"
	"github.com/go-chi/chi/v5"
)

type Basic struct {
	IconURL string `json:"iconURL"`
	Metrics metrics.BasicMetrics
}

func MetricsHandler() http.Handler {
	router := chi.NewRouter()

	router.Get("/basicJSON", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics.GetBasicMetrics())
	})

	router.Get("/basic", func(w http.ResponseWriter, r *http.Request) {
		metrics := metrics.GetBasicMetrics()

		basic := Basic{}

		if metrics.HostINFO.OS.ID != "" {
			basic.IconURL = getIcon(metrics.HostINFO.OS.ID)
		}

		basic.Metrics = metrics

		tmpl := template.Must(template.ParseFiles("templates/basic.html"))
		tmpl.Execute(w, basic)
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

func getIcon(osId string) string {
	iconMap := map[string]string{
		"ubuntu": "/static/icons/ubuntu.svg",
		"debian": "/static/icons/debian.svg",
		"fedora": "/static/icons/fedora.svg",
		"arch":   "/static/icons/arch.svg",
		"alpine": "/static/icons/alpine.svg",
		"pop":    "/static/icons/pop-os.svg", // <- add this
	}

	return iconMap[osId]
}
