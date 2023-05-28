---
code: null
id: MDU6SXNzdWUyODIzMzM4MDY=
databaseId: 282333806
number: 126
url: https://github.com/kamilsk/semaphore/issues/126
title: "WithContext should return ctx.Done channel"
issueType: null
labels:
  - "type:task"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-15T06:39:53Z
updatedAt: 2019-02-26T06:57:53Z
lastEditedAt: null
closedAt: 2019-02-26T06:57:53Z
---

# WithContext should return ctx.Done channel

```go
semaphore.Acquire(semaphore.Multiplex(semaphore.WithContext(req.Context()), semaphore.WithSignal(os.Interrupt)))
```

like as errgroup:

```go
g, ctx := errgroup.WithContext(ctx)
```

https://github.com/golang/sync/tree/master/errgroup

<!-- 2019-02-26T06:57Z https://github.com/kamilsk/semaphore/issues/126#issuecomment-467319602
will be removed by https://github.com/kamilsk/semaphore/issues/133
-->
