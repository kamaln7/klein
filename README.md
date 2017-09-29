<p align="center">
  <img src="/klein.png" alt="klein logo" />
</p>

klein is a URL shortener. Usage:

```
Usage of klein:
  -key string
    	upload API Key
  -length int
    	code length (default 3)
  -listenAddr string
    	listen address (default "127.0.0.1:5556")
  -path string
    	path to urls (default "/srv/www/urls/")
  -root string
    	root redirect
  -template string
    	path to error template (default "./404.html")
  -url string
    	path to public facing url (default "http://127.0.0.1:5556/")
```

The only _absolutely_ necessary options are `template`, `path`, and `url`.
