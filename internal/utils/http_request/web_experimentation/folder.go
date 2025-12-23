package web_experimentation

import (
	"fmt"
	"net/http"
	"strconv"

	models "github.com/flagship-io/abtasty-cli/internal/models/web_experimentation"
	"github.com/flagship-io/abtasty-cli/internal/utils"
	"github.com/flagship-io/abtasty-cli/internal/utils/http_request/common"
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

func (f *FolderRequester) HTTPEditFolder(id int, data []byte) ([]byte, error) {
	return common.HTTPRequest[models.Folder](http.MethodPatch, utils.GetWebExperimentationHost()+"/v1/accounts/"+f.AccountID+"/folders/"+strconv.Itoa(id), data)
}

func (f *FolderRequester) HTTPDeleteFolder(id int) (string, error) {
	_, err := common.HTTPRequest[models.Folder](http.MethodDelete, utils.GetWebExperimentationHost()+"/v1/accounts/"+f.AccountID+"/folders/"+strconv.Itoa(id), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Folder %d deleted", id), nil
}
