<p align="center">

<img  src="https://mk0abtastybwtpirqi5t.kinstacdn.com/wp-content/uploads/picture-solutions-persona-product-flagship.jpg"  width="211"  height="182"  alt="flagship"  />

</p>

<h3 align="center">Bring your features to life</h3>

A Tool to manage your AB Tasty resources built in [Go](https://go.dev/) using the library [Cobra](https://cobra.dev/).

[Website](https://flagship.io) | [Website Documentation](https://docs.developers.flagship.io/docs/flagship-command-line-interface) | [Command Documentation](https://flagship-io.github.io/abtasty-cli/documentation/abtasty-cli) | [Installation Guide](https://docs.developers.flagship.io/docs/flagship-command-line-interface#download-and-install-the-flagship-cli) | [Twitter](https://twitter.com/feature_flags)

[![Apache2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/flagship-io/abtasty-cli)](https://goreportcard.com/report/github.com/flagship-io/abtasty-cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/flagship-io/abtasty-cli.svg)](https://pkg.go.dev/github.com/flagship-io/abtasty-cli)
[![test](https://github.com/flagship-io/abtasty-cli/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/flagship-io/abtasty-cli/actions/workflows/ci.yml)
[![coverage](https://raw.githubusercontent.com/flagship-io/abtasty-cli/badges/.badges/main/coverage.svg)](https://raw.githubusercontent.com/flagship-io/abtasty-cli/badges/.badges/main/coverage.svg)

## Overview

The AB Tasty CLI is a set of commands to create and manage your AB Tasty resources for feature experimentation and web experimentation products such as projects, campaigns, teams, etc... You can use these commands to perform common AB Tasty platform actions from your terminal or through scripts and other automation.

Our CLI is built on top of our Remote Control API, enabling you to use direct API calls or use the CLI.

For example, you can use the AB Tasty CLI to manage :

- Feature experimentation:
  - Projects, campaigns, flags, targeting keys, goals, etc...
  - Users and environments
  - Panic mode
- Web experimentation:
  - campaigns, variations, accounts, elementJS, etc...
  - Pull & Push global codes (account, campaign, variation, elementJS)

## The AB Tasty CLI cheat sheet

For an introduction to the AB Tasty CLI, a list of commonly used commands, and a look at how these commands are structured, see the [AB Tasty cheat sheet](https://docs.developers.flagship.io/docs/ab-tasty-cli-reference-v1xx).

## Download and install AB Tasty CLI

The AB Tasty CLI can be installed and deployed in your infrastructure either by downloading and running the binary, or pulling and running the docker image in your orchestration system.

### Using a binary

- Linux/Darwin

#### With wget

```bash
wget -qO- https://raw.githubusercontent.com/flagship-io/abtasty-cli/main/install.sh | bash
```

#### With curl

```bash
curl -sL https://raw.githubusercontent.com/flagship-io/abtasty-cli/main/install.sh | bash
```

#### With Homebrew

```bash
brew tap flagship-io/abtasty-cli
brew install abtasty-cli
```

- Other Operating Systems

Please download the binary from the [release page](https://github.com/flagship-io/abtasty-cli/releases)

### Using Golang from source

You can pull the project from github and build it using golang latest stable version (+1.18):

    git clone git@github.com:flagship-io/abtasty-cli.git
    cd abtasty-cli
    go build -o abtasty-cli

## Contributors

- Chadi Laoulaou [@Chadiii](https://github.com/chadiii)
- Guillaume Jacquart [@GuillaumeJacquart](https://github.com/guillaumejacquart)

## Contributing

Please read our [contributing guide](./CONTRIBUTING.md).

## License

[Apache License.](https://github.com/flagship-io/abtasty-cli/blob/main/LICENSE)

## About Feature Experimentation

​
<img src="https://www.flagship.io/wp-content/uploads/Flagship-horizontal-black-wake-AB.png" alt="drawing" width="150"/>
​
[Flagship by AB Tasty](https://www.flagship.io/) is a feature flagging platform for modern engineering and product teams. It eliminates the risks of future releases by separating code deployments from these releases :bulb: With Flagship, you have full control over the release process. You can:
​

- Switch features on or off through remote config.
- Automatically roll-out your features gradually to monitor performance and gather feedback from your most relevant users.
- Roll back any feature should any issues arise while testing in production.
- Segment users by granting access to a feature based on certain user attributes.
- Carry out A/B tests by easily assigning feature variations to groups of users.
  ​
  <img src="https://www.flagship.io/wp-content/uploads/demo-setup.png" alt="drawing" width="600"/>
  ​
  Flagship also allows you to choose whatever implementation method works for you from our many available SDKs or directly through a REST API. Additionally, our architecture is based on multi-cloud providers that offer high performance and highly-scalable managed services.
  ​
  **To learn more:**
  ​
- [Solution overview](https://www.flagship.io/#showvideo) - A 5mn video demo :movie_camera:
- [Website Documentation](https://docs.developers.flagship.io/) - Our dev portal with guides, how tos, API and SDK references
- [Command Documentation](https://flagship-io.github.io/abtasty-cli/documentation/abtasty-cli) - Command references
- [Sign up for a free trial](https://www.flagship.io/sign-up/) - Create your free account
- [Guide to feature flagging](https://www.flagship.io/feature-flags/) - Everything you need to know about feature flag related use cases
- [Blog](https://www.flagship.io/blog/) - Additional resources about release management
