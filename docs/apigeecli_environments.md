## apigeecli environments

Manage Apigee environments

### Synopsis

Manage Apigee environments

### Options

```
  -h, --help            help for environments
  -o, --org string      Apigee organization name
  -r, --region string   Apigee control plane region name; default is https://apigee.googleapis.com
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
  -t, --token string     Google OAuth Token
```

### SEE ALSO

* [apigeecli](apigeecli.md)	 - Utility to work with Apigee APIs.
* [apigeecli environments archives](apigeecli_environments_archives.md)	 - Manage archive deployments for the environment
* [apigeecli environments ax-obfuscation](apigeecli_environments_ax-obfuscation.md)	 - Obfuscate analytics fields
* [apigeecli environments create](apigeecli_environments_create.md)	 - Create a new environment
* [apigeecli environments debugmask](apigeecli_environments_debugmask.md)	 - Manage debugmasks for the environment
* [apigeecli environments delete](apigeecli_environments_delete.md)	 - Delete an environment
* [apigeecli environments deployments](apigeecli_environments_deployments.md)	 - Manage deployments for the environment
* [apigeecli environments export](apigeecli_environments_export.md)	 - Export environment details to a file
* [apigeecli environments get](apigeecli_environments_get.md)	 - Get properties of an environment
* [apigeecli environments iam](apigeecli_environments_iam.md)	 - Manage IAM permissions for the environment
* [apigeecli environments import](apigeecli_environments_import.md)	 - Import a file containing environment details
* [apigeecli environments list](apigeecli_environments_list.md)	 - List environments in an Apigee Org
* [apigeecli environments secactions](apigeecli_environments_secactions.md)	 - Manage SecurityActions for Apigee Advanced Security
* [apigeecli environments secactionscfg](apigeecli_environments_secactionscfg.md)	 - Manage SecurityActionsConfig for Apigee Advanced Security
* [apigeecli environments secincidents](apigeecli_environments_secincidents.md)	 - View SecurityIncidents from Apigee Advanced Security
* [apigeecli environments set](apigeecli_environments_set.md)	 - Set environment property
* [apigeecli environments traceconfig](apigeecli_environments_traceconfig.md)	 - Manage Distributed Trace config for the environment

###### Auto generated by spf13/cobra on 1-Jul-2025
