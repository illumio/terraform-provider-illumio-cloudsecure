# Copyright (c) Illumio, Inc.
# SPDX-License-Identifier: MPL-2.0

# Docker image that contains Terraform and the Terraform Provider for Illumio CloudSecure.
#
# To build the image:
# docker build -t illumio-terraform:1.0.0-dev ./
#
# To build the image for another platform, e.g. arm64:
# docker buildx build --platform linux/arm64 -t illumio-terraform:1.0.0-dev ./
#
# To run the image:
# docker run --rm -ti -e TF_VAR_illumio_cloudsecure_client_id -e TF_VAR_illumio_cloudsecure_client_secret --mount type=bind,src="$(pwd)",target=/workspace illumio-terraform:1.0.0-dev <terraform command and arguments...>

FROM golang:1.24.3-bookworm AS build

ARG VERSION=1.0.0-dev

USER root

# Build the terraform-provider-illumio-cloudsecure provider and copy it to /

RUN mkdir -p /build \
    && cd /build \
    && git clone https://github.com/illumio/terraform-provider-illumio-cloudsecure.git \
    && cd terraform-provider-illumio-cloudsecure \
    && git checkout v${VERSION} \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=v${VERSION} -X main.commit=`git rev-parse --short HEAD`" -a -o terraform-provider-illumio-cloudsecure \
    && mv terraform-provider-illumio-cloudsecure / \
    && cd / \
    && rm -rf /build

FROM debian:bookworm

ARG TARGETARCH
ARG TERRAFORM_VERSION=1.11.4-1
ARG VERSION=1.0.0-dev

USER root

# Install basic tools

RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get -y install --no-install-recommends \
    bash-completion \
    binutils \
    gnupg \
    software-properties-common \
    wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Terraform

RUN wget -q -O- https://apt.releases.hashicorp.com/gpg | \
    gpg --dearmor > /usr/share/keyrings/hashicorp-archive-keyring.gpg \
    && gpg --no-default-keyring \
    --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg \
    --fingerprint \
    && echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" > /etc/apt/sources.list.d/hashicorp.list \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get -y install --no-install-recommends \
    terraform=${TERRAFORM_VERSION} \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && terraform -install-autocomplete

# Create a terraform:terraform user

RUN addgroup --gid 10000 terraform && adduser --uid 10000 --gid 10000 --disabled-password --home /terraform terraform \
    && mkdir -p /terraform/.terraform.d/plugins/registry.terraform.io/illumio/illumio-cloudsecure/${VERSION}/linux_${TARGETARCH} \
    && chown -R terraform:terraform /terraform

# Pre-install the terraform-provider-illumio-cloudsecure provider

COPY --from=build /terraform-provider-illumio-cloudsecure /terraform/.terraform.d/plugins/registry.terraform.io/illumio/illumio-cloudsecure/${VERSION}/linux_${TARGETARCH}/

# Create a /workspace volume mount point and use that as the working directory

RUN mkdir -p /workspace \
    && chown -R terraform:terraform /workspace

WORKDIR /workspace

USER terraform

ENTRYPOINT ["/usr/bin/terraform"]
