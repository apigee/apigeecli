## apigeecli apihub apis export

Export API, versions and specifications

### Synopsis

Export API, versions and specifications

```
apigeecli apihub apis export [flags]
```

### Options

```
      --api-id string   API ID
      --folder string   Folder to export the API details
  -h, --help            help for export
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
  -r, --region string    API Hub region name
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli apihub apis](apigeecli_apihub_apis.md)	 - Manage Apigee API Hub APIs

###### Auto generated by spf13/cobra on 23-Jul-2024