# License Audit Service
Aims to execute compliance and security audit for provided github repository

# API
```
GET /health => health check
GET /recent => get 10 recent scans (returns repo names)
GET /analize?url=github.com/akvelon/PowerBI-Stacked-Column-Chart => returns JSON with scan results

POST /analize => execute scan for given url (form-data; name="url") and store results in DB

```

# Development

## Docker

1) Clean current folder from build artifacts: 
```
    make clean
```

2) Build Docker Image
```
    docker build -t akv-audit-srv .   
```

3) Run Docker container: 
```
    docker run --rm -d -p 8000:8000 akv-audit-srv:latest
```

4) Test at http://localhost:8000