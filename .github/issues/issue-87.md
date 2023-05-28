---
code: null
id: MDU6SXNzdWUyNDY1OTQ4MDY=
databaseId: 246594806
number: 87
url: https://github.com/kamilsk/semaphore/issues/87
title: "add example with os.Interrupt"
issueType: null
labels:
  - "type:docs"
assignees:
  - kamilsk
milestone: "[[milestone-3]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-07-30T16:44:15Z
updatedAt: 2017-08-29T15:06:16Z
lastEditedAt: null
closedAt: 2017-08-29T15:06:16Z
---

# add example with os.Interrupt

concept
```go
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
```
