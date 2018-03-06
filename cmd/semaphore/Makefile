.PHONY: cmd-install
cmd-install:
	if ! [ -e main.go ]; then go install ./cmd/semaphore; else go install .; fi



.PHONY: cmd-test-1
cmd-test-1:
	docker run --rm -it \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           -e GO15VENDOREXPERIMENT=1 \
	           golang:1.5 \
	           /bin/sh -c 'go install ./cmd/semaphore; \
	                       semaphore create 1; \
	                       semaphore add -- curl localhost; \
	                       semaphore add -- curl example.com; \
	                       semaphore add -- curl unknown; \
	                       semaphore add -- curl example.com; \
	                       semaphore add -- curl localhost; \
	                       cat /tmp/semaphore.json && echo ""; \
	                       semaphore wait --notify --speed=300 --timeout=10s'

.PHONY: cmd-test-2
cmd-test-2:
	docker run --rm -it \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           -e GO15VENDOREXPERIMENT=1 \
	           golang:1.5 \
	           /bin/sh -c 'go install ./cmd/semaphore \
	                       && semaphore help \
	                       && semaphore -h \
	                       && semaphore --help'

.PHONY: cmd-test-3
cmd-test-3:
	docker run --rm -it \
	           -v '${GOPATH}/src/${GO_PACKAGE}':'/go/src/${GO_PACKAGE}' \
	           -w '/go/src/${GO_PACKAGE}' \
	           -e GO15VENDOREXPERIMENT=1 \
	           golang:1.5 \
	           /bin/sh -c 'go install ./cmd/semaphore; \
	                       semaphore wait --debug --filename="/bad"'



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

.PHONY: cmd-test-2-local
cmd-test-2-local: cmd-install
cmd-test-2-local:
	'${GOPATH}'/bin/semaphore help
	'${GOPATH}'/bin/semaphore -h
	'${GOPATH}'/bin/semaphore --help

.PHONY: cmd-test-3-local
cmd-test-3-local: cmd-install
cmd-test-3-local:
	'${GOPATH}'/bin/semaphore wait --debug --filename="/bad"
