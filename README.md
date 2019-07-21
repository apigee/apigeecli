# apigeecli

This is a tool to interact with [Apigee APIs](https://apigee.googleapis.com). The tool let you manage (get, list) environments, proxies, etc.

# Installation

`apigeecli` is a binary and you can download the appropriate one for your platform from [here](https://github.com/srinandan/apigeepi/releases)

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

`apigeecli` can use the service account directly and obtain an access token. However, the user can also provide an access token. 

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
apigeecli -t $TOKEN
```

Pass the service account

```
apigeecli -a orgadmin.json
```

## Access Token Caching

`apigeecli` caches the OAuth Access token for subsequent calls (until the token expires). The access token is stored in `$HOME/.access_token`. This path must be readable/writeable by the `apigeecli` process. 

<details><summary> Example 1 </summary>
<p>
First command

```
apigeecli orgs list -a orgadmin.json
```

Subsequent command (no token or service account)

```
apigeecli orgs get -o org
```
</p>
</details>


<details><summary> Example 2 </summary>
<p>


First command

```
apigeecli orgs list -t $TOKEN
```

Subsequent command (no token or service account)

```
apigeecli orgs get -o org
```
</p>
</details>

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

---

## <a name="apis"/> apis

* [create](#createapi)
* [delete](#delapi)
* [deploy](#depapi)
* [fetch](#fetchapi)
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
* `--revision -v` (required) API proxy revision
* `--ovr -r` (optional) Forces deployment of the new revision.

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
* `--revision -v` (required) API proxy revision

---

## <a name="apps"/> apps

Supported alias `applications`

* [create](#crtapp)
* [delete](#delapp)
* [genkey](#genkey)
* [get](#getapp)
* [list](#listapps)

### <a name ="crtapp"/> create

Create a developer app

```
apigeecli apps create -o org -n name -e test,prod -p proxy1 
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
* `--name -n` (required) Developer App name

### <a name="listapps"/> list

List all developer apps in an org

```
apigeecli apps list -o org 
```

Parameters
The following parameters are supported. See Common Reference for a list of additional parameters.

* `--org -o` (required) Apigee organization name
* `--expand -x` (optional) Returns an expanded list of apps for the organization.
* `--count -c` (optional) Number of app ids to return.

---