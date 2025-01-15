package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/flagship-io/abtasty-cli/models"
	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/stretchr/testify/assert"
)

var (
	product      = "test_product"
	username     = "test_user"
	clientID     = "client_id"
	clientSecret = "client_secret"
	accessToken  = "access_token"
	refreshToken = "refresh_token"
	scope        = "scope"
	accountID    = "account_id"
	accountEnvID = "account_environment_id"
	identifier   = "identifier"
	email        = "email"
	workingDir   = "workingDir"
)

var authResponse = models.TokenResponse{
	AccessToken:  accessToken,
	RefreshToken: refreshToken,
	Scope:        scope,
}

type TestCampaignTargetingStruct struct {
	name       string
	workingDir string
	want       string
	code       string
	accountID  string
	campaignID string
	wantErr    bool
}

func TestMain(m *testing.M) {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	if _, err := os.Stat(homeDir + "/.abtasty/credentials/" + product); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(homeDir+"/.abtasty/credentials/"+product, os.ModePerm)
	}

	defer os.RemoveAll(currentDir + "/.abtasty")
	defer os.RemoveAll(homeDir + "/.abtasty/credentials/" + product)

	m.Run()
}

func TestCheckABTastyHomeDirectory(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	abtastyHome, err := CheckABTastyHomeDirectory()
	if err != nil {
		t.Errorf("CheckABTastyHomeDirectory() error = %v", err)
	}

	assert.Equal(t, homeDir, abtastyHome)
	assert.Equal(t, homeDir+"/.abtasty/credentials/"+utils.FEATURE_EXPERIMENTATION, abtastyHome+"/.abtasty/credentials/"+utils.FEATURE_EXPERIMENTATION)
	assert.Equal(t, homeDir+"/.abtasty/credentials/"+utils.WEB_EXPERIMENTATION, abtastyHome+"/.abtasty/credentials/"+utils.WEB_EXPERIMENTATION)
	assert.Equal(t, homeDir+"/.abtasty/credentials/"+product, abtastyHome+"/.abtasty/credentials/"+product)

}

func TestCredentialPath(t *testing.T) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	filepath, err := CredentialPath(product, username)
	if err != nil {
		t.Errorf("CredentialPath() error = %v", err)
	}

	expectedPath := homeDir + "/.abtasty/credentials/" + product + "/" + username + ".yaml"
	assert.Equal(t, expectedPath, filepath)

}

func TestGetUsernames(t *testing.T) {

	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Errorf("GetUsernames() error = %v", err)
	}

	fileNames, err := GetUsernames(product)
	if err != nil {
		t.Errorf("GetUsernames() error = %v", err)
	}

	if len(fileNames) != 1 || fileNames[0] != "test_user" {
		t.Errorf("GetUsernames() returned unexpected file names: %v", fileNames)
	}
}

func TestCreateAuthFile(t *testing.T) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	fileContent, err := os.ReadFile(homeDir + "/.abtasty/credentials/" + product + "/" + username + ".yaml")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := fmt.Sprintf(`client_id: %s
client_secret: %s
refresh_token: %s
scope: %s
token: %s
username: %s
`, clientID, clientSecret, refreshToken, scope, accessToken, username)

	assert.Equal(t, expectedContent, string(fileContent))

}

func TestReadAuth(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	v, err := ReadAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, v.GetString("client_id"), clientID)
	assert.Equal(t, v.GetString("client_secret"), clientSecret)
	assert.Equal(t, v.GetString("username"), username)
	assert.Equal(t, v.GetString("token"), authResponse.AccessToken)
	assert.Equal(t, v.GetString("refresh_token"), authResponse.RefreshToken)
}

func TestSelectAuth(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = SelectAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, string(yamlFile), "current_used_credential: test_user\n")
}

func TestSetAccountID(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = SelectAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	err = SetAccountID(product, accountID)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, string(yamlFile), "account_id: account_id\ncurrent_used_credential: test_user\n")
}

func TestSetAccountEnvironmentID(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = SelectAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	err = SetAccountEnvID(product, accountEnvID)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, string(yamlFile), "account_environment_id: account_environment_id\ncurrent_used_credential: test_user\n")
}

func TestReadCredentialsFromFile(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	filepath, err := CredentialPath(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	v, err := ReadCredentialsFromFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, v.GetString("client_id"), clientID)
	assert.Equal(t, v.GetString("client_secret"), clientSecret)
	assert.Equal(t, v.GetString("username"), username)
	assert.Equal(t, v.GetString("token"), authResponse.AccessToken)
	assert.Equal(t, v.GetString("refresh_token"), authResponse.RefreshToken)
}

