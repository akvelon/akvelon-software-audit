build:
	go build -o myapp . && cd cmd && go build .

test:
	go test -cover ./internals/vcs

start:
	bee run

clean:
	rm -f myapp && rm -f akvelon-software-audit && rm -f ./cmd/cmd && rm -rf _repos && rm -f akvelonaudit.db

cleanGHRepos:
	cd _repos && rm -rf github.com/ 
