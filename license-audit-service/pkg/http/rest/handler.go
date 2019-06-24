package rest

import (
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/monitor"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/julienschmidt/httprouter"
)

// Handler handles request using service injected.
func Handler(a licanalize.Service, m *monitor.Monitor) http.Handler {
	log.Println("Register monitor...")
	m.RegisterMonitor()
	router := httprouter.New()

	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/health", checkHealth(a, m))
	router.GET("/recent", getRecentResults(a, m))
	router.GET("/analize", getAnalizedResult(a, m))

	router.POST("/analize", analize(a, m))


	return router
}

func checkHealth(a licanalize.Service, m *monitor.Monitor) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()

		log.Println("Start exec checkHealth...")
		w.Header().Set("Content-Type", "application/json")
		// TODO: check if DB is avaliable too
		if a.CheckHealth() {
			json.NewEncoder(w).Encode("Healthy")
		} else {
			json.NewEncoder(w).Encode("Unhealthy")
		}
	}
}

func getRecentResults(a licanalize.Service, m *monitor.Monitor) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()

		log.Println("Start exec getRecentResults...")
		recent, err := a.GetRecent()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recent)
	}
}

func getAnalizedResult(a licanalize.Service, m *monitor.Monitor) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()

		log.Println("Start exec getAnalizedResult...")
		queryValues := r.URL.Query()
		url := queryValues.Get("url")
		log.Printf("url: %s", url)
		result, err := a.GetRepoResultFromDB(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func analize(a licanalize.Service, m *monitor.Monitor) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()

		log.Println("Start exec analize...")
		repoLink := r.FormValue("url")
		if repoLink == "" {
			http.Error(w, "Failed to parse input parameter, url is missing", http.StatusBadRequest)
			return
		}

		err := a.Scan(licanalize.AnalizedRepo{
			URL: repoLink,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fmt.Sprintf("Finished analizing repo %s", repoLink))
	}
}
