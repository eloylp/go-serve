# Go Serve

Just a static HTTP server with some vitamins.
<p align="center">
<img src="https://raw.githubusercontent.com/lidiackr/gophers/main/eloylp/go-serve/head.png" alt="go-serve" width="300"/>
</p>

## Table of contents

1. [Main features](#main-features)
2. [Binary distributions](#binary-distributions)
3. [Docker images](#docker-images)
4. [Use cases](#use-cases)
    1. [Upload](#upload-file)
    2. [Download](#download-file)
    3. [Upload tar.gz file](#upload-targz-archive)
    4. [Download a directory](#download-a-directory)
5. [Configuration](#configuration)
    1. [Setting up authorization](#setting-up-authorization)
6. [Prometheus metrics](#prometheus-metrics)
7. [The status endpoint](#the-status-endpoint)
8. [Security notes](#security-notes)
9. [Contributing](./CONTRIBUTING.md)

### Main features

* Serve specified folder via the HTTP protocol. Serve the current working directory by default.
* Configure auth for `READ` and `WRITE` operations independently.
* Upload single files.
* Upload an entire directory tree by using `tar.gz`  archive format and specify in server extraction point.
* Download the desired directory tree by using `targ.gz`  archives.
* Basic Prometheus metrics out of the box. Histograms for request duration, response size and upload size.
* Option to serve metrics at an alternative port.
* Status endpoint.
* Cache. Natively provided by the the Go [fileserve](https://github.com/golang/go/blob/acb189ea59d7f47e5db075e502dcce5eac6571dc/src/net/http/fs.go#L838) handler. It uses [If-Modified-Since](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since) header for caching.

### Binary distributions

The [releases](https://github.com/eloylp/go-serve/releases/latest) page will host binaries for the different OS and architectures.

Here is an example install for a Linux machine could be:

```bash
sudo curl -L "https://github.com/eloylp/go-serve/releases/download/v2.0.0/go-serve_2.0.0_Linux_x86_64" \
-o /usr/local/bin/go-serve \
&& sudo chmod +x /usr/local/bin/go-serve
```

Environment variables are the chosen method for [configuration](#configuration). 

### Docker images

There's an available docker image at [eloylp/go-serve](https://hub.docker.com/r/eloylp/go-serve) docker hub repository. You can get a
functional server, serving the specified content root just by:

```bash
docker run --rm \
 -e GOSERVE_DOC_ROOT=/mnt/data \
 -p 8080:8080 \
 -v $(pwd):/mnt/data \
  eloylp/go-serve
```

Environment variables are the chosen method for [configuration](#configuration).

### Use cases

This section will explain some common use cases that are currently covered by Go Serve.

#### Upload file

The file can simply be pushed to the **upload endpoint**. There is more information about  how to configure such endpoint in the [configuration](#configuration) section. Here is an example:

```bash
curl -X POST --location "http://localhost:8080/upload" \
    -H "GoServe-Deploy-Path: /notes.txt" \
    -H "Content-Type: application/octet-stream" \
    --data-binary @tests/root/notes/notes.txt
```

The `GoServe-Deploy-Path` value its always relative to the document root. It points the file in the serve where the bytes must be dropped.

#### Download file

Once service is up and running, files can be fetched as usual:

```bash
curl -X GET --location "http://localhost:8080/v1.2.3/gnu.png" \
    --output ./gnu.png
```

#### Upload tar.gz archive

Entire file and folder trees can be uploaded by creating a `tar.gz` archive. That archive can be pushed the designated **upload
endpoint**. Read how to configure such endpoint in the [configuration](#configuration) section.

```bash
curl -X POST --location "http://localhost:8080/upload" \
    -H "GoServe-Deploy-Path: /v1.2.3" \
    -H "Content-Type: application/tar+gzip" \
    --data-binary @tests/doc-root.tar.gz
```

The `GoServe-Deploy-Path` value its always relative to the document root. It points where the content of the `tar.gz` archive must be extracted.

#### Download a directory

Entire file and folder trees can be downloaded from the server. Fetching the **download endpoint** and requesting the server what type of archive you would like to get. Currently, only `tar.gz` is supported. Read how to configure such endpoint in the [configuration](#configuration) section:

```bash
curl -X GET --location "http://localhost:8080/download" \
    -H "GoServe-Download-Path: /v1.2.3" \
    -H "Accept: application/tar+gzip" \
    --output ./v1.2.4.tar.gz
```

The `GoServe-Download-Path` value its always relative to the document root.

### Configuration

Go serve uses environment variables to configure its internals. Here is a table of the current customizable parts of the server:

| Variable                                 | Description                                                  | Default                                                      |
| ---------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| GOSERVE_LISTEN_ADDR                      | The socket where the server will listen for connections.     | "0.0.0.0:8080"                                               |
| GOSERVE_DOC_ROOT                         | Path to the  document root its going to be served.           | "."                                                          |
| GOSERVE_PREFIX                           | The prefix path under all files will be served. Default value is "/static"  so all files will be served under such path i.e "/static/notes.txt" . This is mandatory and should not interfere with other configured paths. | "/static"                                                    |
| GOSERVE_UPLOAD_ENDPOINT                  | The path in the server where all uploads will take place. If not defined, it will be disabled. By default is **disabled**. | ""                                                           |
| GOSERVE_DOWNLOAD_ENDPOINT                | The path in the server where all downloads will take place. If not defined, it will be disabled. By default is **disabled**. | ""                                                           |
| GOSERVE_SHUTDOWN_TIMEOUT                 | The number of seconds that the server will wait to terminate pending active connections before closing. | "5s"                                                         |
| GOSERVE_READ_TIMEOUT                     | The maximum duration for reading the entire request, including the body. Default is **unlimited**. | "0s"                                                         |
| GOSERVE_WRITE_TIMEOUT                    | The maximum duration before timing out writes of the response. Default is **unlimited**. | "0s"                                                         |
| GOSERVE_READ_AUTHORIZATIONS              | Configures which users are allowed to make idempotent requests to the server. It expects a **base64** string containing a users table generated by the **htpasswd** utility. By default, read authorization is **disabled** so all users can read the entire server. See [authorization](#setting-up-authorization) for more details. | ""                                                           |
| GOSERVE_WRITE_AUTHORIZATIONS             | Configures which users are allowed to make  **non** idempotent requests to the server. It expects a **base64** string containing a users table generated by the **htpasswd** utility. By default, write authorization is **disabled** so unauthorized users can upload files if the  **GOSERVE_UPLOAD_ENDPOINT** variable is defined. See [authorization](#setting-up-authorization) for more details. | ""                                                           |
| GOSERVE_METRICS_ENABLED                  | Configures if the Prometheus metrics are enabled or disabled. | true                                                         |
| GOSERVE_METRICS_PATH                     | Configures in which endpoint the metrics should be served. This can help to hide the metrics endpoint by introducing a more complicated path that only systems will know. | "/metrics"                                                   |
| GOSERVE_METRICS_LISTEN_ADDR              | If configured, another sidecar server will be configured exclusively for serving metrics. This is **disabled** by default. An example of value could be: "0.0.0.0:9091" . | ""                                                           |
| GOSERVE_METRICS_REQUEST_DURATION_BUCKETS | Define the buckets for the histogram of request duration. Expressed in seconds. | "0.005,0.01,0.025, 0.05,0.1,0.25,0.5, 1,2.5,5,10"            |
| GOSERVE_METRICS_SIZE_BUCKETS             | Define the buckets for the histogram of response size and upload size. Expressed in bytes. | "64,256,1024,4096,16384,65536,262144,1048576,4194304,16777216" |

#### Setting up authorization

Both type of authorizations, *GOSERVE_READ_AUTHORIZATIONS* and  *GOSERVE_WRITE_AUTHORIZATIONS* are configured in the same manner. Those variables expect a **base64** encoded file generated by the tool [**htpasswd**](https://httpd.apache.org/docs/2.4/programs/htpasswd.html) .
The passwords must be encrypted by using the **bcrypt** algorithm. The following is an example for creating such value for the user "Alice" with password "password":

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

The authorization fronted in compatible with [rfc7617]( https://tools.ietf.org/html/rfc7617) basic authorization scheme. This is an example
from a curl request:

```bash
curl -X GET --location "http://localhost:8080/v1.2.3/gnu.png" \
    --basic --user alice:password \
    --output ./gnu.png
```

### Prometheus metrics

By default, this server provides various [histograms](https://prometheus.io/docs/practices/histograms/) that will provide a good global view of server operations. You can scrape this metrics at `/metrics` once the server was started. It is possible to have a sidecar HTTP server dedicated to metrics. See the [configuration](#configuration) section for more details. The following is an excerpt of the available metrics:

```prometheus
# HELP http_upload_size Histogram to represent the successful uploads to the server
# TYPE http_upload_size histogram
http_upload_size_bucket{le="64"} 0
http_upload_size_bucket{le="256"} 0
http_upload_size_bucket{le="1024"} 0
http_upload_size_bucket{le="4096"} 0
http_upload_size_bucket{le="16384"} 0
http_upload_size_bucket{le="65536"} 0
http_upload_size_bucket{le="262144"} 0
http_upload_size_bucket{le="1.048576e+06"} 1
http_upload_size_bucket{le="4.194304e+06"} 1
http_upload_size_bucket{le="1.6777216e+07"} 1
http_upload_size_bucket{le="+Inf"} 1
http_upload_size_sum 533766
http_upload_size_count 1
# HELP http_request_duration_seconds 
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.005"} 0
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.01"} 0
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.025"} 0
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.05"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.1"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.25"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="0.5"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="1"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="2.5"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="5"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="10"} 1
http_request_duration_seconds_bucket{code="200",endpoint="/upload",method="POST",le="+Inf"} 1
http_request_duration_seconds_sum{code="200",endpoint="/upload",method="POST"} 0.027193285
http_request_duration_seconds_count{code="200",endpoint="/upload",method="POST"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.005"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.01"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.025"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.05"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.1"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.25"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="0.5"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="1"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="2.5"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="5"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="10"} 1
http_request_duration_seconds_bucket{code="304",endpoint="/static",method="GET",le="+Inf"} 1
http_request_duration_seconds_sum{code="304",endpoint="/static",method="GET"} 0.00024055
http_request_duration_seconds_count{code="304",endpoint="/static",method="GET"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.005"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.01"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.025"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.05"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.1"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.25"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="0.5"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="1"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="2.5"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="5"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="10"} 1
http_request_duration_seconds_bucket{code="404",endpoint="/static",method="GET",le="+Inf"} 1
http_request_duration_seconds_sum{code="404",endpoint="/static",method="GET"} 0.000128526
http_request_duration_seconds_count{code="404",endpoint="/static",method="GET"} 1
# HELP http_response_size 
# TYPE http_response_size histogram
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="64"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="256"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="1024"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="4096"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="16384"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="65536"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="262144"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="1.048576e+06"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="4.194304e+06"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="1.6777216e+07"} 1
http_response_size_bucket{code="200",endpoint="/upload",method="POST",le="+Inf"} 1
http_response_size_sum{code="200",endpoint="/upload",method="POST"} 39
http_response_size_count{code="200",endpoint="/upload",method="POST"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="64"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="256"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="1024"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="4096"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="16384"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="65536"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="262144"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="1.048576e+06"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="4.194304e+06"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="1.6777216e+07"} 1
http_response_size_bucket{code="304",endpoint="/static",method="GET",le="+Inf"} 1
http_response_size_sum{code="304",endpoint="/static",method="GET"} 0
http_response_size_count{code="304",endpoint="/static",method="GET"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="64"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="256"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="1024"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="4096"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="16384"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="65536"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="262144"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="1.048576e+06"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="4.194304e+06"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="1.6777216e+07"} 1
http_response_size_bucket{code="404",endpoint="/static",method="GET",le="+Inf"} 1
http_response_size_sum{code="404",endpoint="/static",method="GET"} 19
http_response_size_count{code="404",endpoint="/static",method="GET"} 1
```

### The status endpoint

The most probably use of this server to place it behind a reverse proxy. In order to facilitate the readiness of the service a status endpoint under `/status` is provided. Here is an example of the information provided:

```json
{
  "status": "ok",
  "info": {
    "name": "go-serve",
    "version": "v2.0.0-rc8",
    "build": "5b7aaf3",
    "build_time": "2021-05-15T17:07:03+0000"
  }
}
```
### Security notes

This server assumes that users with write permissions are trusted ones. They will be able to upload any kind of file to the server document root. Please, if you enable uploads, be sure that you configure write [authorization](#setting-up-authorization) .

