ARG VARIANT=1.22-bookworm
FROM mcr.microsoft.com/vscode/devcontainers/go:${VARIANT}

ARG BUF_VERSION=v1.35.0
ARG GOLANGCI_VERSION=v1.59.0
ARG TERRAFORM_VERSION=1.7.5-1

USER root

# Install basic tools
RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get -y install --no-install-recommends \
    bash-completion \
    git \
    gnupg \
    software-properties-common \
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
    && rm -rf /var/lib/apt/lists/*

USER vscode

# Install Terraform shell autocomplete
RUN terraform -install-autocomplete

RUN go install github.com/bufbuild/buf/cmd/buf@${BUF_VERSION} \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_VERSION}

USER root

RUN /go/bin/buf completion bash > "/etc/bash_completion.d/buf"

USER vscode
