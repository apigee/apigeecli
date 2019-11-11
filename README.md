# apigeecli

[![TravisCI](https://travis-ci.org/srinandan/apigeecli.svg?branch=master)](https://travis-ci.org/srinandan/apigeecli)
[![Go Report Card](https://goreportcard.com/badge/github.com/srinandan/apigeecli)](https://goreportcard.com/report/github.com/srinandan/apigeecli)
[![Version](https://img.shields.io/badge/version-v0.8-green.svg)](https://github.com/srinandan/apigeecli/releases)

This is a tool to interact with [Apigee APIs](https://apigee.googleapis.com). The tool lets you manage (get, list) environments, proxies, etc. The tools also helps you create Service Accounts in Google IAM to operate Apigee hybrid runtime.

## Installation

`apigeecli` is a binary and you can download the appropriate one for your platform from [here](https://github.com/srinandan/apigeecli/releases)

NOTE: Supported platforms are:

* Darwin
* Windows
* Linux

## What you need to know about apigeecli

You must have an account on [Apigee Hybrid](https://docs.apigee.com/hybrid/beta2) to perform any `apigeecli` functions. These functions include: proxies, API Products, Environments,
Org details etc.

You need to be familiar with basic concepts and features of Apigee Edge such as API proxies, organizations, and environments.

For more information, refer to the [Apigee Hybrid API Reference](https://docs.apigee.com/hybrid/beta2/reference/apis/rest/index).

## Service Account

Create a service account with appropriate persmissions. Refer to this [link](https://docs.apigee.com/hybrid/beta2/precog-serviceaccounts) for more details on how to download the JSON file.

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

`apigeecli` caches the OAuth Access token for subsequent calls (until the token expires). The access token is stored in `$HOME/.access_token`. This path must be readable/writeable by the `apigeecli` process. 

```bash
apigeecli token cache -a serviceaccount.json
```

Subsequent commands do not need the token flag

___

## Supported entites

* [apis](#apis)
* [apps](#apps)
* [cache](#cache)
* [developers](#devs)
* [envs](#env)
* [flowhooks](#flow)
* [iam](#iam)
* [keystores](#keystores)
* [keyaliases](#keyaliases)
* [kvms](#kvm)
* [org](#org)
* [projects](#projects)
* [products](#prods)
* [resources](#resources)
* [sharedflows](#sf)
* [sync](#sync)
* [targetservers](#target)
* [token](#token)

___

## <a name="apis"/> apis

* [create](#createapi)
* [delete](#delapi)
* [deploy](#depapi)
* [fetch](#fetchapi)
* [import](#impapis)
* [get](#getapi)
* [export](#expapis)
* [list](#listorgs)
* [listdeploy](#listdeploy)
* [trace](#trace)
* [undeploy](#undepapi)

### <a name="createapi"/> create

Import or create an API Proxy. If a bundle (zip) is supplied, it is imported else, it creates an empty proxy in the Org

```bash
apigeecli apis create -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--proxy -p` (required) API proxy bundle (zip)

### <a name="delapi"/> delete

Deletes an API proxy and all associated endpoints, policies, resources, and revisions. The API proxy must be undeployed before you can delete it.

```bash
apigeecli apis delete -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 

### <a name="depapi"/> deploy

Deploys a revision of an existing API proxy to an environment in an organization.

```bash
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

```bash
apigeecli apis fetch -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 
* `--rev -v` (required) API proxy revision
* `--ovr -r` (optional) Forces deployment of the new revision.

### <a name="getapi"/> get

Get API Proxy information

```bash
apigeecli apis get -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name 

### <a name="impapis"/> import

Upload a folder containing API proxy bundles

```bash
apigeecli apis import -o org -f /tmp
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--folder -f` (required) path containing API proxy bundles

### <a name="expapis"/> export

Export latest revisions of API proxy bundles from an org

```bash
apigeecli apis export -o org
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name


### <a name="listorgs"/> list

List APIs in an Apigee Org

```bash
apigeecli apis list -o org
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (optional) Apigee environment name
* `--rev -r` (optional) Include proxy revisions

If the environment name is passed, lists the deployed proxies in the environment
If the revision flag is enabled, lists all the revisions for the proxy

### <a name="listdeploy"/> listdeploy

Lists all deployments of an API proxy

```bash
apigeecli org listdeploy -o org -n name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--expand -x` (optional) Returns an expanded list of proxies for the organization.
* `--count -c` (optional) Number of apis to return

### <a name="trace"/> trace

Manage debug sessions/trace for API Proxy revisions

* [create](#crttrcapi)
* [get](#gettrcapi)
* [list](#listtrcapi)

#### <a name="crttrcapi"> create

Create a new trace/debug session

```bash
apigeecli apis trace create -o org -e env -n name -v 1 -f "name1=value1,name2=value2"
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) API proxy name
* `--rev -v` (required) API proxy revision
* `--filter -f` (optional) Trace filter; format is: name1=value1,name2=value2

#### <a name="getrcapi"> get

Get details for trace/debug session

```bash
apigeecli apis trace create -o org -e env -n name -v 1 -s uuid
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) API proxy name
* `--rev -v` (required) API proxy revision
* `--ses -s` (required) Trace session ID
* `--msg -m` (optional) Message ID

#### <a name="listtrcapi"> list

List all trace/debug session for a proxy revision in the last 24 hours

```bash
apigeecli apis trace create -o org -e env -n name -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) API proxy name
* `--rev -v` (required) API proxy revision

### <a name="undepapi"/> undeploy

Undeploys a revision of an existing API proxy to an environment in an organization.

```bash
apigeecli apis undeploy -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API proxy name
* `--rev -v` (required) API proxy revision

___

## <a name="apps"/> apps

Supported alias `applications`

* [create](#crtapp)
* [delete](#delapp)
* [genkey](#genkey)
* [get](#getapp)
* [list](#listapps)
* [import](#impapps)
* [export](#expapps)

### <a name ="crtapp"/> create

Create a developer app

```bash
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

```bash
apigeecli apps delete -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer App name

### <a name ="genkey"/> genkey

Create new developer KeyPairs Generates a new consumer key and consumer secret for the named developer app

```bash
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

```bash
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

```bash
apigeecli apps list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--expand -x` (optional) Returns an expanded list of apps for the organization.
* `--count -c` (optional) Number of app ids to return. Default is 10000

### <a name="impapps"/> import

Import Developer app entities into an org

```bash
apigeecli apps import -o org -f filepath.json
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--file -f` (required) A json file containing apps

A sample file format can be found [here](https://github.com/srinandan/apigeecli/blob/master/test/apps_config.json)

### <a name="expapps"/> export

Export Developer app entities from an org

```bash
apigeecli apps export -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

___

## <name ="cache"/> cache

* [delete](#delcache)
* [list](#listlist)

### <a name="delcache"/> delete

Delete a cache resource from an environment

```bash
apigeecli cache delete -o org -e env -n name
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Cache name

### <a name="listcache"/> list

List cache resources in an environment

```bash
apigeecli cache list -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

___

## <a name="devs"/> developers

Supported alias `developers`

* [create](#crtdev)
* [delete](#deldev)
* [get](#getdev)
* [list](#listdevs)
* [import](#impdev)
* [export](#expdev)

### <a name ="crtdev"/> create

Create a new App Developer

```bash
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

```bash
apigeecli devs get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer email


### <a name ="getdev"/> get

Get details of an App Developer

```bash
apigeecli devs get -o org -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) Developer email

### <a name="listdevs"/> list

List all App Developers in an org

```bash
apigeecli devs list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name="impdev"/> import

Import Developer entities into an org

```bash
apigeecli devs import -o org -f filepath.json
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--file -f` (required) A json file containing developers

### <a name="expdev"/> export

Import Developer entities into an org

```bash
apigeecli devs export -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

___

## <a name="env"/> envs

Supported alias `environments`

* [debugmask] (#debugmask)
* [get](#getenv)
* [iam](#iamenv)
* [list](#listenv)

### <a name="debugmask"/> debugmask

* [get](#getdebug)
* [set](#setdebug)

#### <a name="getdebug"/> get

Get debugmask settings for an environment

```bash
apigeecli envs debugmask get -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

#### <a name="setdebug"/> set

Set debugmask settings for an environment

```bash
apigeecli envs debugmask set -o org -e env -m "mask_settings_in_json"
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--mask -m` (required) Mask settings in JSON format

### <a name="iamenv"/> iam

* [get](#getiam)
* [set](#setiam)
* [test](#testiam)

#### <a name="getiam"/> get

Get IAM settings for an environment

```bash
apigeecli envs iam get -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

#### <a name="setiam"/> set

Set IAM settings for an environment

```bash
apigeecli envs iam set -o org -e env -p "iam_policy_in_json"
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--policy -p` (required) IAM policy in JSON format


### <a name="getenv"/> get

Get details of an environment

```bash
apigeecli envs get -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) environment name
* `--config -c` (optional) If set, returns environment configuration 

### <a name="listenv"/> list

List all environments in an org

```bash
apigeecli envs list -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

___

## <a name="flow"/> flowhooks

* [attach](#crtfh)
* [detach](#delfh)
* [get](#getfh)
* [list](#listfh)

### <a name ="crtfh"/> attach

Attach a Flowhook

```bash
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

```bash
apigeecli flowhooks detach -o org -e env -n PreFlow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Name of the flowhook


### <a name ="getfh"/> get

Get a details of a configured Flowhook

```bash
apigeecli flowhooks get -o org -e env -n PreFlow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Name of the flowhook

### <a name ="listfh"/> list

List of configured Flowhooks

```bash
apigeecli flowhooks list -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

___

## <a name="iam"/> iam

* [createall] (#createall)
* [createax] (#createax)
* [createcass] (#createcass)
* [createcass] (#createlogger)
* [createmetric] (#createmetric)
* [createsync] (#createsync)

### <a name ="createall"/> createall

Create a Google IAM Service Account with all the necessary roles to operate the hybrid runtime

```bash
apigeecli iam createall -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

### <a name ="createax"/> createax

Create a Google IAM Service Account for Apigee Analytics Agent

```bash
apigeecli iam createax -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

### <a name ="createcass"/> createcass

Create a Google IAM Service Account for Apigee hybrid Cassandra backup

```bash
apigeecli iam createcass -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

### <a name ="createlogger"/> createlogger

Create a Google IAM Service Account for StackDriver logger

```bash
apigeecli iam createlogger -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

### <a name="createmart"/> createmart

Create a Google IAM Service Account for Apigee MART

```bash
apigeecli iam createmart -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

### <a name="createmetrics"/> createmetrics

Create a Google IAM Service Account for StackDriver Metrics

```bash
apigeecli iam createmetrics -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

### <a name="createsync"/> createsync

Create a Google IAM Service Account for Apigee Sync

```bash
apigeecli iam createsync -p gcp-project-id -n service-account-name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID
* `--name -n` (required) Service Account Name

___

## <a name="keystores"/> keystores

* [create](#crtks)
* [delete](#delks)
* [get](#getks)
* [list](#listks)

### <a name="crtks"/> create

Create a new key store

```bash
apigeecli keystores create -o org -e env -n name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key Store name

### <a name ="delks"/> delete

Delete a key store

```bash
apigeecli keystores delete -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key store name

### <a name ="getks"/> get

Get details of a key store

```bash
apigeecli keystores get -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key store name

### <a name ="listks"/> list

List key stores in an environment

```bash
apigeecli keystores get -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

___

## <a name="keylaliases"/> keyaliases

* [create](#crtka)
* [delete](#delka)
* [get](#getka)
* [list](#listka)

### <a name="crtka"/> create

Create a new key aliases

```bash
apigeecli keyaliases create -o org -e env -n name -s alias
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key Store name
* `--alias -s` (required) Key Alias name

### <a name ="delka"/> delete

Delete a key aliases

```bash
apigeecli keyaliases delete -o org -e env -n name -s alias
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key store name
* `--alias -s` (required) Key Alias name

### <a name ="getka"/> get

Get details of a key aliases

```bash
apigeecli keyaliases get -o org -e env -n name -s alias
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key store name
* `--alias -s` (required) Key Alias name

### <a name ="listka"/> list

List key aliases in a key store

```bash
apigeecli keyaliases get -o org -e env -n name
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Key store name

___

## <a name="kvm"/> kvms

* [create](#crtkvm)
* [delete](#delkvm)
* [list](#listkvm)

### <a name ="crtkvm"/> create

Create a new KV Map

```bash
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

```bash
apigeecli kvms delete -o org -e env -n kvm1 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) KVM Map name

### <a name ="listkvm"/> list

List KVMs in an environment

```bash
apigeecli kvms list -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

___

## <a name="org"/> org

* [create](#createorg)
* [list](#listorgs)
* [get](#getorg)
* [setmart](#setmart)

### <a name="listorgs"/> list

List all the orgs available to the identity (service account)

```bash
apigeecli org list 
```

### <a name="getorg"/> get

Get org details for an Apigee Org

```bash
apigeecli org get -o org 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name="setmart"/> setmart

Configure MART endpoint for an Apigee Org

```bash
apigeecli org get -o org -m http://endpoint
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--mart -m` (required) MART endpoint
* `--whitelist -w` (optional) Enable/disable whitelisting of GCP IP for source connections to MART

___

## <a name="proejcts"/> projects

* [testiam](#testiamp)

### <a name="testiamp"/> testiam

Test IAM permissions for a project

```bash
apigeecli projects testiam -p gcp-project-id
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--prj -p` (required) GCP Project ID

___

## <a name="prods"/> products

Supported alias `prods`

* [create](#crtproduct)
* [delete](#delproduct)
* [get](#getproduct)
* [list](#listproducts)
* [import](#impproducts)

### <a name="crtproduct"/> create

Create an API product

```bash
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

```bash
apigeecli prods delete -o org -n name 
```
Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) API product name

### <a name ="getproduct"/> get

Get details of an API product

```bash
apigeecli prods list -o org -n name 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--name -n` (required) API product name

### <a name ="impproducts"/> import

Import API products from a configuration file

```bash
apigeecli prods list -o org -f file -c connections 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--file -f` (required) File containing API products
`--conn -c` (optional) Number of connections to establish; default is 4

A sample file format can be found [here](https://github.com/srinandan/apigeecli/blob/master/test/products_config.json)

### <a name ="expproducts"/> export

Export API products to a file

```bash
apigeecli prods export -o org -c connections 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--conn -c` (optional) Number of connections to establish; default is 4

___

## <a name="resources"/> resources

* [create] (#crtres)
* [delete] (#delres)
* [get] (#getres)
* [list] (#listres)

### <a name="crtres"> create

Create a resource to an environment

```bash
apigeecli resources create -o org -e env -n test.js -p jsc -r /tmp/test.js 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--env -e` (required) Apigee environment name
`--name -n` (required) Resource name
`--respath -r` (required) Resource path
`--type -p` (required) Resource type

### <a name="delres"> delete

Delete a resource from an environment

```bash
apigeecli resources create -o org -e env -n test.js -p jsc
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--env -e` (required) Apigee environment name
`--name -n` (required) Resource name
`--type -p` (required) Resource type

### <a name="getres"> get

Download a resource from an environment

```bash
apigeecli resources create -o org -e env -n test.js -p jsc
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--env -e` (required) Apigee environment name
`--name -n` (required) Resource name
`--type -p` (required) Resource type

### <a name="listres> list

List resources in an environment

```bash
apigeecli prods list -o org -e env 
```
Required parameters
The following parameters are required. See Common Reference for a list of additional parameters.

`--org -o` (required) Apigee organization name
`--env -e` (required) Apigee environment name
`--type -p` (optional) Filter by resource type

___

## <a name="sf"/> sharedflows

* [create](#createsf)
* [delete](#delsf)
* [deploy](#depsf)
* [export](#expsf)
* [fetch](#fetchsf)
* [import](#impsf)
* [get](#gettsf)
* [list](#listsf)
* [undeploy](#undepsf)

### <a name="createsf"/> create

Import or create a sharedflow. If a bundle (zip) is supplied, it is imported else, it creates an empty proxy in the Org

```bash
apigeecli apis create -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name
* `--proxy -p` (required) sharedflow bundle (zip)

### <a name="delsf"/> delete

Deletes a sharedflow and all policies, resources, and revisions. The sharedflow must be undeployed before you can delete it.

```bash
apigeecli sharedflows delete -o org -n proxy
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 

### <a name="depsf"/> deploy

Deploys a revision of an existing sharedflow to an environment in an organization.

```bash
apigeecli sharedflows deploy -o org -e env -n sharedflow1 -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 
* `--rev -v` (required) sharedflow revision
* `--ovr -r` (optional) Forces deployment of the new revision.

### <a name="expsf"/> export

Export sharedflows as bundles from an organization.

```bash
apigeecli sharedflows export -o org
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name

### <a name="fetchsf"/> fetch

Returns a zip-formatted proxy bundle of code and config files.

```bash
apigeecli apis fetch -o org -e env -n sharedflow -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 
* `--rev -v` (required) API proxy revision

### <a name="getsf"/> get

Get a sharedflow's details

```bash
apigeecli apis get -o org -e env -n sharedflow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name 

### <a name="impsf"/> import

Import sharedflows from dir to an organization.

```bash
apigeecli sharedflows import -o org -f /tmp
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--folder -f` (required) Folder containing sharedflow bundles

### <a name="listsf"/> list

List all sharedflows in an org

```bash
apigeecli apis get -o org -e env -n sharedflow
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (optional) Apigee environment name
* `--rev -r` (optional) Include shared flow revisions 

### <a name="undepsf"/> undeploy

Undeploys a revision of an existing API proxy to an environment in an organization.

```bash
apigeecli apis undeploy -o org -e env -n proxy -v 1
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--name -n` (required) sharedflow name
* `--rev -v` (required) sharedflow revision

___

## <a name="sync"/> sync

* [set](#setsync)
* [get](#getsync)


### <a name="listorgs"/> get

List all the orgs available to the identity (service account)

```bash
apigeecli sync get -o org
```

### <a name="listorgs"/> set

Set identity with access to control plane resources

```bash
apigeecli sync set -o org -i identity 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--ity -i` (required) IAM Identity

___

## <a name="target"/> targetservers

Supported alias `ts`

* [create](#crtts)
* [delete](#delts)
* [export](#expts)
* [import](#impts)
* [get](#getts)
* [list](#listts)

### <a name="createts"/> create

Create a new target server

```bash
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

```bash
apigeecli targetservers delete -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Target server name

### <a name ="expts"/> export

Export a target servers from an environment

```bash
apigeecli targetservers export -o org -e env
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

### <a name ="getts"/> get

Get details of a target server

```bash
apigeecli targetservers get -o org -e env -n name 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name
* `--name -n` (required) Target server name

### <a name ="listts"/> list

List target servers in an environment

```bash
apigeecli targetservers list -o org -e env 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--env -e` (required) Apigee environment name

___

## <a name="token"/> token

* [gen](#gentk)
* [cache](#cachetk)

### <a name ="gettk"/> gen

Generate a new access token

```bash
apigeecli token gen -a serviceaccount.json 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format

### <a name ="cachetk"/> cache

Caches a new access token. Writes the access token to $HOME/.access_token for use by subsequent calls

```bash
apigeecli token cache -a serviceaccount.json 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--account -a` (required) Service Account in json format

___