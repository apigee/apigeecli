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
* View Org details

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
* [products] (#products)
* [apis](#apis)
* [environment](#environment)

## <a name="org"/> org

* [list] (#listorgs)
* [get] (#getorg)

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

## <a name="products"/> products

* [list](#listproducts)
* [get](#getproduct)

### <a name="listproducts"/> list

List all API Products in the org