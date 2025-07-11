## apigeecli products update

Update an API product

### Synopsis

Update an API product

```
apigeecli products update [flags]
```

### Options

```
  -f, --approval string              Approval type
      --attrs stringToString         Custom attributes (default [])
  -d, --desc string                  Description for the API Product
  -m, --display-name string          Display Name of the API Product
  -e, --envs stringArray             Environments to enable
      --gqlopgrp string              File containing GraphQL Operation Group JSON. See samples for how to create the file
      --grpcopgrp string             File containing gRPC Operation Group JSON. See samples for how to create the file
  -h, --help                         help for update
  -i, --interval string              Quota Interval
  -n, --name string                  Name of the API Product
      --opgrp string                 File containing Operation Group JSON. See samples for how to create the file
  -p, --proxies stringArray          API Proxies in product
  -q, --quota string                 Quota Amount
      --quota-counter-scope string   Scope of the quota decides how the quota counter gets applied; can be PROXY or OPERATION
  -s, --scopes stringArray           OAuth scopes
      --space string                 Associated Apigee Space. Pass this if the API Product being updated is part of a space
  -u, --unit string                  Quota Unit
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

* [apigeecli products](apigeecli_products.md)	 - Manage Apigee API products

###### Auto generated by spf13/cobra on 1-Jul-2025
