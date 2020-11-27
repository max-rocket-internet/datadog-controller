# datadog-controller: A kubernetes controller managing Datadog Monitors

[![](img/banner.png)](#)

This is a simple [Kubernetes Controller](https://kubernetes.io/docs/concepts/architecture/controller/) to allow [Datadog Monitors](https://docs.datadoghq.com/monitors/) to be created, updated or deleted from custom resources in Kubernetes.

Here's an example `DatadogMonitor` resource:

```yaml
apiVersion: datadoghq.com/v1beta1
kind: DatadogMonitor
metadata:
  name: apm-error-rate-example
spec:
  name: my-service error rate
  query: "avg(last_5m):sum:trace.http.request.errors{env:stg,service:my-service} / sum:trace.http.request.hits{env:stg,service:my-service} > 1"
  type: "query alert"
  message: Service my-service has a high error rate on env:stg
```

## Installation

You will need a Datadog APP and API key which can be found or created at [app.datadoghq.eu/account/settings](https://app.datadoghq.eu/account/settings#api) or [app.datadoghq.com/account/settings](https://app.datadoghq.com/account/settings#api).

Use the included [Helm](https://helm.sh/) chart in [chart](chart) or install the chart from the [Delivery Hero Helm charts repo](https://github.com/deliveryhero/helm-charts):

```console
helm repo add deliveryhero https://charts.deliveryhero.io/
helm search repo deliveryhero
helm install datadog-controller deliveryhero/datadog-controller --set datadog.client_api_key="YOUR_API_KEY" --set datadog.client_app_key="YOUR_APP_KEY"
```

Or a docker image is available at [maxrocketinternet/datadog-controller](https://hub.docker.com/r/maxrocketinternet/datadog-controller).

## Examples

There are more examples in the [examples](examples) directory.

## Test or run locally

Set your `kubectl` context as required and export required environment variables:

```
export DD_CLIENT_API_KEY="YOUR_API_KEY"
export DD_CLIENT_APP_KEY="YOUR_APP_KEY"
```

Then run `main.go`:

```
go run main.go
```

To run tests you need to [install kubebuilder](https://book.kubebuilder.io/quick-start.html#installation) which includes the required `kube-apiserver` and `etcd` to test the controller:

```
go test ./...
```

## Notes

The contents of this repository are licensed under the Apache License version 2.0.

This project was created using [kubebuilder](https://book-v1.book.kubebuilder.io/).

This repository and its authors are in no way associated with the company Datadog.
