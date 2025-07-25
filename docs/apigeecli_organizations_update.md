## apigeecli organizations update

Update settings of an Apigee Org

### Synopsis

Update settings of an Apigee Org

```
apigeecli organizations update [flags]
```

### Options

```
  -d, --desc string   Apigee org description
  -h, --help          help for update
  -n, --net string    Authorized network; if using a shared VPC format is projects/{host-project-id}/{location}/networks/{network-name} (default "default")
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
      --print-output     Control printing of info log statements (default true)
  -r, --region string    Apigee control plane region name; default is https://apigee.googleapis.com
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli organizations](apigeecli_organizations.md)	 - Manage Apigee Orgs

###### Auto generated by spf13/cobra on 1-Jul-2025
