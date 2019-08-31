# go-serve

An HTTP server for serving local files in a simple way

### How to install this binary

An example for Linux machine could be:
```bash
sudo curl -L "https://github.com/eloylp/go-serve/releases/download/v1.0.0/go-serve_1.0.0_Linux_x86_64" \
-o /usr/local/bin/go-serve \
&& sudo chmod +x /usr/local/bin/go-serve
```

### How to use it

By default "go-serve" command will serve the current working dir as 
its document root and serve its content in all interfaces on port 8080.
Of course you can customize this parameters as in this example:

```bash
go-serve -l 0.0.0.0:3000 -d /home/me/Downloads
```