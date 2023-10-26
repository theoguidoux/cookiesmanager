# Demo

To run this simply do the following:

1. `docker-compose up`
2. `curl localhost:8001 --cookie "cookieName1=value1" --cookie "cookie1=value" --cookie 'cookie2:"foo=bar|oof=rab"'`

The whoami service will simply return information about the request it
received. You can use this to play with the plugin and see how it works.

The config tells the plugin to make sure the cookie `cookie1` contains `foo=bar` and `SameSite` to contains `none`. It also makes sure that `cookie2` does not contain `foo=bar`.

```
  middlewares:
    cookiesmanager:
      plugin:
        cookiesmanager:
          adder:
            - name: cookie1
              value: "foo=bar"
            - name: "SameSite"
              value: "none"
          remover:
            - name: "cookie2"
              value: "foo=bar"
```

After running the curl request above you will see output like below. 

```
Hostname: 6e144267c95b
IP: 127.0.0.1
IP: 172.25.0.2
RemoteAddr: 172.25.0.3:33718
GET / HTTP/1.1
Host: localhost:8001
User-Agent: curl/7.88.1
Accept: */*
Accept-Encoding: gzip
Cookie: cookieName1=value1; cookie1="value foo=bar"; SameSite=none
X-Forwarded-For: 172.25.0.1
X-Forwarded-Host: localhost:8001
X-Forwarded-Port: 8001
X-Forwarded-Proto: http
X-Forwarded-Server: d76a267f2e26
X-Real-Ip: 172.25.0.1
```
