Workstream
==========

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Flexible work time tracking API.

Requirements
------------

* [KUSANAGI framework](http://kusanagi.io) 2.0+
* [Go](https://golang.org/dl/) 1.14+

Installation
------------

Copy the *workstream.ini.dist* file to *workstream.ini*.

Install the support to apply database migration:

```bash
 go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
```

Run the migrations:

```bash
 make migrate
```

Install the binaries using the following command:

```bash
 make all
```

Run
---

Run the realm using:

```bash
 kusanagi realm start -c workstream.yaml -V workstream.ini -e default
```
