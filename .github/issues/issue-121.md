---
code: null
id: MDU6SXNzdWUyNzgwMDkwNjc=
databaseId: 278009067
number: 121
url: https://github.com/kamilsk/semaphore/issues/121
title: "remove FlagSet"
issueType: null
labels:
  - "type:refactoring"
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-4]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2017-11-30T06:11:34Z
updatedAt: 2019-02-26T06:56:19Z
lastEditedAt: null
closedAt: 2019-02-26T06:56:19Z
---

# remove FlagSet

- remove the method from interface https://github.com/kamilsk/semaphore/blob/1e58f96d738b1c6c0a130c7764cac3c67b70ea89/cmd/semaphore/cli.go#L28-L29
- extend the method to accept arguments https://github.com/kamilsk/semaphore/blob/1e58f96d738b1c6c0a130c7764cac3c67b70ea89/cmd/semaphore/cli.go#L35
- extend the method to return parsed arguments https://github.com/kamilsk/semaphore/blob/1e58f96d738b1c6c0a130c7764cac3c67b70ea89/cmd/semaphore/cli.go#L42
- these parts https://github.com/kamilsk/semaphore/blob/1e58f96d738b1c6c0a130c7764cac3c67b70ea89/cmd/semaphore/cli.go#L80 and https://github.com/kamilsk/semaphore/blob/1e58f96d738b1c6c0a130c7764cac3c67b70ea89/cmd/semaphore/cli.go#L105 should be removed
- these parts
```go
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
		c.Flags.BoolVar(&c.Edit, "edit", false, "switch to edit mode to read arguments from input (not implemented yet)")
	}
	return c.Flags
```
should be refactored to stateless

<!-- 2019-02-26T06:56Z https://github.com/kamilsk/semaphore/issues/121#issuecomment-467319202
moved to https://github.com/kamilsk/semaphore.cli/issues/2
-->
