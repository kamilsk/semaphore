---
code: null
id: MDU6SXNzdWUxOTc1MTM3NTg=
databaseId: 197513758
number: 46
url: https://github.com/kamilsk/semaphore/issues/46
title: "improve code quality"
issueType: null
labels:
  - "type:refactoring"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-12-25T20:11:41Z
updatedAt: 2026-11-08T14:08:14Z
lastEditedAt: 2026-11-08T14:08:14Z
closedAt: 2016-12-26T15:40:08Z
---

# improve code quality

```
make docker-tool-gometalinter DEADLINE=60s
docker run --rm \
                   -v "/Users/.../Development/go/src/github.com/kamilsk/semaphore":"/go/src/github.com/kamilsk/semaphore" \
                   -w "/go/src/github.com/kamilsk/semaphore" \
                   kamilsk/go-tools:latest \
                   gometalinter.v1 --vendor --deadline=60s ./...
go17_test.go:93::warning: duplicate of pre_go17_test.go:94-138 (dupl)
pre_go17_test.go:94::warning: duplicate of go17_test.go:93-137 (dupl)
simple/simple.go:7:7:error: could not import github.com/kamilsk/semaphore (can't find import: github.com/kamilsk/semaphore) (gotype)
simple/simple.go:12:2:error: undeclared name: intf (gotype)
simple/simple.go:13:2:error: undeclared name: intf (gotype)
simple/example_test.go:56:25:warning: error return value not checked (defer resp.Body.Close()) (errcheck)
example_test.go:90:25:warning: error return value not checked (defer resp.Body.Close()) (errcheck)
make: *** [docker-tool-gometalinter] Error 1
```
