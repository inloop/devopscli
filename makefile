build:
	go get ./...
	go build -o devops

deploy-local:
	make build-local
	mv devops /usr/local/bin/
