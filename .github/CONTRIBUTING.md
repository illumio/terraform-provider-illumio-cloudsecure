# Contributing Guide

## GitHub Workflow

Non-Illumio contributors to the project should follow this workflow:

1. Fork the repo
2. Create a new branch on the fork
3. Push the branch to your fork
4. Submit a pull request following [GitHub's standard process](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests)

## Bug Reporting

> [!CAUTION]
> If you find a bug or issue that you believe to be a security vulnerability, please see the [SECURITY](SECURITY.md) document for instructions on reporting. **Do not file a public GitHub Issue.**

Please report any bugs you find as GitHub issues.

Before reporting any bugs, please do a quick search to see if it has already been reported. If so, please add a comment on the existing issue rather than creating a new one.

While reporting a bug, please provide a minimal example to reproduce the issue. Include `.tf` files, **making sure to remove any secrets**. If applicable, include the `crash.log` file as well.

## Testing

When submitting a new resource or datasource, please follow the current convention of including acceptance tests that set up, verify, and tear down the target resource. When making changes to existing resources or datasources, update the corresponding tests and add any unit tests you deem necessary to ensure the changes are working as expected and have not introduced regressions.

Refer to the Terraform [testing guideline](https://www.terraform.io/docs/extend/testing/index.html) for instructions on testing resources and datasources.

## Documentation

Documentation is an important aspect of the project. Changes to resources or datasources should be reflected in their respective docs files. Make sure to update the [CHANGELOG](../CHANGELOG.md) describing your changes. This project follows [HashiCorp's Changelog Specification](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#changelog-specification).

## Development

All development in this repository should be done in a [devcontainer](https://containers.dev/) as defined in `/.devcontainer/`, using a compatible IDE, to ensure that you use consistent versions of the development tools.
Recommended IDEs include [GitHub Codespaces](https://github.com/features/codespaces) and [VS Code](https://code.visualstudio.com/).

Run the following command to build the provider:

```
go build -o terraform-provider-illumio-cloudsecure
```

To test this provider binary with Terraform locally, copy the provider binary into the local plugin registry cache directory:

```
export ARCH=`dpkg --print-architecture`
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/illumio/illumio-cloudsecure/0.0.1/linux_${ARCH}
cp terraform-provider-illumio-cloudsecure ~/.terraform.d/plugins/registry.terraform.io/illumio/illumio-cloudsecure/0.0.1/linux_${ARCH}
```

## Release Checklist

1. Implement the `ConfigService` gRPC server as specified in the commit to be released, and deploy that server.
1. Choose a version number for the new release. Follow [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) format `vX.Y.Z`, e.g. `v1.2.3`.
1. Update the version number in `examples/provider/provider.tf` (omit the `v` prefix), and re-generate the docs (`make generate`):

   ```
   terraform {
      required_providers {
         illumio-cloudsecure = {
            source  = "illumio/illumio-cloudsecure"
            version = "~> 1.5.0"
         }
      }
   }
   ```

1. Create a [new GitHub release](https://github.com/illumio/terraform-provider-illumio-cloudsecure/releases) from that commit. Summarize the changes in this release.

   1. Click `Generate release notes` to generate the `What's Changed` section.
   1. Edit the generated `What's Changed` section to follow the [Changelog Specification](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#changelog-specification), esp. to follow the same [categories](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning#categorization) for subsections. Link to the PR(s) for each item. For example:

      ```
      ## What's Changed

      ### BREAKING CHANGES:

      * provider: Bump API version to v2 #12

      ### FEATURES:

      * **New Resource:** `azure_account` #14

      ### BUG FIXES:

      * provider: Fix handling of optional attributed #24
      ```

1. Verify that the `release` workflow completes for this release tag (see [Actions](https://github.com/illumio/terraform-provider-illumio-cloudsecure/actions)) and that it generates and publishes the artifacts (binaries for all systems, etc.) into [this release](https://github.com/illumio/terraform-provider-illumio-cloudsecure/releases).
1. Verify that the [HashiCorp Terraform Registry](https://registry.terraform.io/providers/illumio/illumio-cloudsecure/latest) publishes the new release.
