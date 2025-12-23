package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/flagship-io/abtasty-cli/internal/models"
	"github.com/flagship-io/abtasty-cli/internal/utils"
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
		if viper.GetString("output_format") != "json" {
			fmt.Fprintf(os.Stderr, "error occurred: %v \n", err)
		}

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

func SetIdentifier(product, identifier string) error {
	var v = viper.New()
	configFilepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		return err
	}

	v.SetConfigFile(configFilepath)
	v.MergeInConfig()

	v.Set("identifier", identifier)

	err = v.WriteConfigAs(configFilepath)
	if err != nil {
		return err
	}

	return nil
}

func SetEmail(product string, email string) error {
	var v = viper.New()
	configFilepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		return err
	}

	v.SetConfigFile(configFilepath)
	v.MergeInConfig()

	v.Set("email", email)

	err = v.WriteConfigAs(configFilepath)
	if err != nil {
		return err
	}

	return nil
}

func SetWorkingDir(product, path string) error {
	var v = viper.New()
	configFilepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		return err
	}

	v.SetConfigFile(configFilepath)
	v.MergeInConfig()

	v.Set("working_dir", path)

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

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return fmt.Errorf("read config %s: %w", configFilepath, err)
		}
	}

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

func CampaignTargetingDirectory(workingDir, accountID, campaignID, code string, override bool) (string, error) {
	gcWorkingDir, err := CheckGlobalCodeDirectory(workingDir)
	if err != nil {
		return "", err
	}

	accountCodeDir := gcWorkingDir + "/" + accountID
	campaignCodeDir := accountCodeDir + "/" + campaignID
	targetingCodeDir := campaignCodeDir + "/targeting"

	err = os.MkdirAll(targetingCodeDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	jsonFilePath := targetingCodeDir + "/targeting.json"
	if _, err := os.Stat(jsonFilePath); err == nil {
		if !override {
			fmt.Fprintln(os.Stderr, "File already exists: "+jsonFilePath)
			return jsonFilePath, nil
		}
	}

	err = os.WriteFile(jsonFilePath, []byte(code), os.ModePerm)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout, "File created: "+jsonFilePath)
	return jsonFilePath, nil
}
