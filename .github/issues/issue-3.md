---
code: null
id: MDU6SXNzdWUxODE4NzYyNTg=
databaseId: 181876258
number: 3
url: https://github.com/kamilsk/semaphore/issues/3
title: "add timeline backoff"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-10-09T09:13:15Z
updatedAt: 2016-10-11T19:23:25Z
lastEditedAt: null
closedAt: 2016-10-11T19:23:25Z
---

# add timeline backoff

``` go
backoff.Liniear(semaphore.Acquire(time.Second), time.Minute)
```

<!-- 2016-10-11T18:44Z https://github.com/kamilsk/semaphore/issues/3#issuecomment-253007390
using retry library instead porting its functionality
-->
