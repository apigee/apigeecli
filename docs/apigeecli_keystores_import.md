## apigeecli keystores import

Import a file containing keystores

### Synopsis

Import a file containing keystores

```
apigeecli keystores import [flags]
```

### Examples

```
Import a file containing keystores:
apigeecli keystores import -f samples/keystores.json  -e $env
```

### Options

```
  -c, --conn int      Number of connections (default 4)
  -f, --file string   File containing keystores
  -h, --help          help for import
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --api api          Sets the control plane API. Must be one of prod, autopush or staging; default is prod
      --default-token    Use Google default application credentials access token
      --disable-check    Disable check for newer versions
  -e, --env string       Apigee environment name
      --metadata-token   Metadata OAuth2 access token
      --no-output        Disable printing all statements to stdout
      --no-warnings      Disable printing warnings to stderr
  -o, --org string       Apigee organization name
      --print-output     Control printing of info log statements (default true)
  -r, --region string    Apigee control plane region name; default is https://apigee.googleapis.com
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli keystores](apigeecli_keystores.md)	 - Manage Key Stores

###### Auto generated by spf13/cobra on 1-Jul-2025
