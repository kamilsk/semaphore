> # semaphore/cmd
>
> Package cmd contains CLI tools.

## semaphore

> `semaphore` provides functionality to execute terminal commands in parallel.
> > status: **experimental**

### Concept

```bash
$ semaphore create 4
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --timeout=1m --notify
```
