# go-serve

HTTP server for serving local files to the network simply.

### Motivation

- Simple design and easy configuration
- Low memory footprint
- Low startup and shutdown times
- Share files among machines easily
- Share files on embedded devices (like routers)
- Static web server for applications
- Artifact or assets server for CI/CD pipelines

### How to install

You can go to the [releases](https://github.com/eloylp/go-serve/releases/latest) page for specific OS and 
architecture requirements and download binaries.

An example install for a Linux machine could be:
```bash
sudo curl -L "https://github.com/eloylp/go-serve/releases/download/v1.3.0/go-serve_1.3.0_Linux_x86_64" \
-o /usr/local/bin/go-serve \
&& sudo chmod +x /usr/local/bin/go-serve
```

### How to use it

By default, "go-serve" command will serve the current working dir as 
its document root and serve its content in all interfaces on port 8080.

```bash
$ cd ~
$ go-serve
go-serve v1.3.0
2019/09/02 18:45:02 Starting to serve /home/user at 0.0.0.0:8080 ...
```

Of course, you can customize these parameters as in this full example:
```bash
$ go-serve -l 0.0.0.0:3000 -d /tmp -p /assets
go-serve v1.3.0
2019/09/02 18:47:02 Starting to serve /tmp at 0.0.0.0:3000 ...
```
**Note that the last option is the prefix from where the files will be served.**

#### Basic Authorization
If you need to protect your server with a basic authorization https://tools.ietf.org/html/rfc7617, 
you need to create an .htpasswd file as the following example:
```bash
htpasswd -c users-db user
go-serve -a users-db
```
**The server will watch this file for hot reload**

### How to run tests
```bash
go test ./...
```