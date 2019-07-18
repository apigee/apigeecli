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
* [sync](#sync)
* [kvm](#kvm)
* [flowhooks](#flow)
* [targetservers](#target)

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
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name


## <a name="prods"/> prods

Supported alias `products`

* [list](#listproducts)
* [get](#getproduct)
* [create](#crtproduct)
* [delete](#delproduct)

### <a name="listproducts"/> list

List all API Products in the org

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--expand -x` (optional) Returns an expanded list of products for the organization.
* `--count -c` (optional) Number of app ids to return. Default is 1000

```
apigeeapi prods list -o org 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name


### <a name ="getproduct"/> get

Get details of an API product

```
apigeeapi prods list -o org -n name 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--name -n` (required) API product name

### <a name="crtproduct"/> create

Create an API product

```
apigeeapi prods create -o org -n name -e test,prod -p proxy1,proxy2 -f auto
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API product name
* `--displayName -m` (optional) Display name for API product
* `--approval -f` (required) Approval type for API product
* `--desc -d` (optional) Description for API product
* `--envs -e` (required) A comma separated list of environments to enable
* `--proxies -p` (required) A comma separated list of API proxies
* `--scopes -s` (optional) A comma separated list of OAuth scopes
* `--quota -q` (optional) Quota Amount
* `--interval -i` (optional) Quota Time Interval
* `--unit -u` (optional) Quota Time Unit

### <a name="delprodct"/> delete

Delete an API Product

Get details of an API product

```
apigeeapi prods delete -o org -n name 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API product name


## <a name="apis"/> apis

* [list](#listorgs)
* [listdeploy](#listdeploy)
* [create](#createapi)
* [deploy](#depapi)
* [fetch](#fetchapi)
* [delete](#delapi)
* [undeploy](#undepapi)

### <a name="listorgs"/> list

List APIs in an Apigee Org

```
apigeeapi apis list -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (optional) Apigee environment name

If the environment name is passed, lists the deployed proxies in the environment

### <a name="listdeploy"/> listdeploy

Lists all deployments of an API proxy

```
apigeeapi org listdeploy -o org -n name
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--expand -x` (optional) Returns an expanded list of proxies for the organization.
* `--count -c` (optional) Number of apis to return


### <a name="createapi"/> create

Import or create an API Proxy. If a bundle (zip) is supplied, it is imported else, it creates an empty proxy in the Org

```
apigeeapi apis create -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--proxy -p` (required) API proxy bundle (zip)

### <a name="depapi"/> deploy

Deploys a revision of an existing API proxy to an environment in an organization.

```
apigeeapi apis deploy -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 
* `--revision -v` (required) API proxy revision
* `--ovr -r` (optional) Forces deployment of the new revision.


### <a name="fetchapi"/> fetch

Returns a zip-formatted proxy bundle of code and config files.

```
apigeeapi apis fetch -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 
* `--revision -v` (required) API proxy revision

The downloaded file is {proxyname}.zip and in the folder where the command is executed

### <a name="delapi"/> delete

Deletes an API proxy and all associated endpoints, policies, resources, and revisions. The API proxy must be undeployed before you can delete it.

```
apigeeapi apis delete -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 

### <a name="undepapi"/> undeploy

Undeploys a revision of an existing API proxy to an environment in an organization.

```
apigeeapi apis undeploy -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--revision -v` (required) API proxy revision

## <a name="devs"/> developers

Supported alias `developers`

* [list](#listdevs)
* [get](#getdev)
* [create](#crtdev)
* [delete](#deldev)

### <a name="listdevs"/> list

List all App Developers in an org

```
apigeeapi devs list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name ="getdev"/> get

Get details of an App Developer

```
apigeeapi devs get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer email

## <a name="apps"/> apps

Supported alias `applications`

* [list](#listapps)
* [get](#getapp)
* [create](#crtapp)
* [delete](#delapp)
* [genkey](#genkey)

### <a name="listapps"/> list

List all developer apps in an org

```
apigeeapi apps list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--expand -x` (optional) Returns an expanded list of apps for the organization.
* `--count -c` (optional) Number of app ids to return.

### <a name ="getapp"/> get

Get details of a developer app

```
apigeeapi apps get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name

### <a name ="delapp"/> delete

Delete a developer app

```
apigeeapi apps delete -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name

### <a name ="crtapp"/> create

Delete a developer app

```
apigeeapi apps create -o org -n name -e test,prod -p proxy1 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name
* `--env -e` (required) Developer's email* 
* `--expires -x` (optional) Lifetime of the consumer's key
* `--callabck -c` (optional) OAuth callback url
* `--prods -p` (required) A comma separated list of products
* `--scopes -s` (optional) OAuthe scopes

### <a name ="genkey"/> genkey

Create new developer KeyPairs Generates a new consumer key and consumer secret for the named developer app

```
apigeeapi apps genkey -o org -n name -p proxy1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name
* `--expires -x` (optional) Lifetime of the consumer's key
* `--callabck -c` (optional) OAuth callback url
* `--prods -p` (required) A comma separated list of products
* `--scopes -s` (optional) OAuthe scopes

## <a name="sf"/> sharedflows

* [list](#listsf)
* [get](#getsf)
* [deploy](#depsf)
* [undeploy](#undepsf)

### <a name="listsf"/> list

List all shared flows in an org

```
apigeeapi sharedflows list -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (optional) Apigee environment name

When the environment name is passed, list the deployed shared flows in the environment

### <a name ="getsf"/> get

Get details of a shared flow

```
apigeeapi sharedflow get -o org -n name 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Shared flow name

### <a name ="depsf"/> deploy

Deploy a shared flow

```
apigeeapi sharedflow deploy -o org -n name -v 1 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Shared flow name
* `--revision -v` (required) Shared flow revision

### <a name ="undepsf"/> undeploy

Deploy a shared flow

```
apigeeapi sharedflow undeploy -o org -n name -v 1 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Shared flow name
* `--revision -v` (required) Shared flow revision

## <a name="env"/> envs

Supported alias `environments`

* [list](#listenv)
* [get](#getenv)

### <a name="listenv"/> list

List all environments in an org

```
apigeeapi envs list -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name ="getenv"/> get

Get details of an environment

```
apigeeapi envs get -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) environment name

## <a name="sync"/> sync

* [set](#setsync)
* [get](#getsync)

### <a name="listorgs"/> set

Set identity with access to control plane resources

```
apigeeapi sync set -o org -i identity 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--ity -i` (required) IAM Identity

### <a name="listorgs"/> get

List all the orgs available to the identity (service account)

```
apigeeapi sync set -o org
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name


## <a name="kvm"/> kvm

* [create](#crtkvm)
* [list](#listkvm)
* [delete](#delkvm)

### <a name="crtkvm"/> create

Create a new environment scoped KVM Map

```
apigeeapi kvms create -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) KVM Map name
* `--encrypt -c` (optional) Enable encrypted KVM

### <a name="listkvm"/> list

List all the KVM Maps in an environment

```
apigeeapi kvms create -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

### <a name="delkvm"/> delete

Delete a KVM Map

```
apigeeapi kvms create -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
