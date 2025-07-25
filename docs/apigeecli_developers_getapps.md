## apigeecli developers getapps

Returns the apps owned by a developer by email address

### Synopsis

Returns the apps owned by a developer by email address

```
apigeecli developers getapps [flags]
```

### Options

```
  -x, --expand        expand app details
  -h, --help          help for getapps
  -n, --name string   email of the developer
  -o, --org string    Apigee organization name
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

* [apigeecli developers](apigeecli_developers.md)	 - Manage Apigee App Developers

###### Auto generated by spf13/cobra on 1-Jul-2025
