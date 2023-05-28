---
code: null
id: MDU6SXNzdWUxOTM3NjUwNjk=
databaseId: 193765069
number: 35
url: https://github.com/kamilsk/semaphore/issues/35
title: "fix deps bug"
issueType: null
labels:
  - "type:bug"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-12-06T12:53:49Z
updatedAt: 2016-12-06T14:11:41Z
lastEditedAt: 2016-12-06T12:54:01Z
closedAt: 2016-12-06T14:11:41Z
---

# fix deps bug

```
$ make docker-test-1.5-alpine-gcc
pre_go17.go:5:8: cannot find package "golang.org/x/net/context" in any of:
	/usr/local/go/src/golang.org/x/net/context (from $GOROOT)
	/go/src/golang.org/x/net/context (from $GOPATH)
make: *** [docker-test-1.5-alpine-gcc] Error 1
```

<!-- 2016-12-06T13:09Z https://github.com/kamilsk/semaphore/issues/35#issuecomment-265145576
blocked by https://github.com/kamilsk/shared/issues/14
-->
