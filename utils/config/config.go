package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/spf13/viper"
)

func CheckABTastyHomeDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()

	if _, err := os.Stat(homeDir + "/.abtasty/credentials/" + utils.FEATURE_EXPERIMENTATION); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(homeDir+"/.abtasty/credentials/"+utils.FEATURE_EXPERIMENTATION, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(homeDir + "/.abtasty/credentials/" + utils.WEB_EXPERIMENTATION); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(homeDir+"/.abtasty/credentials/"+utils.WEB_EXPERIMENTATION, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return homeDir, err
}

func CredentialPath(product, username string) (string, error) {
	homeDir, err := CheckABTastyHomeDirectory()
	if err != nil {
		return "", err
	}

	filepath, err := filepath.Abs(homeDir + "/.abtasty/credentials/" + product + "/" + username + ".yaml")
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func GetUsernames(product string) ([]string, error) {
	homeDir, err := CheckABTastyHomeDirectory()
	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`(?P<Username>[^/]+)\.yaml`)
	var fileNames []string

	f, err := os.Open(homeDir + "/.abtasty/credentials/" + product)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occurred: %s", err)
		return nil, err
	}

	files, err := f.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occurred: %s", err)
		return nil, err
	}

	for _, v := range files {
		match := r.FindStringSubmatch(v.Name())
		userName := r.SubexpIndex("Username")
		if len(match) == 0 {
			err := errors.New("Error: File not found")
			return nil, err
		}

		fileNames = append(fileNames, match[userName])
	}
	return fileNames, nil
}

func CreateAuthFile(product, username, clientId, clientSecret string, authenticationResponse models.TokenResponse) error {
	v := viper.New()
	filepath, err := CredentialPath(product, username)
	if err != nil {
		return err
	}

	v.Set("username", username)
	v.Set("client_id", clientId)
	v.Set("client_secret", clientSecret)
	v.Set("token", authenticationResponse.AccessToken)
	v.Set("refresh_token", authenticationResponse.RefreshToken)
	v.Set("scope", authenticationResponse.Scope)

	err = v.WriteConfigAs(filepath)
	if err != nil {
		return err
	}

	return nil
}

func ReadAuth(product, AuthName string) (*viper.Viper, error) {
	v := viper.New()
	configFilepath, err := CredentialPath(product, AuthName)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configFilepath); errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "error occurred: %v \n", err)
	}
	v.SetConfigFile(configFilepath)
	v.MergeInConfig()
	return v, nil
}

func SelectAuth(product, AuthName string) error {
	var v = viper.New()

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		return err
	}

	v.Set("current_used_credential", AuthName)

	err = v.WriteConfigAs(filepath)
	if err != nil {
		return err
	}

	return nil
}

func SetAccountID(product, accountID string) error {
	var v = viper.New()
	configFilepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		return err
	}

	v.SetConfigFile(configFilepath)
	v.MergeInConfig()

	v.Set("account_id", accountID)

	err = v.WriteConfigAs(configFilepath)
	if err != nil {
		return err
	}

	return nil
}

func SetAccountEnvID(product, accountEnvID string) error {
	var v = viper.New()
	configFilepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		return err
	}

	v.SetConfigFile(configFilepath)
	v.MergeInConfig()

	v.Set("account_environment_id", accountEnvID)

	err = v.WriteConfigAs(configFilepath)
	if err != nil {
		return err
	}

	return nil
}

func ReadCredentialsFromFile(AuthFile string) (*viper.Viper, error) {
	var v = viper.New()
	v.SetConfigFile(AuthFile)
	err := v.MergeInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func RewriteToken(product, AuthName string, authenticationResponse models.TokenResponse) error {
	v := viper.New()
	configFilepath, err := CredentialPath(product, AuthName)
	if err != nil {
		return err
	}

	v.SetConfigFile(configFilepath)

	v.MergeInConfig()
	v.Set("token", authenticationResponse.AccessToken)
	v.Set("refresh_token", authenticationResponse.RefreshToken)
	v.Set("scope", authenticationResponse.Scope)

	err = v.WriteConfigAs(configFilepath)
	if err != nil {
		return err
	}

	return nil
}
