# Akvelon Software Audit
Scalable compliance and security audit for modern development.

* license-audit-service - REST service for audit management of github repositories
* ux - UI for interacting with Akvelon Audit Service

### Running With Docker

Run 
```
docker-compose up 
```
to start services. Browse to http://localhost:777 to access portal.

Run
```
docker-compose down
```
to stop services.

Navigate to http://localhost:15672/#/queues to inspect queue's state. 