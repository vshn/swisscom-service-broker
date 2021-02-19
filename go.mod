module github.com/vshn/swisscom-service-broker

go 1.15

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/crossplane/crossplane-runtime v0.12.0
	github.com/gorilla/mux v1.8.0
	github.com/pivotal-cf/brokerapi/v7 v7.5.0
	github.com/stretchr/testify v1.7.0
	github.com/vshn/crossplane-service-broker v0.1.0-rc1
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.20.0
	k8s.io/client-go v0.19.3
	sigs.k8s.io/controller-runtime v0.6.4
)
