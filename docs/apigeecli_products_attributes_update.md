## apigeecli products attributes update

Update an attribute of an API product

### Synopsis

Update an attribute of an API product

```
apigeecli products attributes update [flags]
```

### Options

```
  -k, --attr string    API Product attribute name
  -h, --help           help for update
  -v, --value string   API Product attribute value
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --api api          Sets the control plane API. Must be one of prod, autopush or staging; default is prod
      --default-token    Use Google default application credentials access token
      --disable-check    Disable check for newer versions
      --metadata-token   Metadata OAuth2 access token
  -n, --name string      API Product name
      --no-output        Disable printing all statements to stdout
      --no-warnings      Disable printing warnings to stderr
  -o, --org string       Apigee organization name
      --print-output     Control printing of info log statements (default true)
  -r, --region string    Apigee control plane region name; default is https://apigee.googleapis.com
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli products attributes](apigeecli_products_attributes.md)	 - Manage API Product Attributes

###### Auto generated by spf13/cobra on 1-Jul-2025
