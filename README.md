<br>
<p align="center">
    <img src="./docs/medias/logo.png" alt="Driplimit Icon">
    <h3 align="center">Fast, secure and easy api key management</h3>
</p>

<br>

> **Driplimit is currently in heavy development. API and schemas will break, therefore, you should not run it in production until it reaches v1.0.0**

## What is Driplimit ?

Driplimit is a key management service provided as an API. It keeps your user's keys securely and enables your application to control the validity, the expiration and the rate usage every time they send a token. 

## How to use Driplimit ?

Driplimit comes as a single binary that exposes a web server. Every tasks can be done by calling HTTP RPC endpoints.

See the [How to use section](./docs/how-to-use.md) in order to create your first keys.

### SDK

* [Go Client (official)](https://github.com/i4n-co/driplimit/tree/main/pkg/client)

### RPC API Documentation

Available soon.

## How to install Driplimit ?

The best way to install `driplimit` is to use the official Docker image (available soon).

Alternatively, you can choose to build the binary with the standard Go (1.22+) toolchain:

```bash
go install github.com/i4n-co/driplimit/cmd/driplimit
```

## How to configure Driplimit ?

`driplimit` is configurable via `Env vars`, the complete list of available configurations can be generated via the command line:

```bash
$ driplimit -print-defaults
ADDR=127.0.0.1
CACHE_DURATION=30s
DATABASE_NAME=driplimit.db
DATA_DIR=
GZIP_COMPRESSION=false
KEYS_CACHE_SIZE=65536
LOG_FORMAT=text
LOG_SEVERITY=info
MODE=authoritative
PORT=7131
ROOT_KEYS_CACHE_SIZE=2048
TIMEOUT=5s
UPSTREAM_URL=
```

Driplimit can also be configured via env-file like so:

`$ driplimit -env-file=/etc/driplimit/config.env`

## What are the driplimit modes ?

Driplimit can run with 3 modes:

* `authoritative` - the source of truth, local sqlite database, synchronous (slower, yet very fast)
* `proxy` - proxy requests to `UPSTREAM_URL`, cache some responses in-memory, predicts checks, asynchronous (faster)
* `async_authoritative` - has a local database but maintains an in-memory cache for faster asynchronous responses (when authoritative is not fast enough)

You can have a central authoritative driplimit server while maintaining proxies close to your apps, allowing fast, distributed, key management and rate limiting. Think of DNS infrastructure, but for keys.

Other scalability features will be added soon.


## Open-source, not open-contribution

[Similar to SQLite](https://www.sqlite.org/copyright.html), Driplimit is open
source but closed to contributions, yet. This keeps the code base free of proprietary
or licensed code but it also helps us continue to maintain and build Driplimit.

We are grateful for community involvement, bug reports, & feature requests. We do
not wish to come off as anything but welcoming, however, for now, we've
made the decision to keep this project closed to contributions for long term viability of the project.

## About

This software is distributed under [GNU AFFERO GENERAL PUBLIC LICENSE](./LICENCE.md)
