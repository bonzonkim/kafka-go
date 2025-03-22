# Prerequisite
* helm
* kubectl
* kubernetes(orbstack, minikube, your own kubernetes cluster or whatever)

# Custom Values
Modify the values in `values.yaml` as needed.

The current configuration exposes the ports using NodePort:
* client : 30092
* external : 30095

To change release name, modify `helm install` command in `intsall-kafka.sh`

# Run
```bash
# give permission
chmod +x install-kafka.sh

# run script
./install-kafka.sh
```

This script will install `kafka` using [helm chart](https://artifacthub.io/packages/helm/bitnami/kafka) in namespace kafka.
