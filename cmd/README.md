> # semaphore/cmd
>
> Package cmd contains console tools.

## semaphore

> `semaphore` provides functionality to execute terminal commands in parallel.

[![asciicast](https://asciinema.org/a/135943.png)](https://asciinema.org/a/135943)

### Concept

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --timeout=1m --notify
```

### Documentation

```
Usage: semaphore COMMAND

Semaphore provides functionality to execute terminal commands in parallel.

Commands:

create	is a command to init a semaphore context
  -filename string
    	an absolute path to semaphore context (default "/tmp/semaphore.json")


add	is a command to add a job into a semaphore context
  -edit
    	switch to edit mode to read arguments from input (not implemented yet)
  -filename string
    	an absolute path to semaphore context (default "/tmp/semaphore.json")


wait	is a command to execute a semaphore task
  -filename string
    	an absolute path to semaphore context (default "/tmp/semaphore.json")
  -notify
    	show notification at the end (not implemented yet)
  -timeout duration
    	timeout for task execution (default 1m0s)
```

### Installation

#### Brew

```bash
$ brew install kamilsk/tap/semaphore
```

#### Binary

```bash
$ export SEM_V=4.0.0    # all available versions are https://github.com/kamilsk/semaphore/releases
$ export REQ_OS=Linux   # macOS and Windows also available
$ export REQ_ARCH=64bit # 32bit also available
$ wget -q -O semaphore.tar.gz \
      https://github.com/kamilsk/semaphore/releases/download/${SEM_V}/semaphore_${SEM_V}_${REQ_OS}-${REQ_ARCH}.tar.gz
$ tar xf semaphore.tar.gz -C "${GOPATH}"/bin/
$ rm semaphore.tar.gz
```

#### From source code

```bash
$ go get -d github.com/kamilsk/semaphore
$ cd "${GOPATH}"/src/github.com/kamilsk/semaphore
$ make cmd-deps-local # or cmd-deps if you don't have glide binary but have docker app
$ make cmd-install
```

### Command-line completion

#### Bash

Make sure [bash completion](https://github.com/scop/bash-completion) is installed.

- On a current Linux OS (in a non-minimal installation), bash completion should be available.
- On a Mac, install with `brew install bash-completion`.

Place the completion script in `/etc/bash_completion.d/` (or `$(brew --prefix)/etc/bash_completion.d/` on a Mac):

```bash
$ sudo curl -L https://raw.githubusercontent.com/kamilsk/shared/dotfiles/bash_completion.d/semaphore.sh \
            -o /etc/bash_completion.d/semaphore
```

On a Mac, add the following to your `~/.bash_profile`:

```bash
if [ -f $(brew --prefix)/etc/bash_completion ]; then
    source $(brew --prefix)/etc/bash_completion
fi
```

If you’re using MacPorts instead of brew you'll need to slightly modify your steps to the following:

Run `sudo port install bash-completion` to install bash completion. Add the following lines to `~/.bash_profile`:

```bash
if [ -f /opt/local/etc/profile.d/bash_completion.sh ]; then
    source /opt/local/etc/profile.d/bash_completion.sh
fi
```

You can source your `~/.bash_profile` or launch a new terminal to utilize completion.

#### Zsh

Place the completion script in your `/path/to/zsh/completion`, using e.g. `~/.zsh/completion/`:

```bash
$ mkdir -p ~/.zsh/completion
$ curl -L https://raw.githubusercontent.com/kamilsk/shared/dotfiles/bash_completion.d/semaphore.zsh \
       -o ~/.zsh/completion/_semaphore
```

Include the directory in your `$fpath`, e.g. by adding in `~/.zshrc`:

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
