GIT_ORIGIN?="git@github.com:kamilsk/semaphore.git"
GIT_MIRROR?="git@bitbucket.org:kamilsk/semaphore.git"
GO_TEST_COVERAGE_MODE?="count"
GO_TEST_COVERAGE_FILE_NAME?="coverage.out"
GOFMT_FLAGS?="-s"
GOLINT_MIN_CONFIDENCE?="0.3"


.PHONY: all build
.PHONY: install install-deps
.PHONY: update-deps
.PHONY: test test-with-coverage test-with-coverage-formatted test-with-coverage-profile
.PHONY: clean vet
.PHONY: dev publish

# work in progress

GO_PACKAGE:=github.com/kamilsk/semaphore

include makes/env.mk
include makes/deps.mk
include makes/docker.mk
include makes/flow.mk
include makes/tests.mk
include makes/tools.mk

.PHONY: docker-bench
docker-bench: docker-bench-1.5-alpine-gcc docker-bench-1.6-alpine-gcc docker-bench-1.7-alpine-gcc docker-bench-alpine-gcc

.PHONY: docker-test
docker-test: docker-test-1.5-alpine-gcc docker-test-1.6-alpine-gcc docker-test-1.7-alpine-gcc docker-test-alpine-gcc

# end

all: install-deps build install

build:
	go build -v ./...

install:
	go install ./...

install-deps:
	go get -d -t ./...

update-deps:
	go get -d -t -u ./...

test:
	go test -race -v ./...

test-with-coverage:
	go test -cover -race ./...

test-with-coverage-formatted:
	go test -cover -race ./... | column -t | sort -r

test-with-coverage-profile:
	echo "mode: ${GO_TEST_COVERAGE_MODE}" > ${GO_TEST_COVERAGE_FILE_NAME}
	for package in $$(go list ./...); do \
		go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile "coverage_$${package##*/}.out" -race "$${package}"; \
		sed '1d' "coverage_$${package##*/}.out" >> ${GO_TEST_COVERAGE_FILE_NAME}; \
	done

clean:
	go clean -i -x ./...

vet:
	go vet ./...

dev:
	git remote set-url origin $(GIT_ORIGIN)
	git remote add mirror $(GIT_MIRROR)

publish:
	git push origin master --tags
	git push mirror master --tags
