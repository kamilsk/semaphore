---
code: null
id: MDU6SXNzdWUxODIwMDU5MTI=
databaseId: 182005912
number: 5
url: https://github.com/kamilsk/semaphore/issues/5
title: "simplify Semaphore interface"
issueType: null
labels:
  - "status:research"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-10-10T12:46:16Z
updatedAt: 2016-10-11T08:55:23Z
lastEditedAt: null
closedAt: 2016-10-11T08:55:23Z
---

# simplify Semaphore interface

`Capacity() int, Occupied() int` are required ? Maybe change `Acquire(time.Duration) error` signature to `Acquire(time.Duration) (int, error)` where `int` is current occupied value
