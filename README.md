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
apigeeapi -t $TOKEN
```

Pass the service account

```
apigeeapi -a orgadmin.json
```

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
apigeeapi org list -t $TOKEN
```

### <a name="getorg"/> get

Get org details for an Apigee Org

```
apigeeapi org get -o org -t $TOKEN
```

## <a name="prods"/> prods

Supported alias `products`

* [list](#listproducts)
* [get](#getproduct)

### <a name="listproducts"/> list

List all API Products in the org

```
apigeeapi prods list -o org -t $TOKEN
```

### <a name ="getproduct"/> get

Get details of an API product

```
apigeeapi prods list -o org -n name -t $TOKEN
```

## <a name="devs"/> developers

Supported alias `developers`

* [list](#listdevs)
* [get](#getdev)

### <a name="listdevs"/> list

List all App Developers in an org

```
apigeeapi devs list -o org -t $TOKEN
```

### <a name ="getdev"/> get

Get details of an App Developer

```
apigeeapi devs get -o org -n name -t $TOKEN
```

## <a name="apps"/> apps

Supported alias `applications`

* [list](#listapps)
* [get](#getapp)

### <a name="listapps"/> list

List all developer apps in an org

```
apigeeapi apps list -o org -t $TOKEN
```

### <a name ="getapp"/> get

Get details of a developer app

```
apigeeapi apps get -o org -n name -t $TOKEN
```

## <a name="sf"/> sharedflows

* [list](#listsf)
* [get](#getsf)

### <a name="listsf"/> list

List all shared flows in an org

```
apigeeapi sharedflows list -o org -t $TOKEN
```

### <a name ="getsf"/> get

Get details of an App Developer

```
apigeeapi sharedflow get -o org -n name -t $TOKEN
```

## <a name="env"/> envs

Supported alias `environments`

* [list](#listenv)
* [get](#getsf)

### <a name="listenv"/> list

List all environments in an org

```
apigeeapi envs list -o org -t $TOKEN
```

### <a name ="getenv"/> get

Get details of an environment

```
apigeeapi envs get -o org -e env -t $TOKEN
```

