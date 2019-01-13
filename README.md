# Micro Enterprise [![License](https://img.shields.io/badge/license-enterprise-blue.svg)](https://github.com/micro/enterprise/blob/master/LICENSE)


Micro Enterprise is a turn-key solution for building production ready microservices.

## Overview

Micro Enterprise is an enterprise version of the [micro](https://github.com/micro/micro) toolkit. It builds on the pluggable open source 
toolkit and pre-packages the most useful plugins along with additiona features in a tried and tested single solution critical ready for production use.

<img src="https://micro.mu/micro-enterprise.png" />

This version of Micro requires a [Micro Enterprise License Agreement](LICENSE) commercial subscription.

## Features

Micro Enterprise builds on the all the features of the [Micro Toolkit](https://github.com/micro/micro) along with the following:

- **Dynamic Config** - Load config from environment variables, flags and a config file, all integrated into one interface. Config is 
merged, watched and reloaded as it changes. 

- **Plugin Loading** - Plugins can be built via the command line or built and loaded on the fly. Build your apps and the micro toolkit 
once, swap out plugins at runtime. This enables a flexible and portable runtime.

- **Authentication** - Support for basic, digest, ldap and other forms of authentication. Quickly enable auth on any component of 
the toolkit. Limit the access to the web dashboard or the api gateway easily.

- **ChatOps Inputs** - The micro bot provides ChatOps as a first class citizen. The bot lives within your platform and allows you to 
manage applications via messaging. This includes support for Discord, HipChat, Slack and Telegram.

- **HTTP Proxying** - Micro is an RPC based system. It's most likely you have a multi-protocol architecture and one that heavily 
relies on HTTP. We provide a simple RPC to HTTP service for proxying to http backends. Leverage the micro ecosystem for any language.

- secure by default: tls enabled
- authentication: rbac and service-to-service
- central control plane
- circuit breaking
- distributed tracing
- rate limiting
- smart routing
- instrumentation
- performance tuned

## Pricing

See the [website](https://micro.mu/pricing) for details

## Getting Started

See the [docs](https://micro.mu/docs/enterprise.html) to get started

## License

This version of Micro is distributed under the commercial [Micro Enterprise License Agreement](LICENSE)
