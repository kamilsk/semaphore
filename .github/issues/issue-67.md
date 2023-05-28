---
code: null
id: MDU6SXNzdWUyMzgzMTE2Mzk=
databaseId: 238311639
number: 67
url: https://github.com/kamilsk/semaphore/issues/67
title: "cmd/semaphore"
issueType: null
labels:
  - "type:feature"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-06-24T11:16:45Z
updatedAt: 2017-09-01T12:56:17Z
lastEditedAt: 2017-08-13T11:04:09Z
closedAt: 2017-09-01T12:56:17Z
---

# cmd/semaphore

Concept:

```bash
$ semaphore create 4
$ semaphore add -- docker build ...
$ semaphore add -- docker build ...
...
$ semaphore wait
... show progress
```

If `semaphore create ...` didn't call before then `semaphore add ...` create it with defaults.

If `semaphore add ...` didn't call before then `semaphore wait` returns nothing immediately.

<!-- 2017-06-26T13:27Z https://github.com/kamilsk/semaphore/issues/67#issuecomment-311058703
investigate https://www.gnu.org/software/parallel/
-->

<!-- 2017-08-13T11:03Z https://github.com/kamilsk/semaphore/issues/67#issuecomment-322035385
```
Prototype:
$ semaphore create 4
$ semaphore add -- docker build ...
$ semaphore add -- docker build ...
...
$ semaphore wait | semaphore wait --notify --timeout 1h
... show progress (colored output)
[==>........] 2/10

command `docker build ...`
output:
 ...

command...
```
-->
