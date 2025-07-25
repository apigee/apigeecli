## apigeecli developers update

Update an Apigee developer configuration

### Synopsis

Update an Apigee developer configuration

```
apigeecli developers update [flags]
```

### Options

```
      --attrs stringToString     Custom attributes (default [])
  -n, --email string             The developer's email
  -f, --first string             The first name of the developer
  -h, --help                     help for update
  -s, --last string              The last name of the developer
      --status developerStatus   must be active or inactive
  -u, --user string              The username of the developer
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

* [apigeecli developers](apigeecli_developers.md)	 - Manage Apigee App Developers

###### Auto generated by spf13/cobra on 1-Jul-2025
