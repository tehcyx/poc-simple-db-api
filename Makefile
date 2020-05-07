.PHONY: build

VERSION:=$(if ${VERSION},${VERSION},$(shell cat VERSION))
GIT_COMMIT:=$(if ${GIT_COMMIT},${GIT_COMMIT},$(shell git rev-list -1 HEAD))
CMP_NAME=simple-db-api
FRONTEND_NAME=simple-db-api-frontend
DATABASE=postgres

build: bin/templates test cover backend frontend

backend: bin/templates test cover bench
	go build -ldflags "-X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.Version=${VERSION} -X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/app ./cmd/${CMP_NAME}

frontend: bin/templates test cover bench
	go build -ldflags "-X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.Version=${VERSION} -X github.com/tehcyx/${CMP_NAME}/pkg/${CMP_NAME}/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/frontend ./cmd/${FRONTEND_NAME}

docker: dockerbackend dockerfrontend dockerdb

dockerbackend:
	docker build -t ${CMP_NAME}:${VERSION} -f build/package/Dockerfile.backend --build-arg CMP_NAME="${CMP_NAME}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" .
	docker tag ${CMP_NAME}:${VERSION} ${CMP_NAME}:${GIT_COMMIT}

dockerfrontend:
	docker build -t ${FRONTEND_NAME}:${VERSION} -f build/package/Dockerfile.frontend --build-arg CMP_NAME="${FRONTEND_NAME}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" .
	docker tag ${FRONTEND_NAME}:${VERSION} ${FRONTEND_NAME}:${GIT_COMMIT}

dockerdb:
	docker build -t ${DATABASE}:${VERSION} -f build/package/Dockerfile.${DATABASE} --build-arg CMP_NAME="${DATABASE}" --build-arg VERSION="${VERSION}" --build-arg GIT_COMMIT="${GIT_COMMIT}" ./${DATABASE}
	docker tag ${DATABASE}:${VERSION} ${DATABASE}:${GIT_COMMIT}

tag: tagbackend tagfrontend tagdb

tagbackend: dockerbackend
	docker tag ${CMP_NAME}:${VERSION} tehcyx/${CMP_NAME}:${VERSION}
	docker tag ${CMP_NAME}:${VERSION} tehcyx/${CMP_NAME}:${GIT_COMMIT}

tagfrontend: dockerfrontend
	docker tag ${FRONTEND_NAME}:${VERSION} tehcyx/${FRONTEND_NAME}:${VERSION}
	docker tag ${FRONTEND_NAME}:${VERSION} tehcyx/${FRONTEND_NAME}:${GIT_COMMIT}

tagdb: dockerdb
	docker tag ${DATABASE}:${VERSION} tehcyx/${DATABASE}:${VERSION}
	docker tag ${DATABASE}:${VERSION} tehcyx/${DATABASE}:${GIT_COMMIT}

push: pushbackend pushfrontend pushdb

pushbackend: tagbackend
	docker push tehcyx/${CMP_NAME}:${VERSION}
	docker push tehcyx/${CMP_NAME}:${GIT_COMMIT}

pushfrontend: tagfrontend
	docker push tehcyx/${FRONTEND_NAME}:${VERSION}
	docker push tehcyx/${FRONTEND_NAME}:${GIT_COMMIT}

pushdb: tagdb
	docker push tehcyx/${DATABASE}:${VERSION}
	docker push tehcyx/${DATABASE}:${GIT_COMMIT}

install: build
	go install

bin/templates:
	mkdir -p internal/tmpl
	go run hack/packtemplates.go
	go fmt github.com/tehcyx/${CMP_NAME}/internal/tmpl

test:
	go test -v ./...

test-race:
	go test -v ./... -race

test-unit:
	go test -v ./... -race 2>&1 | go run github.com/jstemmer/go-junit-report > report_test.xml

cover:
	go test -v ./... -cover

cover-race:
	go test -v ./... -cover -race

cover-unit:
	go test -v ./... -cover -race 2>&1 | go run github.com/jstemmer/go-junit-report > report_cover.xml

bench:
	go test -v -bench=. ./...

bench-race:
	go test -v -bench=. ./... -race

bench-unit:
	go test -v -bench=. ./... -race 2>&1 | go run github.com/jstemmer/go-junit-report > report_bench.xml

clean:
	rm -rf bin
