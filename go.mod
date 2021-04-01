module github.com/vshn/swisscom-service-broker

go 1.15

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/crossplane/crossplane-runtime v0.13.0
	github.com/gorilla/mux v1.8.0
	github.com/pivotal-cf/brokerapi/v7 v7.5.0
	github.com/stretchr/testify v1.7.0
	github.com/vshn/crossplane-service-broker v0.2.0
	k8s.io/api v0.20.1
	k8s.io/client-go v0.20.1
	sigs.k8s.io/controller-runtime v0.8.0
)
