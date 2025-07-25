## apigeecli securityprofiles

Manage Adv API Security Profiles

### Synopsis

Manage Adv API Security Profiles

### Options

```
  -h, --help            help for securityprofiles
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
* [apigeecli securityprofiles attach](apigeecli_securityprofiles_attach.md)	 - Attach a security profile to an environment
* [apigeecli securityprofiles compute](apigeecli_securityprofiles_compute.md)	 - Calculates scores for requested time range
* [apigeecli securityprofiles create](apigeecli_securityprofiles_create.md)	 - Create a new Security Profile
* [apigeecli securityprofiles delete](apigeecli_securityprofiles_delete.md)	 - Deletes a security profile
* [apigeecli securityprofiles detach](apigeecli_securityprofiles_detach.md)	 - Detach a security profile from an environment
* [apigeecli securityprofiles export](apigeecli_securityprofiles_export.md)	 - Export Security Profiles to a file
* [apigeecli securityprofiles get](apigeecli_securityprofiles_get.md)	 - Returns a security profile by name
* [apigeecli securityprofiles import](apigeecli_securityprofiles_import.md)	 - Import a folder containing Security Profiles
* [apigeecli securityprofiles list](apigeecli_securityprofiles_list.md)	 - Returns the security profiles in the org
* [apigeecli securityprofiles listrevisions](apigeecli_securityprofiles_listrevisions.md)	 - Returns the revisions of a security profile
* [apigeecli securityprofiles update](apigeecli_securityprofiles_update.md)	 - Update an existing Security Profile

###### Auto generated by spf13/cobra on 1-Jul-2025
