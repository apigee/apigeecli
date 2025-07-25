## apigeecli organizations deployments get

Get deployments for an Apigee org

### Synopsis

Get deployments for an Apigee org

```
apigeecli organizations deployments get [flags]
```

### Options

```
      --all           Return all deployments
  -h, --help          help for get
  -s, --sharedflows   Return sharedflow deployments
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

* [apigeecli organizations deployments](apigeecli_organizations_deployments.md)	 - Manage deployments in an Apigee org

###### Auto generated by spf13/cobra on 1-Jul-2025
