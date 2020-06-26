# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Module
CMD_PATH=github.com/genblue-private/cedrus-backend/cmd

# Cedrus service
CEDRUS_SERVICE=cedrusservice
CEDRUS_PATH=$(CMD_PATH)/$(CEDRUS_SERVICE)
CEDRUS_SERVICE_CMD=./cmd/$(CEDRUS_SERVICE)
CEDRUS_SERVICE_BIN=./bin/$(CEDRUS_SERVICE)

# Email service
EMAIL_SERVICE=emailservice
EMAIL_PATH=$(CMD_PATH)/$(EMAIL_SERVICE)
EMAIL_SERVICE_CMD=./cmd/$(EMAIL_SERVICE)
EMAIL_SERVICE_BIN=./bin/$(EMAIL_SERVICE)

# Seth service
SETH_SERVICE=sethservice
SETH_PATH=$(CMD_PATH)/$(SETH_SERVICE)
SETH_SERVICE_CMD=./cmd/$(SETH_SERVICE)
SETH_SERVICE_BIN=./bin/$(SETH_SERVICE)

.PHONY: all build build-unix test clean clean-cedrus run deps docker-images docker-tag docs

all: test build
build:
	$(GOBUILD) -o $(CEDRUS_SERVICE_BIN) $(CEDRUS_SERVICE_CMD)
	$(GOBUILD) -o $(EMAIL_SERVICE_BIN) $(EMAIL_SERVICE_CMD)
	$(GOBUILD) -o $(SETH_SERVICE_BIN) $(SETH_SERVICE_CMD)
build-unix:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(CEDRUS_SERVICE_BIN) $(CEDRUS_SERVICE_CMD) && chmod +x $(CEDRUS_SERVICE_BIN)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(EMAIL_SERVICE_BIN) $(EMAIL_SERVICE_CMD) && chmod +x $(EMAIL_SERVICE_BIN)
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(SETH_SERVICE_BIN) $(SETH_SERVICE_CMD) && chmod +x $(SETH_SERVICE_BIN)
test:
	set -o pipefail && $(GOTEST) ./... -v | grep -v -E 'GET|POST|PUT'
test-ci:
	$(GOTEST) ./... -v
clean : clean-cedrus
clean-cedrus:
	$(GOCLEAN) $(CEDRUS_PATH)
	rm -f $(CEDRUS_SERVICE_BIN)
	rm -f $(EMAIL_SERVICE_BIN)
	rm -f $(SETH_SERVICE_BIN)
build-docker: build-unix
	docker-compose build
run:
	docker-compose up -d mongo $(CEDRUS_SERVICE) $(SETH_SERVICE)
	while true ; do docker-compose rm -fsv $(EMAIL_SERVICE); docker-compose up $(EMAIL_SERVICE) & sleep 10; done
stop:
	docker-compose down
	if [ $$(docker ps -a | grep service) ]; then docker stop $$(docker ps -a -q); fi
	if [ $$(docker ps -a | grep service) ]; then docker rm $$(docker ps -a -q); fi
deps:
	$(GOGET) -v -d ./...
docker-images: build-unix
	docker build -t $(CEDRUS_SERVICE) -f build/package/$(CEDRUS_SERVICE).Dockerfile .
	docker build -t $(EMAIL_SERVICE) -f build/package/$(EMAIL_SERVICE).Dockerfile .
	docker build -t $(SETH_SERVICE) -f build/package/$(SETH_SERVICE)-multistage.Dockerfile .
docker-tag:
	docker tag $(CEDRUS_SERVICE) ${REGISTRY}:$(CEDRUS_SERVICE)-dev
	docker tag $(EMAIL_SERVICE) ${REGISTRY}:$(EMAIL_SERVICE)-dev
	docker tag $(SETH_SERVICE) ${REGISTRY}:$(SETH_SERVICE)-dev
docker-push:
	echo ${PASSWORD} | docker login -u ${USERNAME} --password-stdin
	docker push ${REGISTRY}:$(CEDRUS_SERVICE)-dev
	docker push ${REGISTRY}:$(EMAIL_SERVICE)-dev
	docker push ${REGISTRY}:$(SETH_SERVICE)-dev
cf-deploy:
	cf push $(CEDRUS_SERVICE) --docker-image ${REGISTRY}:$(CEDRUS_SERVICE)-dev --docker-username ${USERNAME}
	cf push $(SETH_SERVICE) --docker-image ${REGISTRY}:$(SETH_SERVICE)-dev --docker-username ${USERNAME}
	cf push $(EMAIL_SERVICE) --docker-image ${REGISTRY}:$(EMAIL_SERVICE)-dev --docker-username ${USERNAME}
docs:
	if ! which swag; then go get -u github.com/swaggo/swag/cmd/swag ; fi
	swag init --generalInfo rest_controller.go --dir internal/$(CEDRUS_SERVICE)/infrastructure/rest --output api-docs/$(CEDRUS_SERVICE)
	swag init --generalInfo rest_controller.go --dir internal/$(SETH_SERVICE)/infrastructure/rest --output api-docs/$(SETH_SERVICE)
