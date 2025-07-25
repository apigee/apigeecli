## apigeecli envgroups

Manage Apigee environment groups

### Synopsis

Manage Apigee environment groups

### Options

```
  -h, --help            help for envgroups
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
* [apigeecli envgroups attach](apigeecli_envgroups_attach.md)	 - Attach an env to an Environment Group
* [apigeecli envgroups create](apigeecli_envgroups_create.md)	 - Create an Environment Group
* [apigeecli envgroups delete](apigeecli_envgroups_delete.md)	 - Deletes an Environment Group
* [apigeecli envgroups detach](apigeecli_envgroups_detach.md)	 - Detach an env from an Environment Group
* [apigeecli envgroups get](apigeecli_envgroups_get.md)	 - Gets an Environment Group
* [apigeecli envgroups import](apigeecli_envgroups_import.md)	 - Import a file containing environment group definitions
* [apigeecli envgroups list](apigeecli_envgroups_list.md)	 - Returns a list of environment groups
* [apigeecli envgroups listattach](apigeecli_envgroups_listattach.md)	 - List attachments of an Environment Group
* [apigeecli envgroups update](apigeecli_envgroups_update.md)	 - Update hostnames in an Environment Group

###### Auto generated by spf13/cobra on 1-Jul-2025
