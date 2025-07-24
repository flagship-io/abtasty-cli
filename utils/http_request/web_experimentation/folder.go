package web_experimentation

import (
	"net/http"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/flagship-io/abtasty-cli/utils/http_request/common"
)

type FolderRequester struct {
	*common.ResourceRequest
}

func (f *FolderRequester) HTTPListFolder() ([]models.Folder, error) {
	return common.HTTPGetAllPagesWE[models.Folder](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/folders?")
}

func (f *FolderRequester) HTTPGetFolder(id int) (models.Folder, error) {
	return common.HTTPGetItem[models.Folder](utils.GetWebExperimentationHost() + "/v1/accounts/" + f.AccountID + "/folders/" + strconv.Itoa(id))
}

func (f *FolderRequester) HTTPCreateFolder(data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Folder](http.MethodPost, utils.GetWebExperimentationHost()+"/v1/accounts/"+f.AccountID+"/folders", data)
}
