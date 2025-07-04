package handlers

import (
	"encoding/json"
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics.GetDisksMetrics())
	})

	router.Get("/disks", func(w http.ResponseWriter, r *http.Request) {
		metrics := metrics.GetDisksMetrics()
		tmpl := template.Must(template.ParseFiles("templates/disks.html"))
		tmpl.Execute(w, metrics)
	})

	return router
}

func getIcon(osId string) string {
	iconMap := map[string]string{
		"ubuntu":                   "/static/icons/ubuntu.svg",
		"pop":                      "/static/icons/pop-os.svg",
		"pop-os":                   "/static/icons/pop-os.svg",
		"popos":                    "/static/icons/popos.svg",
		"debian":                   "/static/icons/debian.svg",
		"arch":                     "/static/icons/arch.svg",
		"archlinux":                "/static/icons/archlinux.svg",
		"fedora":                   "/static/icons/fedora.svg",
		"manjaro":                  "/static/icons/manjaro.svg",
		"linuxmint":                "/static/icons/linuxmint.svg",
		"linuxmint-cinnamon":       "/static/icons/linuxmint-cinnamon.svg",
		"kubuntu":                  "/static/icons/kubuntu.svg",
		"xubuntu":                  "/static/icons/xubuntu.svg",
		"zorin":                    "/static/icons/zorin.svg",
		"elementary":               "/static/icons/elementary.svg",
		"opensuse":                 "/static/icons/opensuse.svg",
		"void":                     "/static/icons/void.svg",
		"gentoo":                   "/static/icons/gentoo.svg",
		"alpine":                   "/static/icons/alpine.svg",
		"centos":                   "/static/icons/centos-stream.svg",
		"nixos":                    "/static/icons/nixos.svg",
		"freebsd":                  "/static/icons/freebsd.svg",
		"windows 7":                "/static/icons/windows-7.webp",
		"windows vista":            "/static/icons/windows-7.webp",
		"windows server 2008":      "/static/icons/windows-7.webp",
		"windows server 2008 R2":   "/static/icons/windows-7.webp",
		"windows 8":                "/static/icons/windows-8.png",
		"windows server 2012":      "/static/icons/windows-8.png",
		"windows server 2012 R2":   "/static/icons/windows-8.png",
		"windows 8.1":              "/static/icons/windows-8.png",
		"windows 10":               "/static/icons/windows-10.webp",
		"windows server 2016/2019": "/static/icons/windows-10.webp",
		"windows 11":               "/static/icons/windows-11.png",
		"windows server 2022":      "/static/icons/windows-11.png",
		"unknown":                  "/static/icons/unknown.png",
	}

	if icon, ok := iconMap[osId]; ok {
		return icon
	}

	return "/static/icons/unknown.png"
}
