## apigeecli apicategories get

Gets an API Category by ID or name

### Synopsis

Gets an API Category by ID or name

```
apigeecli apicategories get [flags]
```

### Options

```
  -h, --help          help for get
  -i, --id string     API Category ID
  -n, --name string   API Catalog Name
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
  -s, --siteid string    Name or siteid of the portal
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli apicategories](apigeecli_apicategories.md)	 - Manage Apigee API categories that are tagged on catalog items

###### Auto generated by spf13/cobra on 1-Jul-2025
