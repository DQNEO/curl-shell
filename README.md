# Usage

```
$ curl-shell
> set-base-url https://httpbin.org/
> get /get
{
  "args": {},
  "headers": {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
    "Accept-Encoding": "gzip, deflate, br",
 ...
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"
  },
  "url": "https://httpbin.org/get"
}
> post /post '{"foo":"bar"}'

> put /put '{"foo":"bar"}'

> delete /delete
```
