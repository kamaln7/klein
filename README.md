<p align="center">
  <img src="/klein.png" alt="klein logo" width="386" />
</p>

klein is a minimalist URL shortener written in Go. No unnecessary clutter, web UI, features, etc. Just shortening and serving redirections.

## Modularity

klein has three core components that are abstracted into "modules" to allow different functionality:

1. auth
   * Handles authentication, gaurding access to shortening links
   * Comes with two modules:
     * Unauthenticated—shorten URLs without authentication
     * Static Key—require a static key/password
2. alias
   * Handles generating URL aliases.
   * Comes with two modules:
     * Alphanumeric—returns a random alphanumeric string with a configurable length
     * Memorable—returns a configurable amount of English words
3. storage
   * Handles storing and reading shortened URLs.
   * Comes with three modules:
     * File—stores URL data as text files in a directory
     * Bolt—stores URL data in a [bolt](https://github.com/boltdb/bolt) database
     * Redis—stores URL data in a [redis](https://redis.io/) database (ensure you configure save)

## Installation

Grab the latest binary from [the releases page](https://github.com/kamaln7/klein/releases) and drop it in `/usr/local/bin`, `/opt`, or wherever you like.

### Configuration

klein uses CLI options for config.

| option                     | description                                                                                                  | default              |
| -------------------------- | ------------------------------------------------------------------------------------------------------------ | -------------------- |
| `-alphanumeric.length int` | Alias length for the Alphanumeric alias module.                                                              |                      |
| `-alphanumeric.alpha bool` | Include English alphabet characters in the Alphanumeric module aliases.                                      | `true`               |
| `-alphanumeric.num bool`   | Include numbers in the Alphanumeric module aliases.                                                          | `true`               |
| `-auth.key string`         | `key` for the Static Key auth module.                                                                        |                      |
| `-auth.username string`    | Username for the HTTP Basic Auth auth module.                                                                |                      |
| `-auth.password string`    | Password for the HTTP Basic Auth auth module.                                                                |                      |
| `-memorable.length int`    | Alias length for the Memorable alias module.                                                                 |                      |
| `-listenAddr string`       | The network address to listen on.                                                                            | `127.0.0.1:5556`     |
| `-file.path string`        | Path to the storage directory for the File storage module.                                                   |                      |
| `-bolt.path string`        | Path to the bolt database for the Bolt storage module.                                                       |                      |
| `-redis.address string`    | Address:Port for the Redis storage module.                                                                   |                      |
| `-redis.auth string`       | Authentication string for the Redis storage module.                                                          |                      |
| `-redis.db int`            | Database ID for the Redis storage module.                                                                    | `0`                  |
| `-root string`             | The URL to redirect to when the `/` path is accessed. Returns a `404 Not Found` error if left blank.         |                      |
| `-template string`         | Path to 404 document to serve in case a 404 error occurs. Returns a plaintext "404 not found" if left blank. |                      |
| `-url string`              | Base URL to the hosted instance of the klein.                                                                | `http://listenAddr/` |

You must specify one storage provider (`file.path`/`bolt.path`/`redis.address`) and one alias provider (`alphanumeric.length`/`memorable.length`).

If none of `auth.key`, `auth.username`, and `auth.password` are provided, the server is run without authentication.

### Service file

Here's a SystemD service file that you can use with klein:

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

## Usage

Once installed and configured, there are two actions that you can do:

1. Shorten a URL:
   * Send a POST request to `/` with the following two fields:
     1. `url`—the URL to shorten
     2. `key`—if the Static Key auth module is enabled
     3. `alias`—a custom alias to be used instead of a randomly-generated one
   * Example cURL command: `curl -X POST -d 'url=http://github.com/kamaln7/klein' -d 'key=secret_password' -d 'alias=klein_gh' http://localhost:5556/`
     * This will create a short URL at `http://localhost:5556/klein_gh` that redirects to `http://github.com/kamaln7/klein`.
2. Look up a URL/serve a redirect:
   * Browse to `http://[path to klein]/[alias]` to access a short URL.

## License

Copyright 2017 Kamal Nasser

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
