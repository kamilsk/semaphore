---
code: null
id: MDU6SXNzdWUyODIzMzU5NjI=
databaseId: 282335962
number: 127
url: https://github.com/kamilsk/semaphore/issues/127
title: "add streaming"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-15T06:53:56Z
updatedAt: 2019-02-26T06:58:37Z
lastEditedAt: null
closedAt: 2019-02-26T06:58:37Z
---

# add streaming

```go
semaphore wait --stream-to=12345
```

semaphore starts listen 12345 TCP port to stream stdout/stderr of executed commands to it.

<!-- 2017-12-15T06:54Z https://github.com/kamilsk/semaphore/issues/127#issuecomment-351928651
can be useful https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
-->

<!-- 2018-04-01T18:00Z https://github.com/kamilsk/semaphore/issues/127#issuecomment-377804882
some inspiration:
- https://github.com/wagoodman/bashful
-->

<!-- 2019-02-26T06:58Z https://github.com/kamilsk/semaphore/issues/127#issuecomment-467319767
moved to https://github.com/kamilsk/semaphore.cli/issues/4
-->
