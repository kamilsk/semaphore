---
code: null
id: MDU6SXNzdWUyNzYyODgxODU=
databaseId: 276288185
number: 117
url: https://github.com/kamilsk/semaphore/issues/117
title: "idea: barrier"
issueType: null
labels:
  - "type:feature"
  - "status:research"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-11-23T07:53:07Z
updatedAt: 2020-12-30T13:32:57Z
lastEditedAt: null
closedAt: 2020-12-30T13:32:57Z
---

# idea: barrier

Resizeable semaphore potentially can have over-engineered implementation. The main objective to create resizeable semaphore is to make graceful degradation, so, the simplest way is to create barrier:

```go
s := semaphore.New(1000)
w := watcher.New()
w.On(event.Overloading, func () { s.SetBarrier(700) }) // capacity of semaphore will be equal to 300
go w.Watch()
```

<!-- 2017-11-24T21:38Z https://github.com/kamilsk/semaphore/issues/117#issuecomment-346898659
possible implementation:

```go
type sem struct {
  s chan struct{}
  b chan struct{}
}

func (s *sem) SetBarrier(b int) {
  // do validation

  // lock is needed?

  for len(s.b) == b {
    s.b <- s.s
  }
}
```
-->

<!-- 2020-12-30T13:32Z https://github.com/kamilsk/semaphore/issues/117#issuecomment-752605885
won't do
-->
