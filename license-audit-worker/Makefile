build:
	go build -o audit-worker cmd/worker/main.go

start:
	./audit-worker

clean:
	rm -f audit-worker && rm -rf _repos

cleanAllRepos:
	cd _repos && rm -rf github.com/ 
