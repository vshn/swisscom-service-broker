module github.com/vshn/swisscom-service-broker

go 1.16

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/coreos/etcd v3.3.15+incompatible // indirect
	github.com/crossplane/crossplane-runtime v0.15.1-0.20211029211307-c72bcdd922eb
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-getter v1.4.0 // indirect
	github.com/pivotal-cf/brokerapi/v7 v7.5.0
	github.com/stretchr/testify v1.7.0
	github.com/vshn/crossplane-service-broker v0.7.0
	go.uber.org/tools v0.0.0-20190618225709-2cfd321de3ee // indirect
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	k8s.io/api v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/klog v0.4.0 // indirect
	sigs.k8s.io/controller-runtime v0.9.6
	sigs.k8s.io/kustomize/kustomize/v3 v3.10.0
	sigs.k8s.io/structured-merge-diff v0.0.0-20190817042607-6149e4549fca // indirect
)
