package handler

import (
	"fmt"
	"os"

	featureResource "github.com/flagship-io/abtasty-cli/cmd/feature-experimentation/resource"
	webResource "github.com/flagship-io/abtasty-cli/cmd/web-experimentation/resource"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/spf13/viper"
)

// LoadResources loads resources from a resource loader JSON content, processes them, and returns the results.
func LoadWebExperimentationResources(resourceLoaderContent, accountID, token string, dryRun bool) (string, error) {
	v := viper.New()
	var requestConfig = common.RequestConfig{Product: utils.WEB_EXPERIMENTATION}

	v.Set("username", "no-username")
	v.Set("account_id", accountID)
	v.Set("token", token)

	v.Unmarshal(&requestConfig)
	common.Init(requestConfig)

	r := &http_request.ResourceRequester
	r.Init(&requestConfig)

	results, err := webResource.LoadResources(os.Stdout, resourceLoaderContent, "", dryRun)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return results, nil
}

func LoadFeatureExperimentationResources(resourceLoaderContent, accountID, accountEnvID, token string, dryRun bool) (string, error) {
	v := viper.New()
	var requestConfig = common.RequestConfig{Product: utils.FEATURE_EXPERIMENTATION}

	v.Set("username", "no-username")
	v.Set("account_id", accountID)
	v.Set("account_environment_id", accountEnvID)
	v.Set("token", token)

	v.Unmarshal(&requestConfig)
	common.Init(requestConfig)

	r := &http_request.ResourceRequester
	r.Init(&requestConfig)

	results, err := featureResource.LoadResources(os.Stdout, resourceLoaderContent, "", dryRun)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return results, nil
}
