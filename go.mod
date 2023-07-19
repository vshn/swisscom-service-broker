module github.com/vshn/swisscom-service-broker

go 1.16

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/crossplane/crossplane-runtime v0.16.1
	github.com/gorilla/mux v1.8.0
	github.com/pivotal-cf/brokerapi/v7 v7.5.0
	github.com/stretchr/testify v1.7.2
	github.com/vshn/crossplane-service-broker v0.10.0
	k8s.io/api v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/controller-runtime v0.11.0
	sigs.k8s.io/kustomize/kustomize/v3 v3.10.0
)
