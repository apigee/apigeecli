# Unit Tests for apigeecli

## Environment Variables

Set the following environment variables:

* `APIGEE_ORG`
* `APIGEE_ENV`
* `APIGEECLI_PATH` (path to the apigeecli folder)

If the org was created with DRZ options, then also set:

* `APIGEE_REGION`

For tests involving Apigee API Portal, set:

* `APIGEE_SITEID`

For tests involving Apigee API Registry, set:

* `APIGEE_REGION`

## Running the tests

Run `go test internal/client/<package>` or Run `go test internal/client/<package> -cover` to include code coverage
