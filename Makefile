build:
	go build -o myapp . && cd cmd && go build .

start:
	akvelon-software-audit

clean:
	rm myapp && rm ./cmd/cmd
