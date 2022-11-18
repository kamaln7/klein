<p align="center">
  <img src="/klein.png" alt="klein logo" width="386" />
</p>

klein is a minimalist URL shortener written in Go. No unnecessary clutter, web UI, features, etc. Just shortening and serving redirections.

## Modularity

klein has three core components that are abstracted into drivers to allow different functionality:

1. auth
   - Handles authentication, guarding access to shortening links
   - Comes with two drivers:
     - Unauthenticated—shorten URLs without authentication
     - Static Key—require a static key/password
     - HTTP Basic—uses HTTP Basic Auth, require a username and password
2. alias
   - Handles generating URL aliases.
   - Comes with two drivers:
     - Alphanumeric—returns a random alphanumeric string with a configurable length
     - Memorable—returns a configurable amount of English words
3. storage
   - Handles storing and reading shortened URLs.
   - Comes with four drivers:
     - File—stores data as text files in a directory
     - Bolt—stores data in a [bolt](https://github.com/boltdb/bolt) database
     - Redis—stores data in a [redis](https://redis.io/) database (ensure you configure save)
     - Spaces.stateful—stores data as a single file in [DigitalOcean Spaces](https://do.co/spaces)
     - Spaces.stateless—stores each URL as an object in [DigitalOcean Spaces](https://do.co/spaces)
     - PostgreSQL—stores data in a [PostgreSQL](https://www.postgresql.org) database
     - Memory—stores data in a temporary map in memory

## Usage

Once installed and configured, there are two actions that you can do:

1. Shorten a URL:
   - Send a POST request to `/` with the following two fields:
     1. `url`—the URL to shorten
     2. `key`—if the Static Key auth driver is enabled
     3. `alias`—a custom alias to be used instead of a randomly-generated one
   - Example cURL command: `curl -X POST -d url=http://github.com/kamaln7/klein -d key=secret_password -d alias=klein_gh http://localhost:5556/`
     - This will create a short URL at `http://localhost:5556/klein_gh` that redirects to `http://github.com/kamaln7/klein`.
2. Look up a URL/serve a redirect:
   - Browse to `http://[path to klein]/[alias]` to access a short URL.

## Installation

✅ Use the docker image `kamaln7/klein`. The `latest` tag is a good bet. See [the releases page](https://github.com/kamaln7/klein/releases) for version numbers.

Or grab the latest binary from [the releases page](https://github.com/kamaln7/klein/releases) and drop it in `/usr/local/bin`, `/opt`, or wherever you like.

### Configuration

klein uses CLI options or environment variables for config. For environment variables, each option is prefixed with `klein` and both dots and dashes are replaced with underscores, eg the environment variable for the `storage.spaces.access-key` option is `KLEIN_STORAGE_SPACES_ACCESS_KEY`.

Running klein without any configuration will use the following default config:

- Aliases are random 5-character alphanumeric strings
- Listens on 127.0.0.1:5556
- No authentication
- Stores URLs as files in a `urls` directory in the current working directory

#### Full list of config options

```
$ klein --help
klein is a minimalist URL shortener.

Usage:
  klein [flags]

Flags:
      --alias.alphanumeric.alpha                           use letters in code (default true)
      --alias.alphanumeric.length int                      alphanumeric code length (default 5)
      --alias.alphanumeric.num                             use numbers in code (default true)
      --alias.driver string                                what alias generation to use (alphanumeric, memorable) (default "alphanumeric")
      --alias.memorable.length int                         memorable word count (default 3)
      --auth.basic.password string                         password for HTTP basic auth
      --auth.basic.username string                         username for HTTP basic auth
      --auth.driver string                                 what auth backend to use (basic, key, none) (default "none")
      --auth.key string                                    upload API key
      --error-template string                              path to error template
  -h, --help                                               help for klein
      --listen string                                      listen address (default "127.0.0.1:5556")
      --root string                                        root redirect
      --storage.boltdb.path string                         path to use for bolt db (default "bolt.db")
      --storage.driver string                              what storage backend to use (file, boltdb, redis, spaces.stateful, sql.pg, memory) (default "file")
      --storage.file.path string                           path to use for file store (default "urls")
      --storage.redis.address string                       address:port of redis instance (default "127.0.0.1:6379")
      --storage.redis.auth string                          password to access redis
      --storage.redis.db int                               db to select within redis
      --storage.spaces.access-key string                   access key for spaces
      --storage.spaces.region string                       region for spaces
      --storage.spaces.secret-key string                   secret key for spaces
      --storage.spaces.space string                        space to use
      --storage.spaces.stateful.path string                path of the file in spaces (default "klein.json")
      --storage.spaces.stateless.cache-duration duration   time to cache spaces results in memory. 0 to disable (default 1m0s)
      --storage.spaces.stateless.path string               path of the directory in spaces to store urls in (default "/klein")
      --storage.sql.pg.database string                     postgresql database (default "klein")
      --storage.sql.pg.host string                         postgresql host (default "localhost")
      --storage.sql.pg.password string                     postgresql password (default "secret")
      --storage.sql.pg.port int32                          postgresql port (default 5432)
      --storage.sql.pg.sslmode string                      postgresql sslmode (default "prefer")
      --storage.sql.pg.table string                        postgresql table (default "klein")
      --storage.sql.pg.user string                         postgresql user (default "klein")
      --url string                                         path to public facing url
```

### Service file

Here's a Systemd service file that you can use with klein:

```
[Unit]
Description=klein
After=network-online.target

[Service]
Restart=on-failure

User=klein
Group=klein

ExecStart=/usr/local/bin/klein

[Install]
WantedBy=multi-user.target
```

Don't forget to add your config to the `ExecStart` line and update `User` and `Group` if necessary. Make sure that klein has permission to write to the URLs directory.

## Development

To manage dependencies, we use [Go modules](https://blog.golang.org/using-go-modules).

To build the app, run `go build`.  
This will produce a binary named `klein`. You can now run the app by running `./klein`

### ❤️ Contributors

- @LukeHandle
- @DMarby

## License

See [./LICENSE](/LICENSE)
