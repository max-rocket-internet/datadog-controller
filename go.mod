module github.com/max-rocket-internet/datadog-controller

go 1.13

require (
	github.com/DataDog/datadog-api-client-go v1.0.0-beta.10
	github.com/go-logr/logr v0.1.0
	github.com/hashicorp/go-retryablehttp v0.6.7
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/sirupsen/logrus v1.4.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)
