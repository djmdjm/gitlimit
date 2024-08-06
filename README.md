Overview
========

This is a simple program to limit access to git repositories for SSH git clients run from authorized_keys. E.g.

```
restrict,command="/usr/local/bin/gitlimit /git/scripts.git" ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICMvIKQRTt6jgfzrhStlR0kNidaKRBkT4deA21P/zAao djm@semitrusted
```

Installation
============

You can install the tool for the current user with `go install github.com/djmdjm/gitlimit@latest`.

