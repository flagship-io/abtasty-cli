package handler

import (
	"fmt"
	"os"

	"github.com/flagship-io/abtasty-cli/cmd/web-experimentation/resource"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
	"github.com/spf13/viper"
)

// LoadResources loads resources from a resource loader JSON content, processes them, and returns the results.
func LoadResources(resourceFilePath, accountID, token, inputRefFile, outputFile string, dryRun bool) (string, error) {
	v := viper.New()
	var requestConfig = common.RequestConfig{Product: utils.WEB_EXPERIMENTATION}

	v.Set("username", "no-username")
	v.Set("account_id", accountID)
	v.Set("token", token)

	v.Unmarshal(&requestConfig)
	common.Init(requestConfig)

	r := &http_request.ResourceRequester
	r.Init(&requestConfig)

	err := resource.LoadResources(os.Stdout, resourceFilePath, inputRefFile, "", outputFile)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return "Resources loaded successfully", nil
}
