---
code: null
id: MDU6SXNzdWUyNDk4NjcxMTM=
databaseId: 249867113
number: 94
url: https://github.com/kamilsk/semaphore/issues/94
title: "command deps"
issueType: null
labels:
  - "type:task"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-08-13T11:09:30Z
updatedAt: 2017-09-01T12:56:20Z
lastEditedAt: 2017-08-13T11:09:36Z
closedAt: 2017-09-01T12:56:20Z
---

# command deps

need to knit commands to implement next logic:

If semaphore create ... didn't call before then semaphore add ... create it with defaults.
If semaphore add ... didn't call before then semaphore wait returns nothing immediately.
