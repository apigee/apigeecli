# apigeecli
[![TravisCI](https://travis-ci.org/srinandan/apigeecli.svg?branch=master)](https://travis-ci.org/srinandan/apigeecli)
[![Go Report Card](https://goreportcard.com/badge/github.com/srinandan/apigeecli)](https://goreportcard.com/report/github.com/srinandan/apigeecli)
[![Version](https://img.shields.io/badge/version-v0.3-green.svg)](https://github.com/srinandan/apigeecli/releases)

This is a tool to interact with [Apigee APIs](https://apigee.googleapis.com). The tool lets you manage (get, list) environments, proxies, etc.

# Installation

`apigeecli` is a binary and you can download the appropriate one for your platform from [here](https://github.com/srinandan/apigeecli/releases)

NOTE: Supported platforms are:
* Darwin
* Windows
* Linux

# What you need to know about apigeecli

You must have an account on [Apigee Hybrid](https://docs.apigee.com/hybrid/beta2) to perform any `apigeecli` functions. These functions include: proxies, API Products, Environments,
Org details etc.

You need to be familiar with basic concepts and features of Apigee Edge such as API proxies, organizations, and environments.

For more information, refer to the [Apigee Hybrid API Reference](https://docs.apigee.com/hybrid/beta2/reference/apis/rest/index).

## Service Account

Create a service account with appropriate persmissions. Refer to this [link](https://docs.apigee.com/hybrid/beta2/precog-serviceaccounts) for more details on how to download the JSON file.

## Access Token

`apigeecli` can use the service account directly and obtain an access token. 

```
apigeecli token gen -a serviceaccount.json 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format


Use this access token for all subsequent calls (token expires in 1 hour)

# Command Reference

The following options are available for security

Pass the access token
```
apigeecli <flags> -t $TOKEN
```

Pass the service account

```
apigeecli <flags> -a orgadmin.json
```

## Access Token Caching

`apigeecli` caches the OAuth Access token for subsequent calls (until the token expires). The access token is stored in `$HOME/.access_token`. This path must be readable/writeable by the `apigeecli` process. 

```
apigeecli token cache -a serviceaccount.json
```

Subsequent commands do not need the token flag

___

## Supported entites

* [apis](#apis)
* [apps](#apps)
* [developers](#devs)
* [envs](#env)
* [flowhooks](#flow)
* [kvms](#kvm)
* [org](#org)
* [products](#prods)
* [sharedflows](#sf)
* [sync](#sync)
* [targetservers](#target)
* [token](#token)
* [TODO](#todo)

---

## <a name="apis"/> apis

* [create](#createapi)
* [delete](#delapi)
* [deploy](#depapi)
* [fetch](#fetchapi)
* [import](#impapis)
* [list](#listorgs)
* [listdeploy](#listdeploy)
* [undeploy](#undepapi)

### <a name="createapi"/> create

Import or create an API Proxy. If a bundle (zip) is supplied, it is imported else, it creates an empty proxy in the Org

```
apigeecli apis create -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--proxy -p` (required) API proxy bundle (zip)

### <a name="delapi"/> delete

Deletes an API proxy and all associated endpoints, policies, resources, and revisions. The API proxy must be undeployed before you can delete it.

```
apigeecli apis delete -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 

### <a name="depapi"/> deploy

Deploys a revision of an existing API proxy to an environment in an organization.

```
apigeecli apis deploy -o org -e env -n proxy -v 1
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
apigeecli apis fetch -o org -e env -n proxy -v 1
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 
* `--rev -v` (required) API proxy revision
* `--ovr -r` (optional) Forces deployment of the new revision.

### <a name="impapis"/> import

Upload a folder containing API proxy bundles

```
apigeecli apis fetch -o org -f /tmp
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--folder -f` (required) path containing API proxy bundles

### <a name="listorgs"/> list

List APIs in an Apigee Org

```
apigeecli apis list -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (optional) Apigee environment name

If the environment name is passed, lists the deployed proxies in the environment

### <a name="listdeploy"/> listdeploy

Lists all deployments of an API proxy

```
apigeecli org listdeploy -o org -n name
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--expand -x` (optional) Returns an expanded list of proxies for the organization.
* `--count -c` (optional) Number of apis to return

### <a name="undepapi"/> undeploy

Undeploys a revision of an existing API proxy to an environment in an organization.

```
apigeecli apis undeploy -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--rev -v` (required) API proxy revision

---

## <a name="apps"/> apps

Supported alias `applications`

* [create](#crtapp)
* [delete](#delapp)
* [genkey](#genkey)
* [get](#getapp)
* [list](#listapps)
* [import](#impapps)

### <a name ="crtapp"/> create

Create a developer app

```
apigeecli apps create -o org -n name -e test,prod -p proxy1 --attrs "foo1=bar1,foo2=bar2"
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name
* `--email -e` (required) Developer's email 
* `--expires -x` (optional) Lifetime of the consumer's key
* `--callabck -c` (optional) OAuth callback url
* `--prods -p` (optional) A comma separated list of products
* `--scopes -s` (optional) OAuthe scopes
* `--attrs` (optional) Custom Attributes

### <a name ="delapp"/> delete

Delete a developer app

```
apigeecli apps delete -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name

### <a name ="genkey"/> genkey

Create new developer KeyPairs Generates a new consumer key and consumer secret for the named developer app

```
apigeecli apps genkey -o org -n name -p proxy1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name
* `--expires -x` (optional) Lifetime of the consumer's key
* `--callabck -c` (optional) OAuth callback url
* `--prods -p` (required) A comma separated list of products
* `--scopes -s` (optional) OAuthe scopes

### <a name ="getapp"/> get

Get details of a developer app

```
apigeecli apps get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--appId -i` (optional) Developer App Id
* `--name -n` (optional) Developer App Name

NOTE: Either appId or Name must be passed

### <a name="listapps"/> list

List all developer apps in an org

```
apigeecli apps list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--expand -x` (optional) Returns an expanded list of apps for the organization.
* `--count -c` (optional) Number of app ids to return. Default is 10000

### <a name="impapps"/> import

Import Developer app entities into an org
NOTE: This feature is a WIP. It is not fully implemented

```
apigeecli apps import -o org -f filepath.json
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--file -f` (required) A json file containing apps

A sample file format can be found [here](https://github.com/srinandan/apigeecli/blob/master/test/apps_config.json)

---

## <a name="devs"/> developers

Supported alias `developers`

* [create](#crtdev)
* [delete](#deldev)
* [get](#getdev)
* [list](#listdevs)
* [import](#impdev)

### <a name ="crtdev"/> create

Create a new App Developer

```
apigeecli devs create -o org -n email -f firstname -s lastname -u username --attrs "foo1=bar1,foo2=bar2"
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--email -n` (required) Developer email
* `--first -f` (required) Developer firstname
* `--last -s` (required) Developer lastname
* `--user -u` (required) Developer username
* `--attrs` (optional) Custom Attributes

### <a name ="deldev"/> delete

Delete an App Developer

```
apigeecli devs get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer email


### <a name ="getdev"/> get

Get details of an App Developer

```
apigeecli devs get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer email

### <a name="listdevs"/> list

List all App Developers in an org

```
apigeecli devs list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name="impdev"/> import

Import Developer entities into an org
NOTE: This feature is a WIP. It is not fully implemented

```
apigeecli devs import -o org -f filepath.json
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--file -f` (required) A json file containing developers

---

## <a name="env"/> envs

Supported alias `environments`

* [list](#listenv)
* [get](#getenv)

### <a name ="getenv"/> get

Get details of an environment

```
apigeecli envs get -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) environment name

### <a name="listenv"/> list

List all environments in an org

```
apigeecli envs list -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

---

## <a name="flowhooks"/> flowhooks

* [attach](#crtfh)
* [detach](#delfh)
* [get](#getfh)
* [list](#listfh)

### <a name ="crtfh"/> attach

Attach a Flowhook

```
apigeecli flowhooks attach -o org -e env -n PreFlow -n proxy 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Name of the flowhook
* `--desc -d` (optional) Description the flowhook
* `--sharedflow -s` (required) Name of the shared flow
* `--continue -c` (optional) Continue on error

### <a name ="delfh"/> detach

Detach a Flowhook

```
apigeecli flowhooks detach -o org -e env -n PreFlow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Name of the flowhook


### <a name ="getfh"/> get

Get a details of a configured Flowhook

```
apigeecli flowhooks get -o org -e env -n PreFlow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Name of the flowhook

### <a name ="listfh"/> list

List of configured Flowhooks

```
apigeecli flowhooks list -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

---

## <a name="kvm"/> kvms

* [create](#crtkvm)
* [delete](#delkvm)
* [list](#listkvm)

### <a name ="crtkvm"/> create

Create a new KV Map

```
apigeecli kvms create -o org -e env -n kvm1 -c true 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) KVM Map name
* `--encrypt -c` (required) encrypted true or false

### <a name ="delkvm"/> delete

Delete a new KV Map

```
apigeecli kvms delete -o org -e env -n kvm1 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) KVM Map name

### <a name ="listkvm"/> list

List KVMs in an environment

```
apigeecli kvms list -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

---

## <a name="org"/> org

* [list](#listorgs)
* [get](#getorg)
* [setmart](#setmart)

### <a name="listorgs"/> list

List all the orgs available to the identity (service account)

```
apigeecli org list 
```

### <a name="getorg"/> get

Get org details for an Apigee Org

```
apigeecli org get -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name="setmart"/> setmart

Configure MART endpoint for an Apigee Org

```
apigeecli org get -o org -m http://endpoint
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--mart -m` (required) MART endpoint
* `--whitelist -w` (optional) Enable/disable whitelisting of GCP IP for source connections to MART

---

## <a name="prods"/> products

Supported alias `prods`

* [create](#crtproduct)
* [delete](#delproduct)
* [get](#getproduct)
* [list](#listproducts)
* [import](#impproducts)

### <a name="crtproduct"/> create

Create an API product

```
apigeecli prods create -o org -n name -e test,prod -p proxy1,proxy2 -f auto --attrs "foo1=bar1,foo2=bar2"
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
* `--attrs` (optional) Custom Attributes

### <a name="delprodct"/> delete

Delete an API Product

Get details of an API product

```
apigeecli prods delete -o org -n name 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API product name

### <a name ="getproduct"/> get

Get details of an API product

```
apigeecli prods list -o org -n name 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--name -n` (required) API product name

### <a name ="impproducts"/> import

Import API products from a configuration file

```
apigeecli prods list -o org -f file -c connections 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--file -f` (required) File containing API products
`--conn -c` (optional) Number of connections to establish; default is 4

A sample file format can be found [here](https://github.com/srinandan/apigeecli/blob/master/test/products_config.json)

---

## <a name="sf"/> sharedflows

* [create](#createsf)
* [delete](#delsf)
* [deploy](#depsf)
* [fetch](#fetchsf)
* [get](#gettsf)
* [list](#listsf)
* [undeploy](#undepsf)

### <a name="createsf"/> create

Import or create a sharedflow. If a bundle (zip) is supplied, it is imported else, it creates an empty proxy in the Org

```
apigeecli apis create -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name
* `--proxy -p` (required) sharedflow bundle (zip)

### <a name="delsf"/> delete

Deletes a sharedflow and all policies, resources, and revisions. The sharedflow must be undeployed before you can delete it.

```
apigeecli sharedflows delete -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 

### <a name="depsf"/> deploy

Deploys a revision of an existing sharedflow to an environment in an organization.

```
apigeecli sharedflow deploy -o org -e env -n sharedflow1 -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 
* `--rev -v` (required) sharedflow revision
* `--ovr -r` (optional) Forces deployment of the new revision.

### <a name="fetchsf"/> fetch

Returns a zip-formatted proxy bundle of code and config files.

```
apigeecli apis fetch -o org -e env -n sharedflow -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 
* `--rev -v` (required) API proxy revision

### <a name="getsf"/> get

Get a sharedflow's details

```
apigeecli apis get -o org -e env -n sharedflow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 

### <a name="listsf"/> list

List all sharedflows in an org

```
apigeecli apis get -o org -e env -n sharedflow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 

### <a name="undepsf"/> undeploy

Undeploys a revision of an existing API proxy to an environment in an organization.

```
apigeecli apis undeploy -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name
* `--rev -v` (required) sharedflow revision

---

## <a name="sync"/> sync

* [set](#setsync)
* [get](#getsync)


### <a name="listorgs"/> get

List all the orgs available to the identity (service account)

```
apigeecli sync get -o org
```

### <a name="listorgs"/> set

Set identity with access to control plane resources

```
apigeecli sync set -o org -i identity 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--ity -i` (required) IAM Identity

---

## <a name="target"/> targetservers

Supported alias `ts`

* [create](#crtts)
* [delete](#delts)
* [get](#getts)
* [list](#listts)

### <a name="createts"/> create

Create a new target server

```
apigeecli targetservers create -o org -e env -h hostname -p 80 -n ts1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Target server name
* `--desc -d` (optional) Description
* `--host -s` (required) Hostname
* `--port -p` (optional) Port number
* `--enable -b` (optional) Enable or disable

### <a name ="delts"/> delete

Delete a target server

```
apigeecli targetservers get -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Target server name

### <a name ="getts"/> get

Get details of a target server

```
apigeecli targetservers get -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Target server name

### <a name ="listts"/> list

List target servers in an environment

```
apigeecli targetservers list -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

---

## <a name="token"/> token

* [gen](#gentk)
* [cache](#cachetk)

### <a name ="gettk"/> gen

Generate a new access token

```
apigeecli token gen -a serviceaccount.json 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format

### <a name ="cachetk"/> cache

Caches a new access token. Writes the access token to $HOME/.access_token for use by subsequent calls

```
apigeecli token cache -a serviceaccount.json 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format

---

## <a name="todo"/> TODO

* `apigeecli clean`
* `apigeecli export`