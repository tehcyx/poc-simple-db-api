.PHONY: build

VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD)
CMP_NAME=simple-db-api
FRONTEND_NAME=simple-db-api-frontend

build: bin/templates test cover backend frontend

backend: bin/templates test cover
	go build -ldflags "-X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.Version=${VERSION} -X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/app ./cmd/${CMP_NAME}

frontend: bin/templates test cover
	go build -ldflags "-X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.Version=${VERSION} -X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/frontend ./cmd/${FRONTEND_NAME}

docker:
	docker build -t $(CMP_NAME):$(VERSION) -f build/package/Dockerfile.backend --build-arg CMP_NAME="${CMP_NAME}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" .
	docker tag $(CMP_NAME):$(VERSION) $(CMP_NAME):$(GIT_COMMIT)
	docker build -t $(FRONTEND_NAME):$(VERSION) -f build/package/Dockerfile.frontend --build-arg CMP_NAME="${FRONTEND_NAME}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" .
	docker tag $(FRONTEND_NAME):$(VERSION) $(FRONTEND_NAME):$(GIT_COMMIT)

tag:
	docker tag $(CMP_NAME):$(VERSION) tehcyx/$(CMP_NAME):$(VERSION)
	docker tag $(CMP_NAME):$(VERSION) tehcyx/$(CMP_NAME):$(GIT_COMMIT)
	docker tag $(FRONTEND_NAME):$(VERSION) tehcyx/$(FRONTEND_NAME):$(VERSION)
	docker tag $(FRONTEND_NAME):$(VERSION) tehcyx/$(FRONTEND_NAME):$(GIT_COMMIT)

push:
	docker push tehcyx/$(CMP_NAME):$(VERSION)
	docker push tehcyx/$(CMP_NAME):$(GIT_COMMIT)
	docker push tehcyx/$(FRONTEND_NAME):$(VERSION)
	docker push tehcyx/$(FRONTEND_NAME):$(GIT_COMMIT)

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
