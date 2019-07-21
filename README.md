> # 🚦 semaphore
>
> Semaphore pattern implementation with timeout of lock/unlock operations.

[![Build][icon_build]][page_build]
[![Quality][icon_quality]][page_quality]
[![Documentation][icon_docs]][page_docs]
[![Coverage][icon_coverage]][page_coverage]
[![Awesome][icon_awesome]][page_awesome]

## 💡 Idea

The semaphore provides API to control access to a shared resource by multiple goroutines or limit throughput.

```go
releaser, err := semaphore.Acquire(breaker.BreakByTimeout(time.Second))
if err != nil {
	// timeout exceeded
}
defer releaser.Release()
```

Full description of the idea is available [here][design].

## 🏆 Motivation

...

## 🤼‍♂️ How to

### Quick start

```go
limiter := semaphore.New(1000)

http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
	if _, err := limiter.Acquire(
		breaker.BreakByContext(
			context.WithTimeout(req.Context(), time.Second),
		),
	); err != nil {
		http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	defer limiter.Release()

	// handle request
})

log.Fatal(http.ListenAndServe(":80", http.DefaultServeMux))
```

## 🧩 Integration

The library uses [SemVer](https://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe through major releases.
You can use [go modules](https://github.com/golang/go/wiki/Modules) or
[dep](https://golang.github.io/dep/) to manage its version.

The **[master][legacy]** is a feature frozen branch for versions **4.3.x** and no longer maintained.

```bash
$ dep ensure -add github.com/kamilsk/semaphore@4.3.1
```

The **[v4][]** branch is a continuation of the **[master][legacy]** branch for versions **v4.4.x**
to better integration with [go modules](https://github.com/golang/go/wiki/Modules).

```bash
$ go get -u github.com/kamilsk/semaphore/v4@v4.3.1
```

The **[v5][]** branch is an actual development branch.

```bash
$ go get -u github.com/kamilsk/semaphore    # inside GOPATH and for old Go versions

$ go get -u github.com/kamilsk/semaphore/v5 # inside Go module, works well since Go 1.11

$ dep ensure -add github.com/kamilsk/semaphore@v5.0.0-rc1
```

Version **v5** focused on integration with the 🚧 [breaker][] package.

## 🤲 Outcomes

### Console tool for command execution in parallel

This example shows how to execute many console commands in parallel.

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --timeout=1m --notify
```

[![asciicast][cli.preview]][cli.demo]

See more details [here][cli].

---

made with ❤️ for everyone

[icon_awesome]:     https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:       https://travis-ci.org/kamilsk/semaphore.svg?branch=v5
[icon_coverage]:    https://api.codeclimate.com/v1/badges/0261f2170c785702034f/test_coverage
[icon_docs]:        https://godoc.org/github.com/kamilsk/semaphore?status.svg
[icon_quality]:     https://goreportcard.com/badge/github.com/kamilsk/semaphore

[page_awesome]:     https://github.com/avelino/awesome-go#goroutines
[page_build]:       https://travis-ci.org/kamilsk/semaphore
[page_coverage]:    https://codeclimate.com/github/kamilsk/semaphore/test_coverage
[page_docs]:        https://godoc.org/github.com/kamilsk/semaphore
[page_quality]:     https://goreportcard.com/report/github.com/kamilsk/semaphore

[breaker]:          https://github.com/kamilsk/breaker
[cli]:              https://github.com/kamilsk/semaphore.cli
[cli.demo]:         https://asciinema.org/a/136111
[cli.preview]:      https://asciinema.org/a/136111.png
[design]:           https://www.notion.so/octolab/semaphore-7d5ebf715d0141d1a8fa045c7966be3b?r=0b753cbf767346f5a6fd51194829a2f3
[egg]:              https://github.com/kamilsk/egg
[promo]:            https://github.com/kamilsk/semaphore

[legacy]:           https://github.com/kamilsk/semaphore/tree/master
[v4]:               https://github.com/kamilsk/semaphore/tree/v4
[v5]:               https://github.com/kamilsk/semaphore/projects/6

[tmp.docs]:         https://nicedoc.io/kamilsk/semaphore?theme=dark
[tmp.history]:      https://github.githistory.xyz/kamilsk/semaphore/blob/v5/README.md
