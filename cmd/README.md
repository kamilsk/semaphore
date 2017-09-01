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
