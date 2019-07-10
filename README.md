# Akvelon Software Audit
Scalable compliance and security audit for modern development.

* license-audit-service - REST service for audit management of github repositories
* license-audit-worker - console app for executing long running audit tasks
* ux - UI for interacting with Akvelon Audit Service (powered by [Beego](https://beego.me/) framework)
* [MongoDB](https://www.mongodb.com/) for persistant storage

We use [Prometheus](https://prometheus.io/) for monitoring metrics and [RabbitMQ](https://www.rabbitmq.com/) as a message broker for communication among services.

Also we use [Jaeger](https://www.jaegertracing.io/) as an end-to-end distributed tracing tool.

### Running With Docker

Run `docker-compose up` to start services. Browse to http://localhost:777 to access portal.

Run `docker-compose down` to stop services.

### Monitoring Data

* Navigate to http://localhost:15672/#/queues to inspect RabbitMQ queue's state. 

* Navigate to http://localhost:9090/graph to inspect Prometheus graph data (e.g. current number of __Go Threads__, __audit_ux_http_requests_total__ count etc).

* Navigate to http://localhost:16686/search to inspect Jaeger traces.
