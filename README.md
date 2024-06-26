# Terraform Registry Mirror

![Service Logo](https://github.com/jonasasx/terraform-registry-mirror/blob/master/assets/logo.png?raw=true)

## Source Code

https://github.com/jonasasx/terraform-registry-mirror

## Overview

Welcome to the Terraform Registry Mirror Service repository!
This service implements the [Provider Network Mirror Protocol](https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol),
enabling you to mirror repositories of Terraform providers across
your network infrastructure. With Terraform Registry Mirror Service,
you can efficiently manage and distribute Terraform provider
repositories within your organization's network.

## Features

- **Repository Mirroring**: Mirror Terraform provider repositories to your local network.

## Prerequisites

Before using the Terraform Registry Mirror Service, ensure you have the following:

- Network infrastructure capable of hosting the Terraform Registry Mirror Service
- Access to the internet to fetch Terraform provider repositories initially
- Knowledge of networking concepts and basic server administration

## Usage

Once installed, you can start using the Terraform Registry Mirror Service to mirror Terraform provider repositories.

Create `.terraformrc` file in the home directory:

    provider_installation {
        network_mirror {
            url = "https://terraform-registry-mirror.ru/"
        }
    }

You can use any terraform cli tool commands such as `terraform init`.


## Installation

Choose one of the installation methods:

### Docker

    docker run -d -p 8080:8080 ghcr.io/jonasasx/terraform-registry-mirror:0.0.8

### Docker Compose

    version: '3'
    
    services:
      terraform-registry-mirror:
        image: ghcr.io/jonasasx/terraform-registry-mirror:0.0.8
        ports:
          - "8080:8080"

### Helm install

    helm install <RELEASE_NAME> oci://ghcr.io/jonasasx/terraform-registry-mirror/terraform-registry-mirror --version 0.1.0

### ArgoCD

    apiVersion: argoproj.io/v1alpha1
    kind: Application
    metadata:
      name: terraform-registry-mirror
    spec:
      destination:
        namespace: default
        server: https://kubernetes.default.svc
      project: default
      source:
        chart: terraform-registry-mirror
        repoURL: ghcr.io/jonasasx/terraform-registry-mirror
        targetRevision: 0.1.0

## Contributing

Contributions to the Terraform Registry Mirror Service are welcome! If you encounter any issues, have suggestions for improvements, or would like to contribute code, feel free to open an issue or pull request in the GitHub repository.

## License

The Terraform Registry Mirror Service is licensed under the [MIT License](./LICENSE), allowing for both personal and commercial use.

---

Thank you for choosing the Terraform Registry Mirror Service! If you have any questions or need assistance, don't hesitate to reach out.