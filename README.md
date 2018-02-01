# i2c library in Go

[![GoDoc](https://godoc.org/github.com/dasfoo/i2c?status.svg)](http://godoc.org/github.com/dasfoo/i2c)
[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Build Status](https://travis-ci.org/dasfoo/i2c.svg?branch=master)](https://travis-ci.org/dasfoo/i2c)
[![Coverage Status](https://coveralls.io/repos/dasfoo/i2c/badge.svg?branch=master&service=github)](https://coveralls.io/github/dasfoo/i2c?branch=master)

## Installation

The library needs `i2c-dev.h` header file to compile, but the one that is
shipped with `linux-libc-dev` is not sufficient. You have to install the
"full" version of that header file, e.g. on
[Debian](https://packages.debian.org/stretch/all/libi2c-dev/filelist) or
[Ubuntu](https://packages.ubuntu.com/xenial/all/libi2c-dev/filelist) run:

```
# apt install libi2c-dev
```

Then use `go get -v github.com/dasfoo/i2c` normally.
