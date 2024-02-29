# About apigeecli

[![Go Report Card](https://goreportcard.com/badge/github.com/apigee/apigeecli)](https://goreportcard.com/report/github.com/apigee/apigeecli)
[![GitHub release](https://img.shields.io/github/v/release/apigee/apigeecli)](https://github.com/apigee/apigeecli/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is a tool to interact with [Apigee APIs](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest) for [Apigee hybrid](https://cloud.google.com/apigee/docs/hybrid/latest/what-is-hybrid) and [Apigee's managed](https://cloud.google.com/apigee/docs/api-platform/get-started/overview) offering. The tool lets you manage (Create,Get, List, Update, Delete, Export and Import) Apigee entities like proxies, products etc. The tools also helps you create Service Accounts in Google IAM to operate Apigee hybrid runtime.

## Installation

`apigeecli` is a binary and you can download the appropriate one for your platform from [here](https://github.com/apigee/apigeecli/releases). Run this script to download & install the latest version (on Linux or Darwin)

```sh
curl -L https://raw.githubusercontent.com/apigee/apigeecli/main/downloadLatest.sh | sh -
```

## Getting Started

### User Tokens

The simplest way to get started with `apigeecli` is

```
token=$(gcloud auth print-access-token)

apigeecli orgs list -t $token
```

### Metadata OAuth2 Access Tokens

If you are using `apigeecli` on Cloud Shell, GCE instances, Cloud Build, then you can use the metadata to get the access token

```sh
apigeecli orgs list --metadata-token
```

### Google Default Application Credentials

You can configure gcloud to setup/create default application credentials. These credentials can be used by `apigeecli`.

```sh
gcloud auth application-default login
apigeecli orgs list --default-token
```

or through impersonation

```sh
gcloud auth application-default login --impersonate-service-account <SA>
apigeecli orgs list --default-token
```

### Access Token Generation from Service Accounts

`apigeecli` can use the service account directly and obtain an access token.

```bash
apigeecli token gen -a serviceaccount.json
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format

Use this access token for all subsequent calls (token expires in 1 hour)

### Access Token Caching

`apigeecli` caches the OAuth Access token for subsequent calls (until the token expires). The access token is stored in `$HOME/.apigeecli`. This path must be readable/writeable by the `apigeecli` process.

```bash
apigeecli token cache -a serviceaccount.json
```

or

```bash
token=$(gcloud auth print-access-token)
apigeecli token cache -t $token
```

or

```bash
apigeecli token cache --metadata-token
```

## Set Preferences

If you are using the same GCP project for Apigee, then consider setting up preferences so they don't have to be included in every command. Preferences are written to the `$HOME/.apigeecli` folder

```
project=$(gcloud config get-value project | head -n 1)

apigeecli prefs set -o $project
```

Subsequent commands can be like this:

```
token=$(gcloud auth print-access-token)
apigeecli orgs get -t $token #fetches the org details of the org set in preferences
```

The following preferences can be set:

| Flag                   | Description                                           |
| -----------------------| ----------------------------------------------------- |
| `-g, --github string`  | On premises Github URL                                |
| `-o, --org string`     | Apigee organization name                              |
| `-p, --proxy string`   | Use http proxy before contacting the control plane    |
| `--nocheck`            | Don't check for newer versions of cmd                 |

## Container download

The lastest container version for apigeecli can be downloaded via

```sh
docker pull ghcr.io/apigee/apigeecli:latest
```

### Using docker to run commands locally

```sh
docker run -ti ghcr.io/apigee/apigeecli:latest orgs list -t $token
```

### Using apigeecli with Cloud Build

To execute apigeecli commands in cloud build,

```
steps:
- id: 'Run apigeecli commands'
  name: ghcr.io/apigee/apigeecli:latest
  args:
  - 'orgs'
  - 'get'
  - '-o'
  - '$PROJECT_ID'
  - '--metadata-token'
```

If you need the response from the previous command as input to the next, then take advantage of `sh` and `jq` like so:

```yaml
steps:
- id: 'Run apigeecli commands'
  name: ghcr.io/apigee/apigeecli:latest
  entrypoint: 'sh'
  args:
    - -c
    - |
      #setup preferences
      apigeecli prefs set --nocheck=true -o $PROJECT_ID
      apigeecli token cache --metadata-token

      # run other commands here
      apigeecli orgs list | jq
```

### Access shell

```sh
docker run -ti --entrypoint sh ghcr.io/apigee/apigeecli:latest
```

See this [page](https://github.com/apigee/apigeecli/pkgs/container/apigeecli) for other versions.

## What you need to know about apigeecli

You must have an account on [Apigee](https://cloud.google.com/apigee/docs) to perform any `apigeecli` functions. These functions include: proxies, API Products, Environments, Org details etc.

You need to be familiar with basic concepts and features of Apigee such as API proxies, organizations, and environments.

For more information, refer to the [Apigee API Reference](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest).

## Available Commands

Here is a [list](./docs/apigeecli.md) of available commands

## Enviroment Variables

The following environment variables may be set to control the behavior of `apigeecli`. The default values are all `false`

* `APIGEECLI_DEBUG=true` enables debug log
* `APIGEECLI_SKIPCACHE=true` will not cache the access token on the disk
* `APIGEECLI_ENABLE_RATELIMIT=true` enables rate limiting when making Apigee APIs (at 1 API call every 100ms)
* `APIGEECLI_NO_USAGE=true` does not print usage when the command fails
* `APIGEECLI_NO_ERRORS=true` does not print error messages from the CLI (control plane error messages are displayed)
* `APIGEECLI_DRYRUN=true` does not execute Apigee control plane APIs

## Generating API Proxies

`apigeecli` can generate API proxies from:

* OpenAPI 3.0/3.1 Specification
* GraphQL Schema
* A template/stub for Application Integration
* Cloud Endpoints/API Gateway OpenAPI 2.0 specification

### Generating API Proxies from OpenAPI Specs

`apigeecli` allows the user to generate Apigee API Proxy bundles from an OpenAPI spec (3.0.x and 3.1.x are supported). The Apigee control plane does not support custom formats (ex: uuid). If you spec contains custom formats, consider the following flags

* `--add-cors=true`: Add a CORS policy
* `--formatValidation=false`: this disables validation for custom formats.
* `--skip-policy=false`: By default the OAS policy is added to the proxy (to validate API requests). By setting this to false, schema validation is not enabled and the control plane will not reject the bundle due to custom formats.

**NOTE**: The Apigee runtime does not support OAS 3.1.x (specifically the OpenAPI Validation policy). When using OAS 3.1.x, `apigeecli` generates the proxy, but does not include the OAS validation policy.

The following actions are automatically implemented when the API Proxy bundle is generated:

#### Security Policies

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

#### Dynamic target endpoints

`apigeecli` allows the user to dynamically set a target endpoint. These is especially useful when deploying target/backend applications to GCP's serverless platforms like Cloud Run, Cloud Functions etc. `apigeecli` also allows the user to enable Apigee'e [Google authentication](https://cloud.google.com/apigee/docs/api-platform/security/google-auth/overview) before connecting to the backend.

##### Set a dynamic target

```sh
apigeecli apis create -n petstore -f ./test/petstore.yaml --oas-target-url-ref=propertyset.petstore.url
```

This example dynamically sets the `target.url` message context variable. This variable is retrieved from a propertyset file. It is expected the user will separately upload an environment scoped propertyset file with this key.

##### Set a dynamic target for Cloud Run

```sh
apigeecli apis create -n petstore -f ./test/petstore.yaml --oas-google-idtoken-aud-ref=propertyset.petstore.aud --oas-target-url-ref=propertyset.petstore.url
```

This example dynamically sets the Google Auth `audience` and the `target.url` message context variable. These variables are retrieved from a propertyset file. It is expected the user will separately upload an environment scoped propertyset file with these keys. If you do not wish to user a property to set these values later, you can use `--oas-google-idtoken-aud-literal` to set the audience directly in the API Proxy.

While this example shows the use of Google IDToken, Google Access Token is also supported. To use Google Access Token, use the `oas-google-accesstoken-scope-literal` flag instead.

#### Traffic Management

apigeeli allow the user to add [SpikeArrest](https://cloud.google.com/apigee/docs/api-platform/reference/policies/spike-arrest-policy) or [Quota](https://cloud.google.com/apigee/docs/api-platform/reference/policies/quota-policy) policies. Since OpenAPI spec does not natively support the ability to specify such policies, a custom extension is used.

##### Quota custom extension

NOTE: This extension behaves differently when used with Swagger for Cloud Endpoints/API Gateway.

The following configuration allows the user to specify quota parameters in the API Proxy.

```yaml
x-google-quota:
  - name: test1 # this is appended to the quota policy name, ex: Quota-test1
    interval-literal: 1 # specify the interval in the policy, use interval-ref to specify a variable
    timeunit-literal: minute # specify the timeUnit in the policy, use timeUnit-ref to specify a variable
    allow-literal: 1 # specify the allowed rate in the policy, use allow-ref to specify a variable
    identifier-literal: request.headers.api_key # optionally, set an identifier. If not set, the proxy will count every msg. Identifiers are variables in apigee
```

NOTE: literals cannot be combined with variables.

The following configuration allows the user to derive quota parameters from an API Product

```yaml
x-google-quota:
  - name: test1 # this is appended to the quota policy name, ex: Quota-test1
    useQuotaConfigInAPIProduct: Verify-API-Key-api_key # specify the step name that contains the consumer identification. Must be OAuth or VerifyAPIKey step.
```

The above configurations are mutually exclusive.

##### SpikeArrest custom extension

The following configuration allows the user to specify Spike Arrest parameters in the API Proxy.

```yaml
x-google-ratelimit:
  - name: test1 # this is appended to the quota policy name, ex: Spike-Arrest-test1
    rate-literal: 10ps # specify the allowed interval in the policy, use rate-ref to specify a variable
    identifier-ref: request.header.url #optional, specify msg ctx var for the identifier
```

#### Examples

See this [OAS document](./test/petstore-ext1.yaml) for examples

### Generating API Proxies from GraphQL Schemas

`apigeecli` allows the user to generate Apigee API Proxy bundles from a GraphQL schema. When generating a proxy, consider the following flags:

* `--basepath`: Specify a basePath for the GraphQL proxy
* `--skip-policy=false`: By default the GraphQL policy is added to the proxy (to validate API requests). By setting this to false, schema validation is not enabled.
* `--target-url-ref`: Specify a target endpoint location variable. For ex: `--target-url-ref=propertyset.gql.url` implies the GraphQL target location is available in an environment scoped property set called `gql` and the key is `url`

### Generating an API Proxy template for Application Integration

`apigeecli` allows the user to generate an Apigee API Proxy bundle template for [Application Integration](https://cloud.google.com/application-integration/docs/overview). When generating the proxy, consider the following flags:

* `--trigger`: Specify the API trigger name of the Integration. This is also used as the basePath. Don't include `api_trigger/`
* `--integration`: Specify the Name of the Integration
* `--name`: Specify the Name of the API Proxy

### Generating an API Proxy template from Cloud Endpoints/API Gateway artifacts

[Cloud Endpoints](https://cloud.google.com/endpoints) and [API Gateway](https://cloud.google.com/api-gateway) use the [extensions](https://cloud.google.com/endpoints/docs/openapi/openapi-extensions) specified in an OpenAPI (FKA Swagger) document to apply policies. `apigeecli` can use those extensions to generate a similar Apigee API Proxy bundle. When generating the proxy, consider the following flags:

* `--add-cors=true`: Add a CORS policy

#### Limitations

* The `protocol` property in `x-google-backend` is ignored. All upstream/backend is treated as http 1.1
* The `metrics` property in `x-google-management` is not supported
* The quota unit is ignored in `x-google-management` is ignored. See below for quota behavior
* The extension `x-google-endpoints` is ignored. To add CORS, see above
* [Mutiple security requirements](https://cloud.google.com/endpoints/docs/openapi/openapi-limitations#multiple_security_requirements): If more than one security policy, regardless of the security type, is set on a path, then the **first one** is enabled. In the following examples,

```
  /hello:
    get:
      operationId: hello
      security:
        - google_id_token: []
        - api_key: []
```

**or** in the following case,

```
  /hello:
    get:
      operationId: hello
      security:
        - google_id_token: []
          api_key: []
```

it cannot be determined which policy is applied. It is best to avoid Swagger documents with such configurations.

* If more than one `x-google-jwt-locations` are specified, then the first one is used. In the following example,

```
x-google-jwt-locations:
  # Expect header "Authorization": "MyBearerToken <TOKEN>"
  - header: "Authorization"
    value_prefix: "MyBearerToken "
  # expect query parameter "jwt_query_bar=<TOKEN>"
  - query: "jwt_query_bar"
```

query parameters are ignored. By default, if no location is specified, the JWT location is the `Authorization` header and value_prefix is `Bearer <token>`

## How do I import entities using apigeecli

| Operations | Import command |
|---|---|
| apicategories | ``` apigeecli apicategories import -o $org -t $token -f samples/apicategories.json -s $siteId ``` |
| apis | apigeecli apis import -f samples/apis -o $org -t $token |
| appgroups | apigeecli appgroups import -f samples/appgroups.json -o $org -t $token |
| datacollectors | apigeecli datacollectors import -f samples/datacollectors.json -o $org -t $token |
| developers | apigeecli developers import -f samples/developers.json -o $org -t $token |
| kvms | Rename the files under samples/kvms to match your Apigee setup  apigeecli kvms import -f samples/kvms -o $org -t $token |
| products | apigeecli products import -f samples/apiproduct-legacy.json -o $org -t $token apigeecli products import -f samples/apiproduct-gqlgroup.json -o $org -t $token apigeecli products import -f samples/apiproduct-op-group.json -o $org -t $token |
| sharedflows | apigeecli sharedflows import -f samples/sharedflows -o $org -t $token |
| targetservers | apigeecli targetservers import -f samples/targetservers.json -o $org -t $token -e $env |
| keystores | apigeecli keystores import -f samples/keystores.json -o $org -t $token -e $env |
| references | apigeecli references import -f samples/references.json -o $org -t $token -e $env |
| apps | Work In Progress |
| apidocs | Work In Progress |
| environments | Work In Progress |
| organizations | Work In Progress |
| securityprofiles | Work In Progress |

## How do I verify the binary?

All artifacts are signed by [cosign](https://github.com/sigstore/cosign). We recommend verifying any artifact before using them.

You can use the following public key to verify any `apigeecli` binary with:

```sh
cat cosign.pub
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgjcKEyPi18vd6Zk5/ggAkH6CLSy3
C8gzi5q3xsycjI7if5FABk7bfciR4+g32H8xTl4mVHhHuz6I6FBG24/nuQ==
-----END PUBLIC KEY-----

cosign verify-blob --key=cosign.pub --signature apigeecli_<platform>_<arch>.zip.sig apigeecli_<platform>_<arch>.zip
```

Where `platform` can be one of `Darwin`, `Linux` or `Windows` and arch (architecture) can be one of `arm64` or `x86_64`

## How do I verify the apigeecli containers?

All images are signed by [cosign](https://github.com/sigstore/cosign). We recommend verifying any container before using them.

```sh
cat cosign.pub
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgjcKEyPi18vd6Zk5/ggAkH6CLSy3
C8gzi5q3xsycjI7if5FABk7bfciR4+g32H8xTl4mVHhHuz6I6FBG24/nuQ==
-----END PUBLIC KEY-----

cosign verify --key=cosign.pub ghcr.io/apigee/apigeecli:latest
```

___

## Support

This is not an officially supported Google product
