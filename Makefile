SCRIPT_DIR := $(shell pwd)
TMPDIR := ${SCRIPT_DIR}/.tmp
PROJECT_NAME = ports
DOCKER_NAMESPACE ?= informalict
API_PORT ?= 8080

SRC = $(shell find $(SCRIPT_DIR) -name '*.go' -not -path './test/*')

# linter takes care about fmt, imports and many other checks.
.PHONY: linter
linter:
	@${TMPDIR}/bin/golangci-lint run ./...

# download tools required for this project. It should be done once.
.PHONY: tools
tools:
	@mkdir -p "${TMPDIR}"
	@echo ">> Fetching golangci-lint linter"
	@GOBIN=${TMPDIR}/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

# build an application and create final image.
.PHONY: docker
docker: $(SRC)
	docker build -f Dockerfile -t ${DOCKER_NAMESPACE}/${PROJECT_NAME} ${SCRIPT_DIR}

# run ports service in a docker container on host port 8080.
.PHONY: run
run: docker
	docker run --user $(shell id -u):$(shell id -g) \
		--cap-drop=all --memory 200m --cpus "1.0" \
		--name port-svc -t -d -p "${API_PORT}":8080 --rm ${DOCKER_NAMESPACE}/${PROJECT_NAME}

.PHONY: clean
clean:
	docker stop port-svc

.PHONY: run-tests
run-tests:
	API_PORT="${API_PORT}" TEST_FILE="$(shell pwd)/assets/ports.json" go test -v ./test/...

.PHONY: run-unit-tests
run-unit-tests:
	go test -v ./api/... ./pkg/...
