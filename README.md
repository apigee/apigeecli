# apigeeapi

This is a tool to interact with [Apigee APIs](https://apigee.googleapis.com). The tool let you manage (get, list) environments, proxies, etc.

# Installation

`apigeeapi` is a binary and you can download the appropriate one for your platform from [here](https://github.com/srinandan/apigeeapi/releases)

NOTE: Supported platforms are:
* Darwin
* Windows
* Linux

# What you need to know about apigeeapi

You must have an account on [Apigee Hybrid](https://docs.apigee.com/hybrid/beta2) to perform any `apigeeapi` functions. These functions include:

* Manage API proxies,
* Manage API Products,
* Manage Apigee Environments,
* View Org details etc.

You need to be familiar with basic concepts and features of Apigee Edge such as API proxies, organizations, and environments.

For more information, refer to the [Apigee Hybrid API Reference](https://docs.apigee.com/hybrid/beta2/reference/apis/rest/index).

## Service Account

Create a service account with appropriate persmissions. Refer to this [link](https://docs.apigee.com/hybrid/beta2/precog-serviceaccounts) for more details on how to download the JSON file.

## Access Token

`apigeeapi` can use the service account directly and obtain an access token. However, the user can also provide an access token. 

Print the access token 

```
export GOOGLE_APPLICATION_CREDENTIALS=orgadmin.json
gcloud auth application-default print-access-token
```

Use this access token for all subsequent calls (token expires in 1 hour)

# Command Reference

The following options are available for security

Pass the access token
```
apigeeapi 
```

Pass the service account

```
apigeeapi -a orgadmin.json
```

## Access Token Caching

`apigeeapi` caches the OAuth Access token for subsequent calls (until the token expires). The access token is stored in `$HOME/.access_token`. This path must be readable/writeable by the `apigeeapi` process. 

### Example 1

First command

```
apigeeapi orgs list -a orgadmin.json
```

Subsequent command (no token or service account)

```
apigeeapi orgs get -o org
```
### Example 2

First command

```
apigeeapi orgs list -t $TOKEN
```

Subsequent command (no token or service account)

```
apigeeapi orgs get -o org
```

## Supported entites

* [org](#org)
* [prods](#prods)
* [apis](#apis)
* [envs](#env)
* [sharedflows](#sf)
* [apps](#apps)
* [devs](#devs)

## <a name="org"/> org

* [list](#listorgs)
* [get](#getorg)

### <a name="listorgs"/> list

List all the orgs available to the identity (service account)

```
apigeeapi org list 
```

### <a name="getorg"/> get

Get org details for an Apigee Org

```
apigeeapi org get -o org 
```

## <a name="prods"/> prods

Supported alias `products`

* [list](#listproducts)
* [get](#getproduct)

### <a name="listproducts"/> list

List all API Products in the org

Optional parameters:
* `--expand -x` : `true` or `false` - Optional. Returns an expanded list of products for the organization.
* `--count -c` - Optional. Number of app ids to return. Default is 1000

```
apigeeapi prods list -o org 
```

### <a name ="getproduct"/> get

Get details of an API product

```
apigeeapi prods list -o org -n name 
```

## <a name="apis"/> apis

* [list](#listorgs)
* [listdeploy](#listdeploy)

### <a name="listorgs"/> list

List APIs in an Apigee Org

```
apigeeapi org list 
```

### <a name="listdeploy"/> listdeploy

Lists all deployments of an API proxy

Optional parameters:
* `--expand -x` : `true` or `false` - Optional. Returns an expanded list of developers for the organization.
* `--count -c` - Optional. Number of app ids to return. Default is 1000


```
apigeeapi org listdeploy -o org 
```

## <a name="devs"/> developers

Supported alias `developers`

* [list](#listdevs)
* [get](#getdev)

### <a name="listdevs"/> list

List all App Developers in an org

```
apigeeapi devs list -o org 
```

### <a name ="getdev"/> get

Get details of an App Developer

```
apigeeapi devs get -o org -n name 
```

## <a name="apps"/> apps

Supported alias `applications`

* [list](#listapps)
* [get](#getapp)

### <a name="listapps"/> list

List all developer apps in an org

Optional parameters:
* `--expand -x` : `true` or `false` - Optional. Returns an expanded list of apps for the organization.
* `--count -c` - Optional. Number of app ids to return. Default is 10000

```
apigeeapi apps list -o org 
```

### <a name ="getapp"/> get

Get details of a developer app

```
apigeeapi apps get -o org -n name 
```

## <a name="sf"/> sharedflows

* [list](#listsf)
* [get](#getsf)

### <a name="listsf"/> list

List all shared flows in an org

```
apigeeapi sharedflows list -o org 
```

### <a name ="getsf"/> get

Get details of an App Developer

```
apigeeapi sharedflow get -o org -n name 
```

## <a name="env"/> envs

Supported alias `environments`

* [list](#listenv)
* [get](#getsf)

### <a name="listenv"/> list

List all environments in an org

```
apigeeapi envs list -o org 
```

### <a name ="getenv"/> get

Get details of an environment

```
apigeeapi envs get -o org -e env 
```

