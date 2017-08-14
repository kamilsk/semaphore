> # semaphore/cmd
>
> Package cmd contains console tools.

## semaphore

> `semaphore` provides functionality to execute terminal commands in parallel.

### Concept

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --notify --timeout=1m
```
