# Akvelon Software Audit
Scalable compliance and security audit for modern development

# Development

### Build locally
Execute ```make build``` from the root.

Run `./myapp` to start service locally.

### Run tests
Execute ```make test``` from the root.

### CLI
Added now for testing purposes.

Example usage: ``` cd cmd &&  ./cmd  -repo https://github.com/akvelon/PowerBI-Stacked-Column-Chart```

### Docker
1) Build image locally: ```docker build -t akvelon-software-audit .```
2) Run Docker container: ```docker run -p 777:777 akvelon-software-audit ```
3) Expose service at ```http://localhost:777```