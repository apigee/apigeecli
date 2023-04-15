package bundlegen

import (
	apiproxy "internal/bundlegen/apiproxydef"
	"internal/bundlegen/proxies"
)

func GenerateIntegrationAPIProxy(name string,
	integration string,
	apitrigger string,
) (err error) {
	apiproxy.SetDisplayName(name)
	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	apiproxy.AddProxyEndpoint("default")
	apiproxy.AddIntegrationEndpoint("default")
	apiproxy.SetBasePath("/" + apitrigger)

	proxies.NewProxyEndpoint("/"+apitrigger, false)

	proxies.AddStepToPreFlowRequest("set-integration-request")
	apiproxy.AddPolicy("set-integration-request")

	return nil
}
