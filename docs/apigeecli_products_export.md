## apigeecli products export

Export API products to a file

### Synopsis

Export API products to a file

```
apigeecli products export [flags]
```

### Options

```
  -c, --conn int       Number of connections (default 4)
  -h, --help           help for export
      --space string   Apigee Space associated to
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

* [apigeecli products](apigeecli_products.md)	 - Manage Apigee API products

###### Auto generated by spf13/cobra on 1-Jul-2025
