---
code: null
id: MDU6SXNzdWUyNDUxMjM1MTc=
databaseId: 245123517
number: 85
url: https://github.com/kamilsk/semaphore/issues/85
title: "handle os.Interrupt signal"
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
createdAt: 2017-07-24T15:45:17Z
updatedAt: 2017-08-10T17:16:40Z
lastEditedAt: null
closedAt: 2017-08-10T17:16:40Z
---

# handle os.Interrupt signal

```
c := make(chan os.Signal)
signal.Notify(c, os.Interrupt)
```
