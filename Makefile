build:
	go build -o myapp . && cd cmd && go build .

test:
	go test -cover ./internals/vcs

start:
	bee run

clean:
	rm myapp && rm akvelon-software-audit && rm ./cmd/cmd

cleanGHRepos:
	cd _repos && rm -rf github.com/ 
