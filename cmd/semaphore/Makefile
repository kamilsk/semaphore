.PHONY: cmd-deps
cmd-deps:
	docker run --rm \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}/cmd/semaphore' \
	           kamilsk/go-tools:latest \
	           dep ensure

.PHONY: cmd-test-1
cmd-test-1:
	docker run --rm -it \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
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
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           golang:1.8 \
	           /bin/sh -c 'go install ./cmd/semaphore \
	                       && semaphore help \
	                       && semaphore -h \
	                       && semaphore --help; \
	                       echo $$?'

.PHONY: cmd-deps-local
cmd-deps-local:
	if ! [ -e main.go ]; then cd cmd/semaphore && dep ensure; else dep ensure; fi

.PHONY: cmd-install
cmd-install:
	if ! [ -e main.go ]; then go install ./cmd/semaphore; else go install .; fi

.PHONY: cmd-test-1-local
cmd-test-1-local: cmd-install
cmd-test-1-local:
	'${GOPATH}'/bin/semaphore create 1
	'${GOPATH}'/bin/semaphore add -- curl localhost
	'${GOPATH}'/bin/semaphore add -- curl example.com
	'${GOPATH}'/bin/semaphore add -- curl unknown
	'${GOPATH}'/bin/semaphore add -- curl example.com
	'${GOPATH}'/bin/semaphore add -- curl localhost
	'${GOPATH}'/bin/semaphore wait --notify --speed=300 --timeout=10s
	echo $$?

.PHONY: cmd-test-2-local
cmd-test-2-local: cmd-install
cmd-test-2-local:
	'${GOPATH}'/bin/semaphore help
	'${GOPATH}'/bin/semaphore -h
	'${GOPATH}'/bin/semaphore --help
	echo $$?
