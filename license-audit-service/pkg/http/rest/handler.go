package rest

import (
	"akvelon/akvelon-software-audit/license-audit-service/pkg/licanalize"
	"akvelon/akvelon-software-audit/license-audit-service/pkg/monitor"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
)

var (
	rabbitSrv = "amqp://guest:guest@rabbitmq:5672"
	uXAuditQueueName = "audit-queue"
) 

// Handler handles request using service injected.
func Handler(a licanalize.Service, m *monitor.Monitor, t opentracing.Tracer) http.Handler {
	log.Println("Register monitor...")
	m.RegisterMonitor()
	router := httprouter.New()

	router.Handler("GET", "/metrics", promhttp.Handler())
	router.GET("/health", checkHealth(a, m, t))
	router.GET("/recent", getRecentResults(a, m, t))
	router.GET("/analize", getAnalizedResult(a, m, t))

	router.POST("/analize", analize(a, m, t))


	return router
}

func checkHealth(a licanalize.Service, m *monitor.Monitor, t opentracing.Tracer) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func getRecentResults(a licanalize.Service, m *monitor.Monitor, t opentracing.Tracer) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()

		spanCtx, _ := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := t.StartSpan("get-recent-results", ext.RPCServerOption(spanCtx))
		defer span.Finish()
		
		log.Println("Start exec getRecentResults...")
		recent, err := a.GetRecent()
		if err != nil {
			span.LogKV("getting-recent-results", "failed to get recent results from audit db")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recent)
	}
}

func getAnalizedResult(a licanalize.Service, m *monitor.Monitor, t opentracing.Tracer) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()

		spanCtx, _ := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := t.StartSpan("get-analized-repo-result", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		log.Println("Start exec getAnalizedResult...")
		queryValues := r.URL.Query()
		url := queryValues.Get("url")
		log.Printf("url: %s", url)
		result, err := a.GetRepoResultFromDB(url)
		if err != nil {
			span.LogKV("getting-analized-repo-result", "failed to get result from audit db")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func analize(a licanalize.Service, m *monitor.Monitor, t opentracing.Tracer) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m.GetHttpRequestsTotal().Inc()
		log.Println("Start exec analize...")

		spanCtx, _ := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := t.StartSpan("post-analize-repo", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		repoLink := r.FormValue("url")
		if repoLink == "" {
			http.Error(w, "post-analize-repo: Failed to parse input parameter, url is missing", http.StatusBadRequest)
			return
		}

		conn, err := amqp.Dial(rabbitSrv)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a rmq channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			uXAuditQueueName, // name
			false,            // durable
			false,            // delete when unused
			false,            // exclusive
			false,            // no-wait
			nil,              // arguments
		)
		failOnError(err, "Failed to declare a rabbit queue")

		body := repoLink
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fmt.Sprintf("In progress analizing repo %s", repoLink))
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
