> # cmd/semaphore
>
> `semaphore` provides functionality to execute terminal commands in parallel.

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#goroutines)
[![Patreon](https://img.shields.io/badge/patreon-donate-orange.svg)](https://www.patreon.com/octolab)
[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](../../LICENSE)

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

Version 4.2.1 (commit: c3021037717c136851e639a0805900c062c73ce0, build date: 2017-10-29T07:30:15Z, go version: go1.9, compiler: gc, platform: darwin/amd64)
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
$ export SEM_V=4.2.1    # all available versions are on https://github.com/kamilsk/semaphore/releases
$ export REQ_OS=Linux   # macOS and Windows are also available
$ export REQ_ARCH=64bit # 32bit is also available
$ wget -q -O semaphore.tar.gz \
      https://github.com/kamilsk/semaphore/releases/download/${SEM_V}/semaphore_${SEM_V}_${REQ_OS}-${REQ_ARCH}.tar.gz
$ tar xf semaphore.tar.gz -C "${GOPATH}"/bin/
$ rm semaphore.tar.gz
```

### From source code

```bash
$ go get -d github.com/kamilsk/semaphore
$ cd "${GOPATH}"/src/github.com/kamilsk/semaphore
$ make cmd-deps-local # or cmd-deps, if you don't have the dep binary but have the docker
$ make cmd-install
```

## Command-line completion

### Useful articles

- [Command-line completion | Docker Documentation](https://docs.docker.com/compose/completion/)

### Bash

Make sure [bash completion](https://github.com/scop/bash-completion) is installed.

- On a current Linux (in a non-minimal installation), bash completion should be available.
- On a macOS, install by the command `brew install bash-completion`.

Place the completion script in `/etc/bash_completion.d/` (or `$(brew --prefix)/etc/bash_completion.d/` on a macOS):

```bash
$ sudo curl -L https://raw.githubusercontent.com/kamilsk/shared/dotfiles/bash_completion.d/semaphore.sh \
            -o /etc/bash_completion.d/semaphore
```

On a macOS, add the following to your `~/.bash_profile`:

```bash
if [ -f $(brew --prefix)/etc/bash_completion ]; then
    source $(brew --prefix)/etc/bash_completion
fi
```

If you're using MacPorts instead of brew you'll need to slightly modify your steps to the following:

- Run `sudo port install bash-completion` to install bash completion.
- Add the following lines to `~/.bash_profile`:
```bash
if [ -f /opt/local/etc/profile.d/bash_completion.sh ]; then
    source /opt/local/etc/profile.d/bash_completion.sh
fi
```

You can source your `~/.bash_profile` or launch a new terminal to utilize completion.

### Zsh

Place the completion script in your `/path/to/zsh/completion`, using, e.g., `~/.zsh/completion/`:

```bash
$ mkdir -p ~/.zsh/completion
$ curl -L https://raw.githubusercontent.com/kamilsk/shared/dotfiles/bash_completion.d/semaphore.zsh \
       -o ~/.zsh/completion/_semaphore
```

Include the directory in your `$fpath`, e.g., by adding in `~/.zshrc`:

```bash
fpath=(~/.zsh/completion $fpath)
```

Make sure `compinit` is loaded or do it by adding in `~/.zshrc`:

```bash
autoload -Uz compinit && compinit -i
```

Then reload your shell:

```bash
exec $SHELL -l
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/semaphore)
[![@kamilsk](https://img.shields.io/badge/author-%40kamilsk-blue.svg)](https://twitter.com/ikamilsk)
[![@octolab](https://img.shields.io/badge/sponsor-%40octolab-blue.svg)](https://twitter.com/octolab_inc)

## Notes

- made with ❤️ by [OctoLab](https://www.octolab.org/)

[![Analytics](https://ga-beacon.appspot.com/UA-109817251-2/semaphore/cmd)](https://github.com/igrigorik/ga-beacon)
