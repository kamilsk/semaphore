OPEN_BROWSER       =
SUPPORTED_VERSIONS = 1.5 1.6 1.7 1.8 1.9 latest

include makes/env.mk
include makes/local.mk
include makes/docker.mk

.PHONY: check-code-quality
check-code-quality: ARGS = \
	--exclude='.*_test\.go:.*error return value not checked.*\(errcheck\)$' \
	--exclude='duplicate of.*_test.go.*\(dupl\)$' \
	--vendor --deadline=2m ./...
check-code-quality: docker-tool-gometalinter



.PHONY: cmd-deps
cmd-deps:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}/cmd/semaphore' \
	           kamilsk/go-tools:latest \
	           glide install -v
	rm -rf cmd/semaphore/.glide

.PHONY: cmd-test-1
cmd-test-1:
	docker run --rm -it \
	           -v '$(GOPATH)/src/$(GO_PACKAGE)':'/go/src/$(GO_PACKAGE)' \
	           -w '/go/src/$(GO_PACKAGE)' \
	           golang:1.8 \
	           /bin/sh -c 'go install ./cmd/semaphore; \
	                       semaphore create 1; \
	                       semaphore add -- curl localhost; \
	                       semaphore add -- curl example.com; \
	                       semaphore add -- curl unknown; \
	                       semaphore add -- curl example.com; \
	                       semaphore add -- curl localhost; \
	                       cat /tmp/semaphore.json && echo ""; \
	                       semaphore wait --notify --speed=300 --timeout=10s; \
	                       echo $$?'

.PHONY: cmd-test-2
cmd-test-2:
	docker run --rm -it \
	           -v '$(GOPATH)/src/$(GO_PACKAGE)':'/go/src/$(GO_PACKAGE)' \
	           -w '/go/src/$(GO_PACKAGE)' \
	           golang:1.8 \
	           /bin/sh -c 'go install ./cmd/semaphore \
	                       && semaphore help \
	                       && semaphore -h \
	                       && semaphore --help; \
	                       echo $$?'

.PHONY: cmd-deps-local
cmd-deps-local:
	cd cmd/semaphore && glide install -v

.PHONY: cmd-install
cmd-install:
	go install ./cmd/semaphore

.PHONY: cmd-test-1-local
cmd-test-1-local: cmd-install
cmd-test-1-local:
	semaphore create 1
	semaphore add -- curl localhost
	semaphore add -- curl example.com
	semaphore add -- curl unknown
	semaphore add -- curl example.com
	semaphore add -- curl localhost
	semaphore wait --notify --speed=300 --timeout=10s
	echo $$?

.PHONY: cmd-test-2-local
cmd-test-2-local: cmd-install
cmd-test-2-local:
	semaphore help
	semaphore -h
	semaphore --help
	echo $$?



.PHONY: pull-github-tpl
pull-github-tpl:
	rm -rf .github
	(git clone git@github.com:kamilsk/shared.git .github && cd .github && git checkout github-tpl-go-v1 \
	  && echo 'github templates at revision' $$(git rev-parse HEAD) && rm -rf .git)
	rm .github/README.md

.PHONY: pull-makes
pull-makes:
	rm -rf makes
	(git clone git@github.com:kamilsk/shared.git makes && cd makes && git checkout makefile-go-v1 \
	  && echo 'makes at revision' $$(git rev-parse HEAD) && rm -rf .git)
