---
code: null
id: MDU6SXNzdWUyODIzMzY4MjU=
databaseId: 282336825
number: 130
url: https://github.com/kamilsk/semaphore/issues/130
title: "performance review"
issueType: null
labels:
  - "status:research"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-12-15T06:59:28Z
updatedAt: 2019-02-19T08:05:54Z
lastEditedAt: null
closedAt: 2019-02-19T08:05:53Z
---

# performance review

need to define overhead of the package

<!-- 2017-12-29T21:32Z https://github.com/kamilsk/semaphore/issues/130#issuecomment-354502877
`time.NewTimer` instead of `time.After`: investigate "under the hood" benefits (related to GC)
-->

<!-- 2017-12-29T21:34Z https://github.com/kamilsk/semaphore/issues/130#issuecomment-354503083
`time.After` do not release app resource, but `timer` can do it
-->

<!-- 2017-12-29T21:38Z https://github.com/kamilsk/semaphore/issues/130#issuecomment-354503376
```
timer := time.NewTimer(deadline)
semaphore.Acquire(timer.C)
if ok -> timer.Stop() and continue
```
-->

<!-- 2019-02-19T08:05Z https://github.com/kamilsk/semaphore/issues/130#issuecomment-465027860
won't do
-->
