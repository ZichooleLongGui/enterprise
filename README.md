# Micro Enterprise [![License](https://img.shields.io/badge/license-enterprise-blue.svg)](https://github.com/micro/enterprise/blob/master/LICENSE)


Micro Enterprise is a turn-key solution for building production ready microservices.

## Overview

Micro Enterprise is an enterprise version of the [micro](https://github.com/micro/micro) toolkit. It builds on the pluggable open source 
toolkit and pre-packages the most useful plugins along with additiona features in a tried and tested single solution critical ready for production use.

<img src="https://micro.mu/micro-enterprise.png" />

This version of Micro requires a [Micro Enterprise License Agreement](LICENSE) commercial subscription.

## Features

Micro Enterprise builds on the all the features of the [Micro Toolkit](https://github.com/micro/micro) along with the following:

- **Zero Dependency** - Simplified setup with zero external dependencies. There's no need for external service discovery or storage. 
We handle everything internally. Just drop in the api or proxy and get started straight away. 

- **ACME Certificates** - The API Gateway supports ACME TLS certificate management via Let's Encrypt. Simply enable support on the command 
line and run on a secure port. No other configuration needed.

- **Dynamic Config** - Load config from environment variables, flags and a config file, all integrated into one interface. Config is 
merged, watched and reloaded as it changes. 

- **Plugin Loading** - Plugins can be built via the command line or built and loaded on the fly. Build your apps and the micro toolkit 
once, swap out plugins at runtime. This enables a flexible and portable runtime.

- **Authentication** - Support for basic, digest, ldap and other forms of authentication. Quickly enable auth on any component of 
the toolkit. Limit the access to the web dashboard or the api gateway easily.

- **ChatOps Inputs** - The micro bot provides ChatOps as a first class citizen. The bot lives within your platform and allows you to 
manage applications via messaging. This includes support for Discord, HipChat, Slack and Telegram.

- **HTTP Bridge** - Micro is an RPC based system. It's most likely you have a multi-protocol architecture and one that heavily 
relies on HTTP. We provide a simple RPC to HTTP service for proxying to http backends. Leverage the micro ecosystem for any language.

- **CORS Support** - The API Gateway, Web Dashboard and Service Proxy all support the addition of CORS control. This allows you to 
define how Cross-Origin Resource Sharing is dealt with from one place. 

## Roadmap

Features to be integrated:

- secure by default: spifee based x509 identities and mutual tls
- authentication: rbac service to service access control
- central control plane: single location to manage acls, routing, etc
- circuit breaking: fail fast when errors occur
- rate limiting: limit thundering herd issues when things fail
- smart routing: weighted and priority based routing
- built in metrics: record and retrieve stats/debug info 
- distributed tracing: understand the behaviour of requests
- distributed logging: see what happened and when it happened
- performance tuned: optimised from day 1 for high performance
- federated routing: multi-dc networking with minimal config

## Pricing

See the [website](https://micro.mu/pricing) for details

## Getting Started

See the [docs](https://micro.mu/docs/enterprise.html) to get started

## License

This version of Micro is distributed under the commercial [Micro Enterprise License Agreement](LICENSE)
