# ChartMuseum

[![Codefresh build status](https://g.codefresh.io/api/badges/pipeline/chartmuseum/helm%2Fchartmuseum%2Fmaster?type=cf-1)](https://g.codefresh.io/public/accounts/chartmuseum/pipelines/helm/chartmuseum/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/helm/chartmuseum)](https://goreportcard.com/report/github.com/helm/chartmuseum)
[![GoDoc](https://godoc.org/github.com/helm/chartmuseum?status.svg)](https://godoc.org/github.com/helm/chartmuseum)

<p align="center"><img align="center" src="logo2.png"></p><br/>

*ChartMuseum* is an open-source **[Helm Chart Repository](https://github.com/helm/helm-www/blob/master/content/docs/topics/chart_repository.md)** server written in Go (Golang), with support for cloud storage backends, including [Google Cloud Storage](https://cloud.google.com/storage/), [Amazon S3](https://aws.amazon.com/s3/), [Microsoft Azure Blob Storage](https://azure.microsoft.com/en-us/services/storage/blobs/), [Alibaba Cloud OSS Storage](https://www.alibabacloud.com/product/oss), [Openstack Object Storage](https://developer.openstack.org/api-ref/object-store/), [Oracle Cloud Infrastructure Object Storage](https://cloud.oracle.com/storage), [Baidu Cloud BOS Storage](https://cloud.baidu.com/product/bos.html), [Tencent Cloud Object Storage](https://intl.cloud.tencent.com/product/cos), [DigitalOcean Spaces](https://www.digitalocean.com/products/spaces/), [Minio](https://min.io/), and [etcd](https://etcd.io/).

Works as a valid Helm Chart Repository, and also provides an API for uploading charts.

<img width="120" align="right" src="https://github.com/golang-samples/gopher-vector/raw/master/gopher-side_color.png">
<img width="40" align="right" src="https://github.com/golang-samples/gopher-vector/raw/master/gopher-side_color.png">

Powered by some great Go technology:
- [helm/helm](https://github.com/helm/helm) - for working with charts
- [gin-gonic/gin](https://github.com/gin-gonic/gin) - for HTTP routing
- [urfave/cli](https://github.com/urfave/cli) - for command line option parsing
- [spf13/viper](https://github.com/spf13/viper) - for configuration
- [uber-go/zap](https://github.com/uber-go/zap) - for logging
- [chartmuseum/auth](https://github.com/chartmuseum/auth) - for auth
- [chartmuseum/storage](https://github.com/chartmuseum/storage) - for multi-cloud storage

## API

### Helm Chart Repository
- `GET /index.yaml` - retrieved when you run `helm repo add chartmuseum http://localhost:8080/`
- `GET /charts/mychart-0.1.0.tgz` - retrieved when you run `helm install chartmuseum/mychart`
- `GET /charts/mychart-0.1.0.tgz.prov` - retrieved when you run `helm install` with the `--verify` flag

### Chart Manipulation
- `POST /api/charts` - upload a new chart version
- `POST /api/prov` - upload a new provenance file
- `DELETE /api/charts/<name>/<version>` - delete a chart version (and corresponding provenance file)
- `GET /api/charts` - list all charts
- `GET /api/charts/<name>` - list all versions of a chart
- `GET /api/charts/<name>/<version>` - describe a chart version
- `HEAD /api/charts/<name>` - check if chart exists (any versions)
- `HEAD /api/charts/<name>/<version>` - check if chart version exists

### Server Info
- `GET /` - HTML welcome page
- `GET /health` - returns 200 OK

## Uploading a Chart Package
<sub>*Follow **"How to Run"** section below to get ChartMuseum up and running at ht<span>tp:/</span>/localhost:8080*<sub>

First create `mychart-0.1.0.tgz` using the [Helm CLI](https://docs.helm.sh/using_helm/#installing-helm):
```
cd mychart/
helm package .
```

Upload `mychart-0.1.0.tgz`:
```bash
curl --data-binary "@mychart-0.1.0.tgz" http://localhost:8080/api/charts
```

If you've signed your package and generated a [provenance file](https://github.com/helm/helm/blob/master/docs/provenance.md), upload it with:
```bash
curl --data-binary "@mychart-0.1.0.tgz.prov" http://localhost:8080/api/prov
```

Both files can also be uploaded at once (or one at a time) on the `/api/charts` route using the `multipart/form-data` format:

```bash
curl -F "chart=@mychart-0.1.0.tgz" -F "prov=@mychart-0.1.0.tgz.prov" http://localhost:8080/api/charts
```

You can also use the [helm-push plugin](https://github.com/chartmuseum/helm-push):
```
helm push mychart/ chartmuseum
```

## Installing Charts into Kubernetes
Add the URL to your *ChartMuseum* installation to the local repository list:
```bash
helm repo add chartmuseum http://localhost:8080
```

Search for charts:
```bash
helm search chartmuseum/
```

Install chart:
```bash
helm install chartmuseum/mychart
```

## How to Run
### CLI
#### Installation
Install binary using [GoFish](https://gofi.sh/):
```
gofish install chartmuseum
==> Installing chartmuseum...
🐠  chartmuseum 0.11.0: installed in 95.431145ms
```

or manually:
```bash
# on Linux
curl -LO https://s3.amazonaws.com/chartmuseum/release/latest/bin/linux/amd64/chartmuseum

# on macOS
curl -LO https://s3.amazonaws.com/chartmuseum/release/latest/bin/darwin/amd64/chartmuseum

# on Windows
curl -LO https://s3.amazonaws.com/chartmuseum/release/latest/bin/windows/amd64/chartmuseum

chmod +x ./chartmuseum
mv ./chartmuseum /usr/local/bin
```
Using `latest` in URLs above will get the latest binary (built from master branch).

Replace `latest` with `$(curl -s https://s3.amazonaws.com/chartmuseum/release/stable.txt)` to automatically determine the latest stable release (e.g. `v0.11.0`). The stable checksums can be found [here](https://github.com/fishworks/fish-food/blob/master/Food/chartmuseum.lua).

Determine your version with `chartmuseum --version`.

#### Configuration
Show all CLI options with `chartmuseum --help`. Common configurations can be seen below.

All command-line options can be specified as environment variables, which are defined by the command-line option, capitalized, with all `-`'s replaced with `_`'s.

For example, the env var `STORAGE_AMAZON_BUCKET` can be used in place of `--storage-amazon-bucket`.

##### Using config yaml file example
When using a yaml file instead of the command line， Convert the `-` of the parameter in the command to `.`.
Note that `--storage` is converted to `storage.backend`
```yaml
debug: true
port: 8080
storage.backend: local
storage.local.rootdir: <storage_path>
bearerauth: 1
authrealm: <authorization server url>
authservice: <authorization server service name>
authcertpath: <path to authorization server public pem file> 
depth: 2

```

#### Using with Amazon S3 or Compatible services like Minio or DigitalOcean.
Make sure your environment is properly setup to access `my-s3-bucket`

For Amazon S3, `endpoint` is automatically inferred.
```bash
chartmuseum --debug --port=8080 \
  --storage="amazon" \
  --storage-amazon-bucket="my-s3-bucket" \
  --storage-amazon-prefix="" \
  --storage-amazon-region="us-east-1"
```

For S3 compatible services like Minio, set the credentials using environment variables and pass the `endpoint`.

```bash
export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
chartmuseum --debug --port=8080 \
  --storage="amazon" \
  --storage-amazon-bucket="my-s3-bucket" \
  --storage-amazon-prefix="" \
  --storage-amazon-region="us-east-1" \
  --storage-amazon-endpoint="my-s3-compatible-service-endpoint"
```

You need at least the following permissions inside your IAM Policy
```yaml
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowListObjects",
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::my-s3-bucket"
    },
    {
      "Sid": "AllowObjectsCRUD",
      "Effect": "Allow",
      "Action": [
        "s3:DeleteObject",
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::my-s3-bucket/*"
    }
  ]
}
```  

For DigitalOcean, set the credentials using environment variable and pass the `endpoint`.  
Note below, that the region `us-east-1` needs to be set, since that is how the DigitalOcean cli implementation functions. The actual region of your spaces location is defined by the endpoint. Below we are using Frankfurt as an example.
```bash
export AWS_ACCESS_KEY_ID="spaces_access_key"
export AWS_SECRET_ACCESS_KEY="spaces_secret_key"
  chartmuseum --debug --port=8080 \
  --storage="amazon" \
  --storage-amazon-bucket="my_spaces_name" \
  --storage-amazon-prefix="my_spaces_name_subfolder" \
  --storage-amazon-region="us-east-1" \
  --storage-amazon-endpoint="https://fra1.digitaloceanspaces.com"
```
The acces_key and secret_key can be generated from the DigitalOcean console, under the section API/Spaces_access_keys.  

#### Using with Google Cloud Storage
Make sure your environment is properly setup to access `my-gcs-bucket`.

One way to do so is to set the `GOOGLE_APPLICATION_CREDENTIALS` var in your environment, pointing to the JSON file containing your service account key:
```
export GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/[FILE_NAME].json"
```

More info on Google Cloud authentication can be found [here](https://cloud.google.com/docs/authentication/getting-started).

```bash
chartmuseum --debug --port=8080 \
  --storage="google" \
  --storage-google-bucket="my-gcs-bucket" \
  --storage-google-prefix=""
```

#### Using with Microsoft Azure Blob Storage

Make sure your environment is properly setup to access `mycontainer`.

To do so, you must set the following env vars:
- `AZURE_STORAGE_ACCOUNT`
- `AZURE_STORAGE_ACCESS_KEY`

```bash
chartmuseum --debug --port=8080 \
  --storage="microsoft" \
  --storage-microsoft-container="mycontainer" \
  --storage-microsoft-prefix=""
```

#### Using with Alibaba Cloud OSS Storage

Make sure your environment is properly setup to access `my-oss-bucket`.

To do so, you must set the following env vars:
- `ALIBABA_CLOUD_ACCESS_KEY_ID`
- `ALIBABA_CLOUD_ACCESS_KEY_SECRET`

```bash
chartmuseum --debug --port=8080 \
  --storage="alibaba" \
  --storage-alibaba-bucket="my-oss-bucket" \
  --storage-alibaba-prefix="" \
  --storage-alibaba-endpoint="oss-cn-beijing.aliyuncs.com"
```

#### Using with Openstack Object Storage

Make sure your environment is properly setup to access `mycontainer`.

To do so, you must set the following env vars (depending on your openstack version):
- `OS_AUTH_URL`
- either `OS_PROJECT_NAME` or `OS_TENANT_NAME` or `OS_PROJECT_ID` or `OS_TENANT_ID`
- either `OS_DOMAIN_NAME` or `OS_DOMAIN_ID`
- either `OS_USERNAME` or `OS_USERID`
- `OS_PASSWORD`

```bash
chartmuseum --debug --port=8080 \
  --storage="openstack" \
  --storage-openstack-container="mycontainer" \
  --storage-openstack-prefix="" \
  --storage-openstack-region="myregion"
```

#### Using with Oracle Cloud Infrastructure Object Storage

Make sure your environment is properly setup to access `my-ocs-bucket`.

More info on Oracle Cloud Infrastructure authentication can be found [here](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/apisigningkey.htm).

```bash
chartmuseum --debug --port=8080 \
  --storage="oracle" \
  --storage-oracle-bucket="my-ocs-bucket" \
  --storage-oracle-prefix="" \
  --storage-oracle-compartmentid="ocid1.compartment.oc1..1234"
```

#### Using with Baidu Cloud BOS Storage

Make sure your environment is properly setup to access `my-bos-bucket`.

To do so, you must set the following env vars:
- `BAIDU_CLOUD_ACCESS_KEY_ID`
- `BAIDU_CLOUD_ACCESS_KEY_SECRET`

```bash
chartmuseum --debug --port=8080 \
  --storage="baidu" \
  --storage-baidu-bucket="my-bos-bucket" \
  --storage-baidu-prefix="" \
  --storage-baidu-endpoint="bj.bcebos.com"
```

#### Using with Tencent Cloud COS Storage

Make sure your environment is properly setup to access `my-cos-bucket`.

To do so, you must set the following env vars:
- `TENCENT_CLOUD_COS_SECRET_ID`
- `TENCENT_CLOUD_COS_SECRET_KEY`

```bash
chartmuseum --debug --port=8080 \
  --storage="tencent" \
  --storage-tencent-bucket="my-cos-bucket" \
  --storage-tencent-prefix="" \
  --storage-tencent-endpoint="cos.ap-beijing.myqcloud.com"
```

#### Using with etcd

To use etcd as backend you need the CA certificate and the signed key pair.
See [here](https://coreos.com/etcd/docs/latest/op-guide/security.html)

```bash
chartmuseum --debug --port=8080 \
  --storage="etcd" \
  --storage-etcd-cafile="/path/to/ca.crt" \
  --storage-etcd-certfile="/path/to/server.crt" \
  --storage-etcd-keyfile="/path/to/server.key" \
  --storage-etcd-prefix="" \
  --storage-etcd-endpoint="http://localhost:2379"
```

#### Using with local filesystem storage
Make sure you have read-write access to `./chartstorage` (will create if doesn't exist on first upload)
```bash
chartmuseum --debug --port=8080 \
  --storage="local" \
  --storage-local-rootdir="./chartstorage"
```

#### Basic Auth
If both of the following options are provided, basic http authentication will protect all routes:
- `--basic-auth-user=<user>` - username for basic http authentication
- `--basic-auth-pass=<pass>` - password for basic http authentication

You may want basic auth to only be applied to operations that can change Charts, i.e. PUT, POST and DELETE.  So to avoid basic auth on GET operations use

- `--auth-anonymous-get` - allow anonymous GET operations

#### Bearer/Token Auth

If all of the following options are provided, bearer auth will protect all routes:
- `--bearer-auth` - enables bearer auth
- `--auth-realm=<realm>` - authorization server url
- `--auth-service=<service>` - authorization server service name
- `--auth-cert-path=<path>` - path to authorization server public pem file

Using options above, *ChartMuseum* is configured with a public key, and will accept RS256 JWT tokens signed by the associated private key, passed in the `Authorization` header. You can use the [chartmuseum/auth](https://github.com/chartmuseum/auth) Go library to generate valid JWT tokens.

In order to gain access to a specific resource, the JWT token must contain an `access` section in the claims. This section indicates which resources the user is able to access. Here is an example token payload:

```json
{
  "exp": 1543995770,
  "iat": 1543995470,
  "access": [
    {
      "type": "artifact-repository",
      "name": "org1/repo1",
      "actions": [
        "pull"
      ]
    }
  ]
}
```

The `type` is always "artifact-repository", the `name` is the namespace/tenant (just use the string "repo" if using single-tenant server), and `actions` is an array of actions the user can perform ("pull" and/or "push).

For more information about how this works, please see [chartmuseum/auth-server-example](https://github.com/chartmuseum/auth-server-example).


#### HTTPS
If both of the following options are provided, the server will listen and serve HTTPS:
- `--tls-cert=<crt>` - path to tls certificate chain file
- `--tls-key=<key>` - path to tls key file

##### HTTPS with Client Certificate Authentication
If the above HTTPS values are provided in addition to below, the server will listen and serve HTTPS and authenticate client requests against the CA certificate:
-  `--tls-ca-cert=<cacert>` - path to tls certificate file

#### Just generating index.yaml
You can specify the `--gen-index` option if you only wish to use _ChartMuseum_ to generate your index.yaml file. Note that this will only work with `--depth=0`.

The contents of index.yaml will be printed to stdout and the program will exit. This is useful if you are satisfied with your current Helm CI/CD process and/or don't want to monitor another web service.

#### Other CLI options
- `--log-json` - output structured logs as json
- `--log-health` - log incoming /health requests
- `--disable-api` - disable all routes prefixed with /api
- `--disable-delete` - explicitely disable the delete chart route
- `--disable-statefiles` - disable use of index-cache.yaml
- `--allow-overwrite` - allow chart versions to be re-uploaded without ?force querystring
- `--disable-force-overwrite` - do not allow chart versions to be re-uploaded, even with ?force querystring
- `--chart-url=<url>` - absolute url for .tgzs in index.yaml
- `--storage-amazon-endpoint=<endpoint>` - alternative s3 endpoint
- `--storage-amazon-sse=<algorithm>` - s3 server side encryption algorithm
- `--storage-openstack-cacert=<path>` - path to a custom ca certificates bundle for openstack
- `--chart-post-form-field-name=<field>` - form field which will be queried for the chart file content
- `--prov-post-form-field-name=<field>` - form field which will be queried for the provenance file content
- `--index-limit=<number>` - limit the number of parallel indexers
- `--context-path=<path>` - base context path (new root for application routes)
- `--depth=<number>` - levels of nested repos for multitenancy
- `--cors-alloworigin=<value>` - value to set in the Access-Control-Allow-Origin HTTP header

### Docker Image
Available via [Docker Hub](https://hub.docker.com/r/chartmuseum/chartmuseum/).

Example usage (local storage):
```bash
docker run --rm -it \
  -p 8080:8080 \
  -e DEBUG=1 \
  -e STORAGE=local \
  -e STORAGE_LOCAL_ROOTDIR=/charts \
  -v $(pwd)/charts:/charts \
  chartmuseum/chartmuseum:latest
```

Example usage (S3):
```bash
docker run --rm -it \
  -p 8080:8080 \
  -e DEBUG=1 \
  -e STORAGE="amazon" \
  -e STORAGE_AMAZON_BUCKET="my-s3-bucket" \
  -e STORAGE_AMAZON_PREFIX="" \
  -e STORAGE_AMAZON_REGION="us-east-1" \
  -v ~/.aws:/home/chartmuseum/.aws:ro \
  chartmuseum/chartmuseum:latest
```

### Helm Chart
There is a [Helm chart for *ChartMuseum*](https://github.com/helm/charts/tree/master/stable/chartmuseum) itself which can be found in the official Helm charts repository.

You can also view it on [Helm Hub](https://hub.helm.sh/charts/stable/chartmuseum).

To install:
```bash
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm install stable/chartmuseum
```

If interested in making changes, please submit a PR to [helm/charts](https://github.com/helm/charts). Before doing any work, please check for any [currently open pull requests](https://github.com/helm/charts/pulls?q=is%3Apr+is%3Aopen+chartmuseum). Thanks!

## Multitenancy
Multitenancy is supported with the `--depth` flag.

To begin, start with a directory structure such as
```
charts
├── org1
│   ├── repoa
│   │   └── nginx-ingress-0.9.3.tgz
├── org2
│   ├── repob
│   │   └── chartmuseum-0.4.0.tgz
```

This represents a storage layout appropriate for `--depth=2`. The organization level can be eliminated by using `--depth=1`. The default depth is 0 (singletenant server).

Start the server with `--depth=2`, pointing to the `charts/` directory:
```
chartmuseum --debug --depth=2 --storage="local" --storage-local-rootdir=./charts
```

This example will provide two separate Helm Chart Repositories at the following locations:
- `http://localhost:8080/org1/repoa`
- `http://localhost:8080/org2/repob`

This should work with all supported storage backends.

To use the chart manipulation routes, simply place the name of the repo directly after "/api" in the route:

```bash
curl -F "chart=@mychart-0.1.0.tgz" http://localhost:8080/api/org1/repoa/charts
```

You may also experiment with the `--depth-dynamic` flag, which should allow for dynamic depth levels (i.e. all of `/api/charts`, `/api/myrepo/charts`, `/api/org1/repoa/charts`).

## Pagination

For large chart repositories, you may wish to paginate the results from the `GET /api/charts` route. 

To do so, add the `offset` and `limit` query params to the request. For example, to retrieve a list of 5 charts total, skipping the first 5 charts, you could use the following:

```
GET /api/charts?offset=5&limit=5
```

## Cache

By default, the contents of `index.yaml` (per-tenant) will be stored in memory. This means that memory usage will continue to grow indefinitely as more charts are added to storage.

You may wish to offload this to an external cache store, especially for large, multitenant installations.

### Using Redis

Example of using Redis as an external cache store:
```bash
chartmuseum --debug --port=8080 \
  --storage="local" \
  --storage-local-rootdir="./chartstorage" \
  --cache="redis" \
  --cache-redis-addr="localhost:6379" \
  --cache-redis-password="" \
  --cache-redis-db=0
```


## Prometheus Metrics

ChartMuseum exposes its [Prometheus metrics](https://prometheus.io/docs/concepts/metric_types/) at the `/metrics` route on the main port. This can be disabled with the `--disable-metrics` command-line flag or the `DISABLE_METRICS` environment variable.

> Note that the Kubernetes chart currently disables metrics by default (`DISABLE_METRICS=true` is set in the chart).

Below are the current application metrics exposed. Note that there is a per tenant (repo) label. The repo label corresponds to the depth parameter, so a depth=2 as the example above would
have repo labels named `org1/repoa` and `org2/repob`.

| Metric                                   | Type  | Labels     | Description                              |
| ---------------------------------------- | ----- | ---------- | ---------------------------------------- |
| chartmuseum_charts_served_total          | Gauge | {repo="*"} | Total number of charts                   |
| chartmuseum_charts_versions_served_total | Gauge | {repo="*"} | Total number of chart versions available |

*: see above for repo label

There are other general global metrics harvested (per process, hence for all tenants). You can get the complete list by using the `/metrics` route.

| Metric                                     | Type    | Labels                                                | Description                               |
| ------------------------------------------ | ------- | ----------------------------------------------------- | ----------------------------------------- |
| chartmuseum_request_duration_seconds       | Summary | {quantile="0.5"}, {quantile="0.9"}, {quantile="0.99"} | The HTTP request latencies in seconds     |
| chartmuseum_request_duration_seconds_sum   |         |                                                       |                                           |
| chartmuseum_request_duration_seconds_count |         |                                                       |                                           |
| chartmuseum_request_size_bytes             | Summary | {quantile="0.5"}, {quantile="0.9"}, {quantile="0.99"} | The HTTP request sizes in bytes           |
| chartmuseum_request_size_bytes_sum         |         |                                                       |                                           |
| chartmuseum_request_size_bytes_count       |         |                                                       |                                           |
| chartmuseum_response_size_bytes            | Summary | {quantile="0.5"}, {quantile="0.9"}, {quantile="0.99"} | The HTTP response sizes in bytes          |
| chartmuseum_response_size_bytes_sum        |         |                                                       |                                           |
| chartmuseum_response_size_bytes_count      |         |                                                       |                                           |
| go_goroutines                              | Gauge   |                                                       | Number of goroutines that currently exist |


## Notes on index.yaml
The repository index (index.yaml) is dynamically generated based on packages found in storage. If you store your own version of index.yaml, it will be completely ignored.

`GET /index.yaml` occurs when you run `helm repo add chartmuseum http://localhost:8080` or `helm repo update`.

If you manually add/remove a .tgz package from storage, it will be immediately reflected in `GET /index.yaml`.

You are no longer required to maintain your own version of index.yaml using `helm repo index --merge`.

The `--gen-index` CLI option (described above) can be used to generate and print index.yaml to stdout.

Upon index regeneration, *ChartMuseum* will, however, save a statefile in storage called `index-cache.yaml` used for cache optimization. This file is only meant for internal use, but may be able to be used for migration to simple storage.

## Mirroring the official Kubernetes repositories
Please see `scripts/mirror_k8s_repos.sh` for an example of how to download all .tgz packages from the official Kubernetes repositories (both stable and incubator).

You can then use *ChartMuseum* to serve up an internal mirror:
```
scripts/mirror_k8s_repos.sh
chartmuseum --debug --port=8080 --storage="local" --storage-local-rootdir="./mirror"
 ```

## Original Logo

<sub>**_"Preserve your precious artifacts... in the cloud!"_**<sub>

![](./logo.png)

## Subprojects

The following subprojects are maintained by *ChartMuseum*:

- [chartmuseum/helm-push](https://github.com/chartmuseum/helm-push) - Helm plugin to push chart package to ChartMuseum
- [chartmuseum/storage](https://github.com/chartmuseum/storage) - Go library providing common interface for working across multiple cloud storage backends
- [chartmuseum/auth](https://github.com/chartmuseum/auth) - Go library for generating ChartMuseum JWT Tokens, authorizing HTTP requests, etc.
- [chartmuseum/auth-server-example](https://github.com/chartmuseum/auth-server-example) - Example server providing JWT tokens for ChartMuseum auth
- [chartmuseum/testbed](https://github.com/chartmuseum/testbed) - Docker testbed for continuous integration
- [chartmuseum/www](https://github.com/chartmuseum/www) - chartmuseum.com static site source code
- [chartmuseum/ui](https://github.com/chartmuseum/ui) - ChartMuseum frontend UI

## Community
You can reach the *ChartMuseum* community and developers in the [Kubernetes Slack](https://slack.k8s.io) **#chartmuseum** channel.
