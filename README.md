# Akvelon Software Audit
Scalable compliance and security audit for modern development

# Development

### Build locally
Execute ```make build``` from the root.

Run ```./myapp ```t o start service locally.

### Run tests
Execute ```make test``` from the root.

### CLI
Added now for testing purposes.

Example usage: 
``` 
    cd cmd &&  ./cmd  -repo https://github.com/akvelon/PowerBI-Stacked-Column-Chart
```

### Local Development With Docker
1) Clean root folder from build artifacts: 
```
    make clean
```
2) Build image locally: 
```
    docker build -t akvelon-software-audit .
```
3) Run Docker container: 
```
    docker run -p 777:777 
        -v $(pwd)/views:/app/views:rw 
        -v $(pwd)/controllers:/app/controllers:rw 
        -v $(pwd)/internals:/app/internals:rw  
        -d akvelon-software-audit
```
4) Expose service at 
```
    http://localhost:777
``` 
and edit source code whenever needed.


