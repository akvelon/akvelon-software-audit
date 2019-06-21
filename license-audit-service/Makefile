build:
	go build -o audit-srv cmd/server/main.go && go build -o audit-worker cmd/worker/main.go

start:
	./audit-srv & ./audit-worker

clean:
	rm -f audit-srv && rm -f audit-worker && rm -rf _repos && rm -f akvelonaudit.db

cleanAllRepos:
	cd _repos && rm -rf github.com/ 
