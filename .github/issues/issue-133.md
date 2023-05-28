---
code: null
id: MDU6SXNzdWUyODg1NzQ5OTY=
databaseId: 288574996
number: 133
url: https://github.com/kamilsk/semaphore/issues/133
title: "goroutine leak"
issueType: null
labels:
  - "type:critical"
  - "status:research"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: OPEN
stateReason: null
duplicateOf: null
createdAt: 2018-01-15T12:15:43Z
updatedAt: 2018-05-04T09:05:18Z
lastEditedAt: 2018-01-15T12:27:36Z
closedAt: null
---

# goroutine leak

- [ ] https://github.com/kamilsk/semaphore/blob/master/channel.go#L42-L46
- [ ] https://github.com/kamilsk/semaphore/blob/master/context.go#L11-L12
- [ ] https://github.com/kamilsk/semaphore/blob/master/channel.go#L52-L63 - channel + cancel, because `time.After` is bad
- [ ] https://github.com/kamilsk/semaphore/blob/master/channel.go#L35-L49 - the same as above, because we don't need to listen signal when work is done
- [ ] investigate more

<!-- 2018-05-04T09:05Z https://github.com/kamilsk/semaphore/issues/133#issuecomment-386544156
https://go101.org/article/concurrent-common-mistakes.html
-->
