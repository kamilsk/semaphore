---
code: null
id: MDU6SXNzdWUxOTc1MTA0NjY=
databaseId: 197510466
number: 44
url: https://github.com/kamilsk/semaphore/issues/44
title: "race condition"
issueType: null
labels:
  - "type:critical"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2016-12-25T18:10:50Z
updatedAt: 2026-11-08T14:08:05Z
lastEditedAt: 2026-11-08T14:08:05Z
closedAt: 2016-12-25T18:49:30Z
---

# race condition

```
docker run --rm \
                   -v "/Users/.../Development/go/src/github.com/kamilsk/semaphore":"/go/src/github.com/kamilsk/semaphore" \
                   -w "/go/src/github.com/kamilsk/semaphore" \
                   golang:1.7 \
                   go get -d -t `go list ./... | grep -v /vendor/` && \
                   go test `go list ./... | grep -v /vendor/` -race -v

...

=== RUN   TestBinarySemaphore
==================
WARNING: DATA RACE
Write at 0x00c4200664e0 by goroutine 14:
  runtime.closechan()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/runtime/chan.go:303 +0x0
  github.com/kamilsk/semaphore/classic.semaphore.Flush()
      /Users/.../Development/go/src/github.com/kamilsk/semaphore/classic/classic_test.go:11 +0x38
  github.com/kamilsk/semaphore/classic.TestBinarySemaphore()
      /Users/.../Development/go/src/github.com/kamilsk/semaphore/classic/classic_test.go:89 +0x425
  testing.tRunner()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/testing/testing.go:610 +0xc9

Previous read at 0x00c4200664e0 by goroutine 17:
  runtime.chansend()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/runtime/chan.go:122 +0x0
  github.com/kamilsk/semaphore/classic.semaphore.P()
      /Users/.../Development/go/src/github.com/kamilsk/semaphore/classic/classic.go:25 +0x5f
  github.com/kamilsk/semaphore/classic.semaphore.Lock()
      /Users/.../Development/go/src/github.com/kamilsk/semaphore/classic/classic.go:74 +0x41
  github.com/kamilsk/semaphore/classic.TestBinarySemaphore.func1.1()
      /Users/.../Development/go/src/github.com/kamilsk/semaphore/classic/classic_test.go:73 +0x70

Goroutine 14 (running) created at:
  testing.(*T).Run()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/testing/testing.go:646 +0x52f
  testing.RunTests.func1()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/testing/testing.go:793 +0xb9
  testing.tRunner()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/testing/testing.go:610 +0xc9
  testing.RunTests()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/testing/testing.go:799 +0x4ba
  testing.(*M).Run()
      /usr/local/Cellar/go/1.7.4_1/libexec/src/testing/testing.go:743 +0x12f
  main.main()
      github.com/kamilsk/semaphore/classic/_test/_testmain.go:66 +0x1b8

Goroutine 17 (running) created at:
  github.com/kamilsk/semaphore/classic.TestBinarySemaphore.func1()
      /Users/.../Development/go/src/github.com/kamilsk/semaphore/classic/classic_test.go:76 +0x168
==================
--- FAIL: TestBinarySemaphore (0.00s)
        classic_test.go:86: [first second] is not equal to [first second third]

...

make: *** [docker-test-1.7] Error 1
```
