# apigeecli

[![TravisCI](https://travis-ci.org/srinandan/apigeecli.svg?branch=master)](https://travis-ci.org/srinandan/apigeecli)
[![Go Report Card](https://goreportcard.com/badge/github.com/srinandan/apigeecli)](https://goreportcard.com/report/github.com/srinandan/apigeecli)
[![GitHub release](https://img.shields.io/github/v/release/srinandan/apigeecli)](https://github.com/srinandan/apigeecli/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is a tool to interact with [Apigee APIs](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest) for [Apigee hybrid](https://cloud.google.com/apigee/docs/hybrid/v1.3/what-is-hybrid) and [Apigee's managed](https://cloud.google.com/apigee/docs/api-platform/get-started/overview) offering. The tool lets you manage (Create,Get, List, Update, Delete, Export and Import) Apigee entities like proxies, products etc. The tools also helps you create Service Accounts in Google IAM to operate Apigee hybrid runtime.

## Installation

### Brew

To install via brew,

```bash
brew tap srinandan/homebrew-tap
brew install apigeecli
```

### Others

`apigeecli` is a binary and you can download the appropriate one for your platform from [here](https://github.com/srinandan/apigeecli/releases)

NOTE: Supported platforms are:

* Darwin
* Windows
* Linux

## What you need to know about apigeecli

You must have an account on [Apigee](https://cloud.google.com/apigee/docs) to perform any `apigeecli` functions. These functions include: proxies, API Products, Environments, Org details etc.

You need to be familiar with basic concepts and features of Apigee such as API proxies, organizations, and environments.

For more information, refer to the [Apigee API Reference](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest).

## Service Account

Create a service account with appropriate persmissions. Use `apigeecli` to create service accounts (`apigeectl iam`). Read more [here](https://cloud.google.com/apigee/docs/hybrid/v1.3/sa-about)

## Access Token

`apigeecli` can use the service account directly and obtain an access token.

```bash
apigeecli token gen -a serviceaccount.json 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format


Use this access token for all subsequent calls (token expires in 1 hour)

## Command Reference

The following options are available for security

Pass the access token

```bash
apigeecli <flags> -t $TOKEN
```

Pass the service account

```bash
apigeecli <flags> -a orgadmin.json
```

## Access Token Caching

`apigeecli` caches the OAuth Access token for subsequent calls (until the token expires). The access token is stored in `$HOME/.apigeecli`. This path must be readable/writeable by the `apigeecli` process.

```bash
apigeecli token cache -a serviceaccount.json
```

or

```bash
apigeecli orgs get -o org-name -a serviceaccount.json
```

Subsequent commands do not need the token or service account flag

## Preferences

Users can set a default org via preferences and that org name will be used for all subsequent commands

```bash
apigeecli prefs set -o org-name

apigeecli orgs get
```

NOTE: the second command uses the org name from perferences

## Apigee Client Library

apigeecli is can also be used as a golang based client library. Look at this [sample](./samples) for more details

## Docker

Use apigecli via docker

```bash
docker run --name apigeecli -v path-to-service-account.json:/etc/client_secret.json --rm nandanks/apigeecli:v{Tag} orgs list -a /etc/client_secret.json
```

___

## Available Commands

Here is a [list](./docs/apigeecli.md) of available commands

___

## Support

This is not an officially supported Google product
