#!/bin/bash

set -e

install-kafka() {
  if ! kubectl get ns kafka > /dev/null 2>&1; then
    echo "Namespace kafka does not exist. Create namespace"
    kubectl create ns kafka
    echo "Successfully created namespace. Install kafka"
    helm upgrade --install my-kafka oci://registry-1.docker.io/bitnamicharts/kafka -f ./values.yaml -n kafka
  else 
    echo "Name space kafka already exists. Install kafka"
    helm upgrade --install my-kafka oci://registry-1.docker.io/bitnamicharts/kafka -f ./values.yaml -n kafka
  fi
}

main() {
  local answer
  local context=$(kubectl config current-context)
  read -p "You sure want to proceed with $context ? ([y|Y] others will exit) " answer

  case "$answer" in
    [yY]) install-kafka ;;
    *) echo "Shutting down program ..."; return 1 ;;
  esac
}

main
