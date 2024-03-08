package feature_experimentation

import (
	"net/http"

	"github.com/flagship-io/flagship/models"
	"github.com/flagship-io/flagship/utils"
	"github.com/spf13/viper"
)

func HTTPListTargetingKey() ([]models.TargetingKey, error) {
	return HTTPGetAllPages[models.TargetingKey](utils.GetHost() + "/v1/accounts/" + viper.GetString("account_id") + "/targeting_keys")
}

func HTTPGetTargetingKey(id string) (models.TargetingKey, error) {
	return HTTPGetItem[models.TargetingKey](utils.GetHost() + "/v1/accounts/" + viper.GetString("account_id") + "/targeting_keys/" + id)
}

func HTTPCreateTargetingKey(data string) ([]byte, error) {
	return HTTPRequest(http.MethodPost, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys", []byte(data))
}

func HTTPEditTargetingKey(id, data string) ([]byte, error) {
	return HTTPRequest(http.MethodPatch, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys/"+id, []byte(data))
}

func HTTPDeleteTargetingKey(id string) error {
	_, err := HTTPRequest(http.MethodDelete, utils.GetHost()+"/v1/accounts/"+viper.GetString("account_id")+"/targeting_keys/"+id, nil)
	return err
}
