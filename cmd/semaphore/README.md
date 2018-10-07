> # cmd/semaphore [![Tweet][icon_twitter]][twitter_publish]
> [![Analytics][analytics_pixel]][page_promo]
> `semaphore` provides functionality to execute terminal commands in parallel.

[![Awesome][icon_awesome]](https://github.com/avelino/awesome-go#goroutines)
[![Patreon][icon_patreon]](https://www.patreon.com/octolab)
[![Build Status][icon_build]][page_build]
[![Code Coverage][icon_coverage]][page_quality]
[![Code Quality][icon_quality]][page_quality]
[![GoDoc][icon_docs]][page_docs]
[![License][icon_license]](../../LICENSE)

## Concept

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --timeout=1m --notify
```

[![asciicast](https://asciinema.org/a/135943.png)](https://asciinema.org/a/135943)

## Documentation

```
Usage: semaphore COMMAND

Semaphore provides functionality to execute terminal commands in parallel.

Commands:

create	is a command to init a semaphore context
  -debug
    	show error stack trace
  -filename string
    	an absolute path to semaphore context (default "/tmp/semaphore.json")


add	is a command to add a job into a semaphore context
  -debug
    	show error stack trace
  -edit
    	switch to edit mode to read arguments from input (not implemented yet)
  -filename string
    	an absolute path to semaphore context (default "/tmp/semaphore.json")


wait	is a command to execute a semaphore task
  -debug
    	show error stack trace
  -filename string
    	an absolute path to semaphore context (default "/tmp/semaphore.json")
  -notify
    	show notification at the end (not implemented yet)
  -speed int
    	a velocity of report output (characters per second)
  -timeout duration
    	timeout for task execution (default 1m0s)

Version 4.0.0 (commit: ..., build date: ..., go version: go1.9, compiler: gc, platform: darwin/amd64)
```

### Complex example

```bash
$ semaphore create 2
$ semaphore add -- bash -c "cd /tmp; \
    git clone git@github.com:kamilsk/semaphore.git \
    && cd semaphore \
    && echo 'semaphore at revision' \$(git rev-parse HEAD) \
    && rm -rf /tmp/semaphore"
$ semaphore add -- bash -c "cd /tmp; \
    git clone git@github.com:kamilsk/retry.git \
    && cd retry \
    && echo 'retry at revision' \$(git rev-parse HEAD) \
    && rm -rf /tmp/retry"
$ semaphore wait
```

## Installation

### Brew

```bash
$ brew install kamilsk/tap/semaphore
```

### Binary

```bash
$ export REQ_VER=4.0.0  # all available versions are on https://github.com/kamilsk/semaphore/releases
$ export REQ_OS=Linux   # macOS and Windows are also available
$ export REQ_ARCH=64bit # 32bit is also available
$ # wget -q -O semaphore.tar.gz
$ curl -sL -o semaphore.tar.gz \
       https://github.com/kamilsk/semaphore/releases/download/"${REQ_VER}/semaphore_${REQ_VER}_${REQ_OS}-${REQ_ARCH}".tar.gz
$ tar xf semaphore.tar.gz -C "${GOPATH}"/bin/ && rm semaphore.tar.gz
```

### From source code

```bash
$ egg github.com/kamilsk/semaphore@^4.0.0 -- make test install
$ # or use mirror
$ egg bitbucket.org/kamilsk/semaphore@^4.0.0 -- make test install
```

> [egg](https://github.com/kamilsk/egg)<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

### Bash and Zsh completions

You can find completion files [here](https://github.com/kamilsk/shared/tree/dotfiles/bash_completion.d) or
build your own using these commands

```bash
$ semaphore completion bash > /path/to/bash_completion.d/semaphore.sh
$ semaphore completion zsh  > /path/to/zsh-completions/_semaphore.zsh
```

<sup id="egg">1</sup> The project is still in prototyping. [↩](#anchor-egg)

---

[![Gitter][icon_gitter]](https://gitter.im/kamilsk/semaphore)
[![@kamilsk][icon_tw_author]](https://twitter.com/ikamilsk)
[![@octolab][icon_tw_sponsor]](https://twitter.com/octolab_inc)

made with ❤️ by [OctoLab](https://www.octolab.org/)

[analytics_pixel]: https://ga-beacon.appspot.com/UA-109817251-2/semaphore/dev?pixel

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/semaphore.svg?branch=dev
[icon_coverage]:   https://scrutinizer-ci.com/g/kamilsk/semaphore/badges/coverage.png?b=dev
[icon_docs]:       https://godoc.org/github.com/kamilsk/semaphore?status.svg
[icon_gitter]:     https://badges.gitter.im/Join%20Chat.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_quality]:    https://scrutinizer-ci.com/g/kamilsk/semaphore/badges/quality-score.png?b=dev
[icon_research]:   https://img.shields.io/badge/research-in%20progress-yellow.svg
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
[icon_twitter]:    https://img.shields.io/twitter/url/http/shields.io.svg?style=social

[page_build]:      https://travis-ci.org/kamilsk/semaphore
[page_docs]:       https://godoc.org/github.com/kamilsk/semaphore
[page_promo]:      https://github.com/kamilsk/semaphore
[page_quality]:    https://scrutinizer-ci.com/g/kamilsk/semaphore/?branch=dev

[twitter_publish]: https://twitter.com/intent/tweet?text=Semaphore%20pattern%20implementation%20with%20a%20timeout%20of%20lock%2Funlock%20operations%20based%20on%20channels&url=https://github.com/kamilsk/semaphore&via=ikamilsk&hashtags=go,semaphore,throughput,limiter
