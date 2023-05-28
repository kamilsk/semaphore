---
code: null
id: MDU6SXNzdWUzMTIwMjQ4NTY=
databaseId: 312024856
number: 142
url: https://github.com/kamilsk/semaphore/issues/142
title: "add stress-test use case"
issueType: null
labels:
  - "type:docs"
  - "type:feature"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: OPEN
stateReason: null
duplicateOf: null
createdAt: 2018-04-06T15:44:02Z
updatedAt: 2020-02-15T09:11:18Z
lastEditedAt: null
closedAt: null
---

# add stress-test use case

https://github.com/kamilsk/form-api/issues/47

```
semaphore create 2
semaphore add -- curl -O http://user:pass@localhost:8080/debug/
semaphore add -- hey -cpu=2 localhost:8080/api/v1/{schema}
semaphore wait
```

extend this: provide a profile to https://github.com/google/pprof

<!-- 2019-03-19T20:27Z https://github.com/kamilsk/semaphore/issues/142#issuecomment-474568799
https://github.com/loov/hrtime
-->
