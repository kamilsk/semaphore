> # semaphore/cmd
>
> Package cmd contains CLI tools.

## semaphore

> `semaphore` provides functionality to execute terminal commands in parallel.
> > status: **experimental**

### Concept

```bash
$ semaphore create --capacity=4 --timeout=1m
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait
```
