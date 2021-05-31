# Go Serve

Just a static HTTP server with some vitamins.
<p align="center">
<img src="art/gopher.png" alt="go-serve" width="300"/>
</p>

<p align="right">
<span style="color:silver">From the original gopher by Renee French.
</span>
</p>

## Table of contents

1. [Main features](#main-features)
2. [Binary distributions](#binary-distributions)
3. [Docker images](#docker-images)
4. [Use cases](#use-cases)
    1. [Upload tar.gz file](#upload-targz-file)
    2. [File download](#ordinary-file-download)
    3. [Download a directory](#download-a-directory)
5. [Configuration](#configuration)
    1. [Setting up authorization](#setting-up-authorization)
6. [Prometheus metrics](#prometheus-metrics)
7. [The status endpoint](#the-status-endpoint)
8. [Contributing](./CONTRIBUTING.md)

### Main features

* Serve specified folder via the HTTP protocol.
* Add users authorization for `READ` and `WRITE` operations independently.
* Upload `tar.gz` files and extract them under the specified path in the document root.
* Download folders and files of your document root using `tar.gz` files as archive.
* Basic Prometheus metrics out of the box.
* Option to serve metrics on an alternative port.
* Status endpoint.

### Binary distributions

You can go to the [releases](https://github.com/eloylp/go-serve/releases/latest) page for specific OS and architecture requirements and
download binaries.

An example install for a Linux machine could be:

```bash
sudo curl -L "https://github.com/eloylp/go-serve/releases/download/v2.0.0/go-serve_2.0.0_Linux_x86_64" \
-o /usr/local/bin/go-serve \
&& sudo chmod +x /usr/local/bin/go-serve
```

Environment vars are the chosen method for configuration. See this section for more info about [configuration](#configuration).

### Docker images

There's an available docker image at [eloylp/go-serve](https://hub.docker.com/r/eloylp/go-serve) docker hub repository. You can get a
functional server, serving the current content root just by:

```bash
docker run --rm \
 -e GOSERVE_DOCROOT=/mnt/data \
 -p 8080:8080 \
 -v $(pwd):/mnt/data \
  eloylp/goserve
```

As you may notice, environment vars are the chosen method for configuration. See this section for more info
about [configuration](#configuration).

### Use cases

This section will explain some common use cases that are currently covered by Go Serve.

#### Upload tar.gz file

You can upload files to the document root of the server at runtime. Just create your tar.gz and push it to the designated **upload
endpoint**. Read how to configure such endpoint in the [configuration](#configuration) section.

```bash
curl -X POST --location "http://localhost:8080/upload" \
    -H "GoServe-Deploy-Path: /v1.2.3" \
    -H "Content-Type: application/tar+gzip" \
    -d @tests/doc-root.tar.gz
```

The `GoServe-Deploy-Path` value its always relative to the document root.

#### Ordinary file download

Once service is up and running, you can fetch resources as usual you will do with any HTTP server:

```bash
curl -X GET --location "http://localhost:8080/v1.2.3/gnu.png" \
    --output ./gnu.png
```

#### Download a directory

You can download a directory by just fetching the **download endpoint** and requesting the server what type of archive you would like to
get. Currently, only `tar.gz` is supported. Read how to configure such endpoint in the [configuration](#configuration) section:

```bash
curl -X GET --location "http://localhost:8080/download" \
    -H "GoServe-Download-Path: /v1.2.3" \
    -H "Accept: application/tar+gzip" \
    --output ./v1.2.4.tar.gz
```

The `GoServe-Download-Path` value its always relative to the document root.

### Configuration

Go serve uses environment variables to configure its internals. Here is a table of the current customizable parts of the server:

| Variable                                 | Description                                                  | Default                                           |
| ---------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------- |
| GOSERVE_LISTEN_ADDR                      | The socket where the server will listen for connections.     | "0.0.0.0:8080"                                    |
| GOSERVE_DOC_ROOT                         | Path to the  document root we are going to serve.            | "."                                               |
| GOSERVE_PREFIX                           | The prefix path under all files will be served. Defaults in value is "/static"  so all files will be served under such path i.e "/static/notes.txt" . This is mandatory and should not interfere with other configured paths. | "/static"                                         |
| GOSERVE_UPLOAD_ENDPOINT                  | The path in the server where all uploads will take place. If not defined, it will be disabled. By default is **
disabled** . | ""                                                |
| GOSERVE_DOWNLOAD_ENDPOINT                | The path in the server where all downloads will take place. If not defined, it will be disabled. By default is **
disabled** . | ""                                                |
| GOSERVE_SHUTDOWN_TIMEOUT                 | The number of seconds that the server will wait to terminate pending active connections before closing. | "5s"                                              |
| GOSERVE_READ_TIMEOUT                     | The maximum duration for reading the entire request, including the body. Default is **
unlimited**. | "0s"                                              |
| GOSERVE_WRITE_TIMEOUT                    | The maximum duration before timing out writes of the response. Default is **
unlimited**. | "0s"                                              |
| GOSERVE_READ_AUTHORIZATIONS              | Configures which users are allowed to make idempotent requests to the server. It expects a **
base64** string containing a users table generated by the **htpasswd** utility. By default, read authorization is **
disabled** so all users can read the entire server. See [authorization](#setting-up-authorization) for more details. | ""                                                |
| GOSERVE_WRITE_AUTHORIZATIONS             | Configures which users are allowed to make  *
non* idempotent requests to the server. It expects a **base64** string containing a users table generated by the **
htpasswd** utility. By default, write authorization is **disabled** so unauthorized users can upload files if the  **
GOSERVE_UPLOADENDPOINT** variable is defined. See [authorization](#setting-up-authorization) for more details. | ""                                                |
| GOSERVE_METRICS_ENABLED                  | Configures if the Prometheus metrics are enabled or disabled. | true                                              |
| GOSERVE_METRICS_PATH                     | Configures in which endpoint the metrics should be served. This can help to hide the metrics endpoint by introducing a more complicated path that only systems will know. | "/metrics"                                        |
| GOSERVE_METRICS_LISTEN_ADDR              | If configured, another sidecar server will be configured exclusively for serving metrics. This is **
disabled** by default. An example of value could be: "0.0.0.0:9091" . | ""                                                |
| GOSERVE_METRICS_REQUEST_DURATION_BUCKETS | Default metrics on this server is a histogram of request duration time. Here a user can customize the buckets where distribution ranges are going to be defined. | "0.005,0.01,0.025, 0.05,0.1,0.25,0.5, 1,2.5,5,10" |

#### Setting up authorization

Both type of authorizations, *GOSERVE_READAUTHORIZATIONS* and  *GOSERVE_WRITEAUTHORIZATIONS* are configured in the same manner. Those
variables expects a **base64** encoded file generated by the tool [**htpasswd**](https://httpd.apache.org/docs/2.4/programs/htpasswd.html) .
The passwords must be encrypted by using the **bcrypt** algorithm. The following is an example for creating such value for the user "alice"
with password "password":

```bash
$ htpasswd -B -c auth.txt alice
New password: password           ## note this is an interactive step
Re-type new password: password   ## note this is an interactive step
Adding password for user alice

$ cat auth.txt  ## Check the content of the file
alice:$2y$05$039J5egx9S9ayeGQTYQ5nex3SmMuXho7oXbIMInW9EX9UIywjIJJa

## Create the needed value for GOSERVE_READAUTHORIZATIONS or GOSERVE_WRITEAUTHORIZATIONS
$ cat auth.txt | base64
YWxpY2U6JDJ5JDA1JDAzOUo1ZWd4OVM5YXllR1FUWVE1bmV4M1NtTXVYaG83b1hiSU1Jblc5RVg5
VUl5d2pJSkphCg==
```

Once the file is created following the above steps, we can also add more users with the following command and continue with the same
previous steps.

```bash
htpasswd -B auth.txt bob
```

F.A.Q: In [Kubernetes secrets](https://kubernetes.io/es/docs/concepts/configuration/secret/) you need to double encode in base64 the value,
as Kubernetes requires to wrap all the secrets in this encoding.

#### Using authorization in requests

The authorization frontend in compatible with [rfc7617]( https://tools.ietf.org/html/rfc7617) basic authorization scheme. This is an example
from a curl request:

```bash
curl -X GET --location "http://localhost:8080/v1.2.3/gnu.png" \
    --basic --user alice:password \
    --output ./gnu.png
```

### Prometheus metrics

By default this server provides basic Prometheus metrics. It includes an [histogram](https://prometheus.io/docs/practices/histograms/) that
represents the request duration in seconds. By default you can scrape this metrics at `/metrics` once the server was started. It is possible
to have a sidecar HTTP server dedicated to metrics. See the [configuration](#configuration) section for more details.