# UI for interacting with Akvelon Audit Servie

### Local Development With Docker

1) Clean root folder from build artifacts: 
```
    make clean
```
2) Build image locally: 
```
    docker build -t akv-audit-ux .
```
3) Run Docker container: 
```
    docker run -p 777:777 -v $(pwd)/views:/app/views -v $(pwd)/conf:/app/conf -d akv-audit-ux
```
4) Expose service at 
```
    http://localhost:777
``` 
and edit source code whenever needed.
