GIT_ORIGIN:="git@github.com:kamilsk/semaphore.git"
GIT_MIRROR:="git@bitbucket.org:kamilsk/semaphore.git"
GO_PACKAGE:="github.com/kamilsk/semaphore"

include makes/env.mk
include makes/deps.mk
include makes/docker.mk
include makes/flow.mk
include makes/tests.mk
include makes/tools.mk

.PHONY: all
all: install-deps build install

.PHONY: docker-bench
docker-bench: docker-bench-1.5-alpine-gcc docker-bench-1.6-alpine-gcc docker-bench-1.7-alpine-gcc docker-bench-alpine-gcc

.PHONY: docker-test
docker-test: docker-test-1.5-alpine-gcc docker-test-1.6-alpine-gcc docker-test-1.7-alpine-gcc docker-test-alpine-gcc
