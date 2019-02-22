> # 🚦 semaphore
>
> Semaphore pattern implementation with timeout of lock/unlock operations.

[![Awesome][icon_awesome]][awesome]
[![Patreon][icon_patreon]][support]
[![GoDoc][icon_docs]][docs]
[![Research][icon_research]][research]
[![License][icon_license]][license]

## Important news

The **[master][legacy]** is a feature frozen branch for versions **4.3.x** and no longer maintained.

```bash
$ dep ensure -add github.com/kamilsk/semaphore@4.3.1
```

The **[v4][]** branch is a continuation of the **[master][legacy]** branch for versions **v4.4.x**
to better integration with [Go Modules][gomod].

```bash
$ go get -u github.com/kamilsk/semaphore/v4@v4.3.1
```

The **[v5][]** branch is an actual development branch.

```bash
$ go get -u github.com/kamilsk/semaphore/v5

$ dep ensure -add github.com/kamilsk/semaphore@v5.0.0-rc1
```

Version **v5.x.y** focused on integration with the 🚧 [breaker][] and the 🧰 [platform][] packages.

## Usage

### Quick start

```go
limiter := semaphore.New(1000)

http.HandleFunc("/", func(rw http.ResponseWriter, _ *http.Request) {
	if _, err := limiter.Acquire(semaphore.WithTimeout(time.Minute)); err != nil {
		http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	defer limiter.Release()
	// handle request
})

log.Fatal(http.ListenAndServe(":80", nil))
```

### Console tool for command execution in parallel

This example shows how to execute many console commands in parallel.

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --timeout=1m --notify
```

[![asciicast](https://asciinema.org/a/136111.png)](https://asciinema.org/a/136111)

See more details [here](cmd/semaphore).

## Installation

```bash
$ go get github.com/kamilsk/semaphore
$ # or use mirror
$ egg bitbucket.org/kamilsk/semaphore
```

> [egg][]<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

## Update

This library is using [SemVer](https://semver.org/) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.

<sup id="egg">1</sup> The project is still in prototyping. [↩](#anchor-egg)

---

[![Gitter][icon_gitter]][gitter]
[![@kamilsk][icon_tw_author]][author]
[![@octolab][icon_tw_sponsor]][sponsor]

made with ❤️ by [OctoLab][octolab]

[awesome]:         https://github.com/avelino/awesome-go#goroutines
[build]:           https://travis-ci.org/kamilsk/semaphore
[docs]:            https://godoc.org/github.com/kamilsk/semaphore
[gitter]:          https://gitter.im/kamilsk/semaphore
[license]:         LICENSE
[promo]:           https://github.com/kamilsk/semaphore
[quality]:         https://scrutinizer-ci.com/g/kamilsk/semaphore/?branch=v5
[research]:        https://github.com/kamilsk/go-research/tree/master/projects/semaphore
[legacy]:          https://github.com/kamilsk/semaphore/tree/master
[v4]:              https://github.com/kamilsk/semaphore/tree/v4
[v5]:              https://github.com/kamilsk/semaphore/projects/6

[egg]:             https://github.com/kamilsk/egg
[breaker]:         https://github.com/kamilsk/breaker
[gomod]:           https://github.com/golang/go/wiki/Modules
[platform]:        https://github.com/kamilsk/platform

[author]:          https://twitter.com/ikamilsk
[octolab]:         https://www.octolab.org/
[sponsor]:         https://twitter.com/octolab_inc
[support]:         https://www.patreon.com/octolab

[analytics]:       https://ga-beacon.appspot.com/UA-109817251-2/semaphore/v5?pixel
[tweet]:           https://twitter.com/intent/tweet?text=Semaphore%20pattern%20implementation%20with%20a%20timeout%20of%20lock%2Funlock%20operations&url=https://github.com/kamilsk/semaphore&via=ikamilsk&hashtags=go,semaphore,throughput,limiter

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/semaphore.svg?branch=v5
[icon_coverage]:   https://scrutinizer-ci.com/g/kamilsk/semaphore/badges/coverage.png?b=v5
[icon_docs]:       https://godoc.org/github.com/kamilsk/semaphore?status.svg
[icon_gitter]:     https://badges.gitter.im/Join%20Chat.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_quality]:    https://scrutinizer-ci.com/g/kamilsk/semaphore/badges/quality-score.png?b=v5
[icon_research]:   https://img.shields.io/badge/research-in%20progress-yellow.svg
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
[icon_twitter]:    https://img.shields.io/twitter/url/http/shields.io.svg?style=social
