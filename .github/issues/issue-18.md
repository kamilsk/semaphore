---
code: null
id: MDU6SXNzdWUxODYwNzYyMzg=
databaseId: 186076238
number: 18
url: https://github.com/kamilsk/semaphore/issues/18
title: "mistake with monitored semaphore example"
issueType: null
labels:
  - "type:bug"
  - "type:docs"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-10-29T12:35:12Z
updatedAt: 2016-10-29T14:23:02Z
lastEditedAt: null
closedAt: 2016-10-29T14:23:02Z
---

# mistake with monitored semaphore example

`func (sem *monitoredSemaphore) Acquire(timeout time.Duration) error {` without mutex protection
