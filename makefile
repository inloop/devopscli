OWNER=inloopeu
IMAGE_NAME=devops
QNAME=$(OWNER)/$(IMAGE_NAME)

GIT_TAG=$(QNAME):$(TRAVIS_COMMIT)
BUILD_TAG=$(QNAME):0.1.$(TRAVIS_BUILD_NUMBER)
LATEST_TAG=$(QNAME):latest

lint:
	docker run -it --rm -v "$(PWD)/Dockerfile:/Dockerfile:ro" redcoolbeans/dockerlint

build:
	go get ./...
	# gox -osarch="linux/amd64" -output="bin/devops-alpine"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/devops-alpine .
	docker build -t $(GIT_TAG) .
	docker build -t $(GIT_TAG)-golang ./Dockerfile.golang

tag:
	docker tag $(GIT_TAG) $(BUILD_TAG)
	docker tag $(GIT_TAG)-golang $(BUILD_TAG)-golang
	docker tag $(GIT_TAG) $(LATEST_TAG)
	docker tag $(GIT_TAG)-golang $(LATEST_TAG)-golang

login:
	@docker login -u "$(DOCKER_USER)" -p "$(DOCKER_PASS)"
push: login
	# docker push $(GIT_TAG)
	# docker push $(BUILD_TAG)
	docker push $(LATEST_TAG)
	docker push $(LATEST_TAG)-golang


build-local:
	go get ./...
	go build -o devops

deploy-local:
	make build-local
	mv devops /usr/local/bin/
