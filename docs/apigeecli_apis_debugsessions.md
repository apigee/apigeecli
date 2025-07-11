## apigeecli apis debugsessions

Manage debusessions of Apigee API proxies

### Synopsis

Manage debusessions of Apigee API proxy revisions deployed in an environment

### Options

```
  -e, --env string   Apigee environment name
  -h, --help         help for debugsessions
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --api api          Sets the control plane API. Must be one of prod, autopush or staging; default is prod
      --default-token    Use Google default application credentials access token
      --disable-check    Disable check for newer versions
      --metadata-token   Metadata OAuth2 access token
      --no-output        Disable printing all statements to stdout
      --no-warnings      Disable printing warnings to stderr
  -o, --org string       Apigee organization name
      --print-output     Control printing of info log statements (default true)
  -r, --region string    Apigee control plane region name; default is https://apigee.googleapis.com
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli apis](apigeecli_apis.md)	 - Manage Apigee API proxies in an org
* [apigeecli apis debugsessions create](apigeecli_apis_debugsessions_create.md)	 - Create a new debug session for an API proxy
* [apigeecli apis debugsessions get](apigeecli_apis_debugsessions_get.md)	 - Get a debug session for an API proxy revision
* [apigeecli apis debugsessions list](apigeecli_apis_debugsessions_list.md)	 - List all debug sessions for an API proxy revision

###### Auto generated by spf13/cobra on 1-Jul-2025
