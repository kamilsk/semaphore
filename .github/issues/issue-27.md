---
code: null
id: MDU6SXNzdWUxOTIxMjA2MTQ=
databaseId: 192120614
number: 27
url: https://github.com/kamilsk/semaphore/issues/27
title: "separate Releaser"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-11-28T21:22:44Z
updatedAt: 2016-12-03T20:15:27Z
lastEditedAt: null
closedAt: 2016-12-03T20:15:27Z
---

# separate Releaser

the main reason is to don't forget Release() semaphore.

```go
if releaser, err := sem.Acquire(ctx); err != {
   // forget
}
defer releaser.Release()
```
