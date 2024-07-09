# Contributing

We're excited for your contributions to the illumio cloudsecure terraform provider!
This document outlines how you can get involved.

## Dev Setup
### Prerequisites
vscode, docker

### Instructions
1. Clone the repository
2. Open the repository in vscode
3. Start the project in dev container
   1. vscode will prompt you to reopen the project in a dev container use it.
   2. If not, you can open the command palette and run `Remote-Containers: Reopen in Container`
3. Run `make clean` to clean the project
4. Run `make generate` to generate the code from generators
5. Run `make run` to run the tests
6. Run `make build` to build the project and save to the provider to the local terraform registry
