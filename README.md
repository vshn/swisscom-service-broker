[![Build](https://img.shields.io/github/workflow/status/vshn/swisscom-service-broker/Pull%20Request)][build]
![Go version](https://img.shields.io/github/go-mod/go-version/vshn/swisscom-service-broker)
[![Version](https://img.shields.io/github/v/release/vshn/swisscom-service-broker)][releases]
[![GitHub downloads](https://img.shields.io/github/downloads/vshn/swisscom-service-broker/total)][releases]
[![Docker image](https://img.shields.io/docker/pulls/vshn/swisscom-service-broker)][dockerhub]
[![License](https://img.shields.io/github/license/vshn/swisscom-service-broker)][license]

# Swisscom Service Broker

[Open Service Broker](https://github.com/openservicebrokerapi/servicebroker) API which provisions
Redis and MariaDB instances via [crossplane](https://crossplane.io/).

Based on [crossplane-service-broker](https://github.com/vshn/crossplane-service-broker).

## Documentation

Most of the explanation on how this all works together currently lives in the [VSHN Knowledgebase](https://kb.vshn.ch/app-catalog/explanations/crossplane_service_broker.html).

## Contributing

You'll need:

- A running kubernetes cluster (minishift, minikube, k3s, ... you name it) with crossplane installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and [kustomize](https://kubernetes-sigs.github.io/kustomize/installation/)
- Go development environment
- Your favorite IDE (with a Go plugin)
- docker
- make

These are the most common make targets: `build`, `test`, `docker-build`, `run`.

### Folder structure overview

```
.
├── cmd
│   └── swisscom-service-broker      # main file
├── deploy
│   └── base                         # deployment files
├── docs                             # antora docs
├── e2e                              # e2e testing files
├── pkg
│   ├── custom                       # custom API with swisscom specifics
└── testdata                         # integration testing files
```

### Run the service broker

You can run the operator in different ways:

1. using `make run` (provide your own env variables)
1. using `make kind-run` (uses KIND to install a cluster in docker and provides its own kubeconfig in `testbin/`)
1. using a configuration of your favorite IDE (see below for VSCode example)

Example VSCode run configuration:

```
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/swisscom-service-broker/main.go",
      "env": {
        "KUBECONFIG": "path/to/kubeconfig",
        "OSB_USERNAME": "test",
        "OSB_PASSWORD": "TEST",
        "OSB_SERVICE_IDS": "PROVIDE-SERVICE-UUIDS-HERE",
        "OSB_NAMESPACE": "test"
      },
      "args": []
    }
  ]
}
```

## Run integration tests

"Integration" testing is done using [envtest](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest) and [crossplane's integration test helper](https://github.com/crossplane/crossplane-runtime/tree/master/pkg/test/integration).

```
make integration-test
```

## Run E2E tests

The e2e tests currently only test if the deployment works. They do not represent a real
e2e test as of now but are meant as a base to build upon.

You need `node` and `npm` to run the tests, as it runs with [DETIK][detik].

To run e2e tests for newer K8s versions run

```
make e2e-test
```

To remove the local KIND cluster and other resources, run

```
make clean
```

[build]: https://github.com/vshn/swisscom-service-broker/actions?query=workflow%3APull%20Request
[releases]: https://github.com/vshn/swisscom-service-broker/releases
[license]: https://github.com/vshn/swisscom-service-broker/blob/master/LICENSE
[dockerhub]: https://hub.docker.com/r/vshn/swisscom-service-broker
[detik]: https://github.com/bats-core/bats-detik
