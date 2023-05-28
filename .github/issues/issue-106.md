---
code: null
id: MDU6SXNzdWUyNTYxODE0MDI=
databaseId: 256181402
number: 106
url: https://github.com/kamilsk/semaphore/issues/106
title: "rate limiter example"
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
createdAt: 2017-09-08T08:40:55Z
updatedAt: 2017-10-06T20:50:44Z
lastEditedAt: 2017-09-08T08:41:33Z
closedAt: 2017-10-06T20:50:44Z
---

# rate limiter example

The concept:
```go
sem := limiter.ForUser(request.Context())
release, err := sem.Acquire(semaphore.WithTimeout(config.SLA))
if err != nil {
	http.Error(rw, "operation timeout", http.StatusGatewayTimeout)
	return
}
deadline := semaphore.WithDeadline(time.Now().Add(config.LimitRange))
go func() {
	release()
	release, err = sem.Acquire(deadline)
	if err != nil { return }
	<-deadline
	release()
}()
// handle user request
```
