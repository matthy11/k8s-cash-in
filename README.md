# heypay-cash-in-server

This server handles cash-in logic for received and reversed transfers. The implemented endpoints and their documentation can be found [here](ENDPOINTS_README.md).

## Table of Contents
1. [Setup](#setup)
2. [Configuration](#configuration)
3. [Run](#run)
4. [Deploy](#deploy)
5. [Testing](#testing)
7. [Release](#release)

<a name="setup"></a>
## Setup

First you'll need to set up Go development environment, currently this project uses:

- Go v1.15.X

Check this [guide](https://golang.org/doc/install) to install Golang. 

### Generate gRPC Proto files

This server calls services through gRPC. To add and/or update new services you will need to run the following commands to generate `<proto-entity>.pb.go` files associated with the service:

```sh
$ cd internal/services/<proto-entity>service/grpc
$ protoc --proto_path=. --go_out=plugins=grpc:. ./<proto-entity>.proto
```

<a name="configuration"></a>
## Configuration

Check `config/env.json` file to view available environment variables.

Also, you will need two service accounts:

- `k8s/service-account-dev.json`: App project Firebase default service account.
- `k8s/grpc-client-accounts-engine-service-account.json`: Accounts engine service account with `Cloud Run Invoker` role.

<a name="run"></a>
## Run

To run the project locally, use the following command on project root folder:

```sh
$ go run .

# You should see the following logs indicating that the server is running
INFO Initializing firebase service...              filename=firebase_service.go function=init.0 ip=undefined line=33 loggerId=undefined severity=info step=undefined version=1.2.0
INFO Initializing identity toolkit service...      filename=identity_toolkit_service.go function=init.0 ip=undefined line=15 loggerId=undefined severity=info step=undefined version=1.2.0
INFO Initializing accounts engine service...       filename=account_engine_service.go function=init.0 ip=undefined line=25 loggerId=undefined severity=info step=undefined version=1.2.0
INFO Starting the service...                       filename=main.go function=main ip=undefined line=11 loggerId=undefined severity=info step=undefined version=1.2.0
INFO The service is ready to listen and serve.     filename=main.go function=main ip=undefined line=13 loggerId=undefined severity=info step=undefined version=1.2.0
```

<a name="deploy"></a>
## Deploy

This server is configured to run as containerized app on Google Kubernetes Engine. Since this project it's configured to work with Google Cloud Build, when pushing to a specific branch it will trigger the deployment.

*Note: When deploying to production environment, Cloud Build trigger it's configured to listen for pushes to `master` branch, but it's disabled for security reasons, so it must be run manually.*

### Connect to GKE cluster

Before working with kubectl, make sure you are connected to the desired cluster by running the following command:

```sh
$ gcloud container clusters get-credentials <cluster_name> --zone <cluster_zone> --project <gcp_project_name>
```

*Note: You can get this command ready-to-use by going to Google Cloud Platform project, then to Kubernetes > Clusters dashboard, click on the desired cluster to use and press "Connect" button.*

### Configure Kubernetes

```sh
$ kubectl apply -f ./k8s/deploy.yaml
$ kubectl apply -f ./k8s/backend-config.yaml
$ kubectl apply -f ./k8s/cash-in-service.yaml
$ kubectl apply -f ./k8s/ingress.yaml
```

### Server Configuration and Secrets

```sh
# Create cash-in-config secret with values from env.json file containing environment variables
$ kubectl create secret generic cash-in-config --from-literal=config="$(cat config/env.json)" --namespace="cash-in"

# Get cash-in-config secret
$ kubectl get secret cash-in-config -o jsonpath="{.data.config}" --namespace="cash-in" | base64 --decode

# Create gcloud-cash-in-key secret from service-account.json file containing service account
$ kubectl create secret generic gcloud-cash-in-key --from-file=k8s/service-account.json --namespace="cash-in"

# Get gcloud-cash-in-key secret
$ kubectl get secret gcloud-cash-in-key -o jsonpath="{.data['service-account\.json']}" --namespace="cash-in" | base64 --decode

# Create gcloud-grpc-client-accounts-engine-key secret from grpc-client-accounts-engine-service-account.json file containing service account
$ kubectl create secret generic gcloud-grpc-client-accounts-engine-key --from-file=k8s/grpc-client-accounts-engine-service-account.json --namespace="cash-in"

# Get gcloud-grpc-client-accounts-engine-key secret
$ kubectl get secret gcloud-grpc-client-accounts-engine-key -o jsonpath="{.data['grpc-client-accounts-engine-service-account\.json']}" --namespace="cash-in" | base64 --decode
```

<a name="testing"></a>
## Testing

To execute tests, run the following commands:
```sh
# In order to apply execution permissions, should be ran once
$ chmod a+x tests/test.sh
$ chmod a+x tests/test_debug.sh

# Run one of the following
$ ./tests/tests.sh
$ ./tests/tests_debug.sh
```

<a name="release"></a>
## Release

To generate a new release, follow this steps:

- From `develop` branch, fetch the latest changes
- Change `V.Version` value in file `settings/settings.go`
- `git add .`
- `git commit -m "Version bump —version <version>"`
- Update CHANGELOG.md with output of `git log --pretty=format:"%H - %s" --no-merges <last-version>..`
- `git add .`
- `git commit -m "CHANGELOG updated <version>"`
- `git push`
- From `staging` branch, fetch the latest changes and merge `develop` branch
- Draft new release
- From `master` branch, fetch the latest changes and merge `staging` branch
