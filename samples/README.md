# apigeecli command Samples

The following table contains some examples of apigeecli.

Set up apigeecli with preferences: `apigeecli prefs set -o $org`

| Operations | Command |
|---|---|
| apicategories | `apigeecli apicategories import --file samples/apicategories.json --siteid $siteId  --default-token`|
| apis | `apigeecli apis import -f samples/apis  --default-token` |
| apis | `apigeecli apis create oas -n petstore --oas-base-folderpath=./samples --oas-name=petstore.yaml --add-cors=true --google-idtoken-aud-literal=https://sample.run.app  --default-token` |
| apis | `apigeecli apis create oas -n petstore -f ./samples/petstore.yaml --add-cors=true --env=$env --wait=true --default-token` |
| appgroups | `apigeecli appgroups import -f samples/appgroups.json --default-token` |
| datacollectors | `apigeecli datacollectors import -f samples/datacollectors.json --default-token`  |
| developers | `apigeecli developers import -f samples/developers.json`  |
| kvms | `apigeecli kvms import -f samples/kvms`  |
| products | `apigeecli products import -f samples/apiproduct-legacy.json  --default-token` |
| products | `apigeecli products import -f samples/apiproduct-gqlgroup.json  --default-token` |
| products | `apigeecli products import -f samples/apiproduct-op-group.json  --default-token` |
| products | `apigeecli products create --name $product_name --display-name $product_name --opgrp $ops_file --envs $env --approval auto --attrs access=public  --default-token` |
| products | `apigeecli products create --name $product_name --display-name $product_name --opgrp $ops_file --envs $env --approval auto --attrs access=public --quota 100 --interval 1 --unit minute  --default-token` |
| sharedflows | `apigeecli sharedflows import -f samples/sharedflows` |
| targetservers | `apigeecli targetservers import -f samples/targetservers.json  -e $env` |
| keystores | `apigeecli keystores import -f samples/keystores.json  -e $env` |
| references | `apigeecli references import -f samples/references.json  -e $env` |
| apps | `apigeecli apps import -f samples/references.json -d samples/developers.json` |
| apidocs | `apigeecli apidocs import -f samples/apidocs  -s $siteId`  |


NOTE: This file is auto-generated during a release. Do not modify.