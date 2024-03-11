## apigeecli registry apis deployments list

List deployments for an API in Apigee Registry

### Synopsis

List deployments for an API in Apigee Registry

```
apigeecli registry apis deployments list [flags]
```

### Options

```
      --api-name string     API Name
      --filter string       An expression that can be used to filter the list
  -h, --help                help for list
  -n, --name string         Deployment Name
      --order-by string     A comma-separated list of fields to be sorted; ex: foo desc
      --page-token string   A page token, received from a previous list call
```

### Options inherited from parent commands

```
  -a, --account string     Path Service Account private key in JSON
      --default-token      Use Google default application credentials access token
      --disable-check      Disable check for newer versions
      --metadata-token     Metadata OAuth2 access token
      --no-output          Disable printing all statements to stdout
      --no-warnings        Disable printing warnings to stderr
  -o, --org string         Apigee organization name
      --print-output       Control printing of info log statements (default true)
  -p, --projectID string   Apigee Registry Project ID; Use if Apigee Orgniazation not provisioned in this project
      --region string      Region where Apigee Registry is provisioned
  -t, --token string       Google OAuth Token
```

### SEE ALSO

* [apigeecli registry apis deployments](apigeecli_registry_apis_deployments.md)	 - Manage API Deployments in Apigee Registry

###### Auto generated by spf13/cobra on 6-Mar-2024