func TestRewriteToken(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, models.TokenResponse{})
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = RewriteToken(product, username, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	v, err := ReadAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, v.GetString("client_id"), clientID)
	assert.Equal(t, v.GetString("client_secret"), clientSecret)
	assert.Equal(t, v.GetString("username"), username)
	assert.Equal(t, v.GetString("token"), authResponse.AccessToken)
	assert.Equal(t, v.GetString("refresh_token"), authResponse.RefreshToken)
}

func TestSetIdentifier(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = SelectAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	err = SetIdentifier(product, identifier)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, string(yamlFile), "current_used_credential: test_user\nidentifier: identifier\n")
}

func TestSetEmail(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = SelectAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	err = SetEmail(product, email)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, string(yamlFile), "current_used_credential: test_user\nemail: email\n")
}

func TestSetWorkingDir(t *testing.T) {
	err := CreateAuthFile(product, username, clientID, clientSecret, authResponse)
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	err = SelectAuth(product, username)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	err = SetWorkingDir(product, workingDir)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	filepath, err := CredentialPath(product, utils.HOME_CLI)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	assert.Equal(t, string(yamlFile), "current_used_credential: test_user\nworking_dir: workingDir\n")
}

func TestCampaignTargetingDirectory(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []TestCampaignCodeStruct{
		{
			name:       "ExistingDirectory",
			workingDir: currentDir,
			code:       "{\"url_scopes\":[{\"condition\":40,\"include\":true,\"value\":\"https://abtasty.com\"},{\"condition\":40,\"include\":false,\"value\":\"https://abtasty.com\"}]}\n", // Content of JSON file
			accountID:  "123456",
			campaignID: "100000",
			want:       currentDir + "/.abtasty/" + mockAccountID + "/" + mockCampaignID + "/targeting/targeting.json",
			wantErr:    false,
		},
		{
			name:       "NonExistingDirectory",
			workingDir: "/path/to/nonexistent/directory",
			code:       "{\"url_scopes\":[{\"condition\":40,\"include\":true,\"value\":\"https://abtasty.com\"},{\"condition\":40,\"include\":false,\"value\":\"https://abtasty.com\"}]}\n", // Content of JSON file
			accountID:  "123456",
			campaignID: "100000",
			want:       "",
			wantErr:    true,
		},
	}

	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				got, err := CampaignTargetingDirectory(tt.workingDir, tt.accountID, tt.campaignID, tt.code, true)
				if (err != nil) != tt.wantErr {
					t.Errorf("CampaignTargetingDirectory() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("CampaignTargetingDirectory() = %v, want %v", got, tt.want)
				}
			})

		}
	}
}

func TestHashString(t *testing.T) {
	s1 := "Hello, World!"
	s2 := "Hello, World!"
	s3 := "Different string"

	hash1 := HashString(s1)
	hash2 := HashString(s2)
	hash3 := HashString(s3)

	if !reflect.DeepEqual(hash1, hash2) {
		t.Errorf("hashString should produce the same hash for identical strings.\nGot:   %x\nWant:  %x", hash1, hash2)
	}

	if reflect.DeepEqual(hash1, hash3) {
		t.Errorf("hashString should produce different hashes for different strings.\nString1: %q\nString3: %q\nHash: %x", s1, s3, hash1)
	}
}

func TestHashFile(t *testing.T) {
	content := "File content to be hashed."
	expectedHash := HashString(content)

	tmpFile, err := os.CreateTemp("", "hash_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	fileHash, err := HashFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to hash file: %v", err)
	}

	if fileHash != expectedHash {
		t.Errorf("File hash and string hash do not match.\nFile Hash: %x\nExpected: %x", fileHash, expectedHash)
	}

	differentContent := "Some other content"
	differentHash := HashString(differentContent)

	tmpFile.Seek(0, 0)
	tmpFile.Truncate(0)
	if _, err := tmpFile.Write([]byte(differentContent)); err != nil {
		t.Fatalf("Failed to write different content to temp file: %v", err)
	}

	newFileHash, err := HashFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to hash file with new content: %v", err)
	}

	if newFileHash == expectedHash {
		t.Error("Expected different hash for changed file content, but got the same.")
	}

	if newFileHash != differentHash {
		t.Errorf("Expected hash of file content to match hash of string with same content.\nGot: %x\nWant: %x", newFileHash, differentHash)
	}
}
