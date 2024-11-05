## apigeecli observations sources list

List Observation Sources

### Synopsis

List Observation Sources

```
apigeecli observations sources list [flags]
```

### Options

```
  -h, --help                help for list
      --page-size int       The maximum number of versions to return (default -1)
      --page-token string   A page token, received from a previous call
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
  -r, --region string    API Observation region name
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli observations sources](apigeecli_observations_sources.md)	 - Manage Observation sources for Shadow API Discovery

###### Auto generated by spf13/cobra on 30-Oct-2024