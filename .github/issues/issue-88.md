---
code: null
id: MDU6SXNzdWUyNDc0NzA5MDk=
databaseId: 247470909
number: 88
url: https://github.com/kamilsk/semaphore/issues/88
title: "goroutine pool example"
issueType: null
labels:
  - "type:docs"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-08-02T17:51:57Z
updatedAt: 2017-08-30T10:20:38Z
lastEditedAt: null
closedAt: 2017-08-30T10:20:38Z
---

# goroutine pool example

```go
package gpool

func New(size int) *Pool {
    return &Pool{
        work: make(chan func()),
        sem:  make(chan struct{}, size),
    }
}

func (p *Pool) Schedule(task func()) error {
    select {
    case p.work <- task:
    case p.sem <- struct{}{}:
        go p.worker(task)
    }
}

func (p *Pool) worker(task func()) {
    defer func() { <-p.sem }
    for {
        task()
        task = <-p.work
    }
}
```
