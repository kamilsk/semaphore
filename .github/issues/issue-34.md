---
code: null
id: MDU6SXNzdWUxOTM1MzAxOTM=
databaseId: 193530193
number: 34
url: https://github.com/kamilsk/semaphore/issues/34
title: "Special case for simple semaphore with zero timeout"
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
createdAt: 2016-12-05T15:36:24Z
updatedAt: 2016-12-06T13:05:19Z
lastEditedAt: null
closedAt: 2016-12-06T13:05:19Z
---

# Special case for simple semaphore with zero timeout

https://github.com/kamilsk/semaphore/blob/master/simple/simple.go#L33
When zero timeout passed semaphore should wait infinitely.
Any objections?
