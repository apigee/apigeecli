## apigeecli registry apis versions specs artifacts create

Create an artifact for an API version

### Synopsis

Create an artifact for an API version

```
apigeecli registry apis versions specs artifacts create [flags]
```

### Options

```
      --annotations stringToString   Annotations attach non-identifying metadata to resources (default [])
      --api-name string              Name of the API
      --api-version string           API Version ID
  -f, --file string                  Path to a file containing Artifact Contents
  -h, --help                         help for create
  -i, --id string                    Apigee Registry Artifact ID
      --labels stringToString        Labels attach identifying metadata to resources (default [])
  -n, --name string                  Name of the Artifact
      --spec-name string             API Version Spec name
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

* [apigeecli registry apis versions specs artifacts](apigeecli_registry_apis_versions_specs_artifacts.md)	 - Manage artifacts for an API version's spec

###### Auto generated by spf13/cobra on 6-Mar-2024