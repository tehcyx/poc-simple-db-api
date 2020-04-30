.PHONY: build

VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD)
CMP_NAME=simple-db-api
FRONTEND_NAME=simple-db-api-frontend

build: bin/templates test cover
	go build -ldflags "-X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.Version=${VERSION} -X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/app ./cmd/${CMP_NAME}
	go build -ldflags "-X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.Version=${VERSION} -X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/frontend ./cmd/${FRONTEND_NAME}

docker: bin/templates test cover 
	docker build -t $(CMP_NAME):$(VERSION) -f build/package/Dockerfile --build-arg CMP_NAME="${CMP_NAME}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" ../../
	docker build -t $(CMP_NAME):$(GIT_COMMIT) -f build/package/Dockerfile --build-arg CMP_NAME="${CMP_NAME}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" ../../

install: build
	go install

bin/templates:
	mkdir -p internal/tmpl
	go run hack/packtemplates.go
	go fmt github.com/tehcyx/${CMP_NAME}/internal/tmpl

test:
	go test ./...

test-race:
	go test ./... -race

cover:
	go test ./... -cover

cover-race:
	go test ./... -cover -race

clean:
	rm -rf bin
