---
code: null
id: MDU6SXNzdWUyODMxNDYxMTg=
databaseId: 283146118
number: 131
url: https://github.com/kamilsk/semaphore/issues/131
title: "fix code quality problems"
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
createdAt: 2017-12-19T08:49:08Z
updatedAt: 2019-02-26T06:59:53Z
lastEditedAt: null
closedAt: 2019-02-26T06:59:53Z
---

# fix code quality problems

```
cmd/semaphore/cli.go:158:37:warning: error return value not checked (color.New(color.FgYellow).Fprintln(os.Stderr, "edit component is not ready yet")) (errcheck)
cmd/semaphore/cli.go:166::warning: Expect file permissions to be 0600 or less,MEDIUM,HIGH (gas)
cmd/semaphore/cli.go:172::warning: declaration of "err" shadows declaration at cli.go:166 (vetshadow)
cmd/semaphore/cli.go:285::warning: cyclomatic complexity 11 of function (*WaitCommand).Do() is high (> 10) (gocyclo)
cmd/semaphore/cli.go:286::warning: Expect file permissions to be 0600 or less,MEDIUM,HIGH (gas)
cmd/semaphore/cli.go:292::warning: declaration of "err" shadows declaration at cli.go:286 (vetshadow)
cmd/semaphore/cli.go:363:37:warning: error return value not checked (color.New(color.FgYellow).Fprintln(os.Stderr, "notify component is not ready yet")) (errcheck)
cmd/semaphore/cli.go:422:33:warning: error return value not checked (color.New(color.FgRed).Fprintf(c.Output, format, c.Error)) (errcheck)
cmd/semaphore/main.go:51:10:warning: error return value not checked (help.Do()) (errcheck)
cmd/semaphore/task.go:100::warning: Errors unhandled.,LOW,HIGH (gas)
cmd/semaphore/task.go:106::warning: Subprocess launching with variable.,HIGH,HIGH (gas)
cmd/semaphore/task.go:93::warning: Errors unhandled.,LOW,HIGH (gas)
cmd/semaphore/task.go:98:17:warning: error return value not checked (io.WriteString(s, j.String())) (errcheck)
cmd/semaphore/task.go:98::warning: Errors unhandled.,LOW,HIGH (gas)
```

<!-- 2019-02-26T06:59Z https://github.com/kamilsk/semaphore/issues/131#issuecomment-467320080
will be fixed in https://github.com/kamilsk/semaphore.cli
-->
