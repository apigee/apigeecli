# apigeecli

[![Go Report Card](https://goreportcard.com/badge/github.com/apigee/apigeecli)](https://goreportcard.com/report/github.com/apigee/apigeecli)
[![GitHub release](https://img.shields.io/github/v/release/apigee/apigeecli)](https://github.com/apigee/apigeecli/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is a tool to interact with [Apigee APIs](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest) for [Apigee hybrid](https://cloud.google.com/apigee/docs/hybrid/latest/what-is-hybrid) and [Apigee's managed](https://cloud.google.com/apigee/docs/api-platform/get-started/overview) offering. The tool lets you manage (Create,Get, List, Update, Delete, Export and Import) Apigee entities like proxies, products etc. The tools also helps you create Service Accounts in Google IAM to operate Apigee hybrid runtime.

## Installation

`apigeecli` is a binary and you can download the appropriate one for your platform from [here](https://github.com/apigee/apigeecli/releases)

NOTE: Supported platforms are:

* Darwin
* Windows
* Linux

## What you need to know about apigeecli

You must have an account on [Apigee](https://cloud.google.com/apigee/docs) to perform any `apigeecli` functions. These functions include: proxies, API Products, Environments, Org details etc.

You need to be familiar with basic concepts and features of Apigee such as API proxies, organizations, and environments.

For more information, refer to the [Apigee API Reference](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest).

## Available Commands

Here is a [list](./docs/apigeecli.md) of available commands

## Service Account

Create a service account with appropriate persmissions. Use `apigeecli` to create service accounts (`apigeecli iam`). Read more [here](https://cloud.google.com/apigee/docs/api-platform/system-administration/apigee-roles) about IAM roles in Apigee 

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

## Generating API Proxies from OpenAPI Specs

apigeecli allows the user to generate Apigee API Proxy bundles from an OpenAPI spec (only 3.0.x supported). The Apigee control plane does not support custom formats (ex: uuid). If you spec contains custom formats, consider the following flags

* `--formatValidation=false`: this disables validation for custom formats.
* `--skip-policy=false`: By default the OAS policy is added to the proxy (to validate API requests). By setting this to false, schema validation is not enabled and the control plane will not reject the bundle due to custom formats.

The following actions are automatically implemented when the API Proxy bundle is generated:

### Security Policies

If the spec defines securitySchemes, for ex the following snippet:

```yaml
components:
  securitySchemes:
    petstore_auth:
      type: oauth2
      flows:
        implicit:
          authorizationUrl: 'http://petstore.swagger.io/api/oauth/dialog'
          scopes:
            'write:pets': modify pets in your account
            'read:pets': read your pets
    api_key:
      type: apiKey
      name: api_key
      in: header
```

is interpreted as OAuth-v20 (verification only) policy and the VerifyAPIKey policy.


These security schemes can be added to the PreFlow by enabling the scheme globally

```yaml
security:
  - api_key: []
```

Or within a Flow Condition like this

```yaml
  '/pet/{petId}/uploadImage':
    post:
      ...
      security:
        - petstore_auth:
            - 'write:pets'
            - 'read:pets'

```

### Dynamic target endpoints

apigeecli allows the user to dynamically set a target endpoint. These is especially useful when deploying target/backend applications to GCP's serverless platforms like Cloud Run, Cloud Functions etc. apigeecli also allows the user to enable Apigee'e [Google authentication](https://cloud.google.com/apigee/docs/api-platform/security/google-auth/overview) before connecting to the backend.

#### Set a dynamic target

```sh
apigeecli apis create -n petstore -f ./test/petstore.yaml --oas-target-url-ref=propertyset.petstore.url
```

This example dynamically sets the `target.url` message context variable. This variable is retrieved from a propertyset file. It is expected the user will separately upload an environment scoped propertyset file with this key.

#### Set a dynamic target for Cloud Run

```sh
apigeecli apis create -n petstore -f ./test/petstore.yaml --oas-google-idtoken-aud-ref=propertyset.petstore.aud --oas-target-url-ref=propertyset.petstore.url
```

This example dynamically sets the Google Auth `audience` and the `target.url` message context variable. These variables are retrieved from a propertyset file. It is expected the user will separately upload an environment scoped propertyset file with these keys. If you do not wish to user a property to set these values later, you can use `--oas-google-idtoken-aud-literal` to set the audience directly in the API Proxy.

While this example shows the use of Google IDToken, Google Access Token is also supported. To use Google Access Token, use the `oas-google-accesstoken-scope-literal` flag instead.

### Traffic Management

apigeeli allow the user to add [SpikeArrest](https://cloud.google.com/apigee/docs/api-platform/reference/policies/spike-arrest-policy) or [Quota](https://cloud.google.com/apigee/docs/api-platform/reference/policies/quota-policy) policies. Since OpenAPI spec does not natively support the ability to specify such policies, a custom extension is used.

#### Quota custom extension

The following configuration allows the user to specify quota parameters in the API Proxy.

```yaml
x-google-quota:
  - name: test1 # this is appended to the quota policy name, ex: Quota-test1
    interval-literal: 1 # specify the interval in the policy, use interval-ref to specify a variable
    timeunit-literal: minute # specify the timeUnit in the policy, use timeUnit-ref to specify a variable
    allow-literal: 1 # specify the allowed rate in the policy, use allow-ref to specify a variable
```

NOTE: literals cannot be combined with variables.

The following configuration allows the user to derive quota parameters from an API Product

```yaml
x-google-quota:
  - name: test1 # this is appended to the quota policy name, ex: Quota-test1
    useQuotaConfigInAPIProduct: Verify-API-Key-api_key # specify the step name that contains the consumer identification. Must be OAuth or VerifyAPIKey step.
```

The above configurations are mutually exclusive.

#### SpikeArrest custom extension

The following configuration allows the user to specify Spike Arrest parameters in the API Proxy.

```yaml
x-google-ratelimit: 
  - name: test1 # this is appended to the quota policy name, ex: Spike-Arrest-test1
    rate-literal: 10ps # specify the allowed interval in the policy, use rate-ref to specify a variable
    identifier-ref: request.header.url #optional, specify msg ctx var for the identifier
```

### Examples

See this [OAS document](./test/petstore-ext1.yaml) for examples

## Generating API Proxies from GraphQL Schemas

apigeecli allows the user to generate Apigee API Proxy bundles from a GraphQL schema. When generating a proxy, consider the following flags:

* `--basepath`: Specify a basePath for the GraphQL proxy
* `--skip-policy=false`: By default the GraphQL policy is added to the proxy (to validate API requests). By setting this to false, schema validation is not enabled.
* `--target-url-ref`: Specify a target endpoint location variable. For ex: `--target-url-ref=propertyset.gql.url` implies the GraphQL target location is available in an environment scoped property set called `gql` and the key is `url`.  

___

## Support

This is not an officially supported Google product