build:
	go build -o audit-serv cmd/server/main.go

start:
	./audit-serv

clean:
	rm -f audit-serv && rm -rf _repos && rm -f akvelonaudit.db

cleanAllRepos:
	cd _repos && rm -rf github.com/ 
