package config

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	name           string
	workingDir     string
	want           string
	code           string
	accountID      string
	campaignID     string
	variationID    string
	modificationID string
	selector       string
	wantErr        bool
}

var (
	mockAccountID      = "123456"
	mockCampaignID     = "100000"
	mockVariationID    = "200000"
	mockModificationID = "300000"
	mockSelector       = "document.querySelector('main')"
)

func TestCheckWorkingDirectory(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	tests := []TestStruct{
		{
			name:       "ExistingDirectory",
			workingDir: currentDir,
			want:       currentDir,
			wantErr:    false,
		},
		{
			name:       "NonExistingDirectory",
			workingDir: "/path/to/nonexistent/directory",
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckWorkingDirectory(tt.workingDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckWorkingDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckWorkingDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckGlobalCodeDirectory(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	tests := []TestStruct{
		{
			name:       "ExistingDirectory",
			workingDir: currentDir,
			want:       currentDir + "/.abtasty",
			wantErr:    false,
		},
		{
			name:       "NonExistingDirectory",
			workingDir: "/path/to/nonexistent/directory",
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckGlobalCodeDirectory(tt.workingDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckGlobalCodeDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckGlobalCodeDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountGlobalCodeDirectory(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []TestStruct{
		{
			name:       "ExistingDirectory",
			workingDir: currentDir,
			code:       "console.log('Hello, World!')", // Content of JavaScript file
			accountID:  mockAccountID,
			want:       currentDir + "/.abtasty/" + mockAccountID + "/accountGlobalCode.js",
			wantErr:    false,
		},
		{
			name:       "NonExistingDirectory",
			workingDir: "/path/to/nonexistent/directory",
			code:       "console.log('Hello, World!')", // Content of JavaScript file
			accountID:  mockAccountID,
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AccountGlobalCodeDirectory(tt.workingDir, tt.accountID, tt.code, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountGlobalCodeDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("AccountGlobalCodeDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCampaignGlobalCodeDirectory(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []TestStruct{
		{
			name:       "ExistingDirectory",
			workingDir: currentDir,
			code:       "console.log('Hello, World!')", // Content of JavaScript file
			accountID:  "123456",
			campaignID: "100000",
			want:       currentDir + "/.abtasty/" + mockAccountID + "/" + mockCampaignID + "/campaignGlobalCode.js",
			wantErr:    false,
		},
		{
			name:       "NonExistingDirectory",
			workingDir: "/path/to/nonexistent/directory",
			code:       "console.log('Hello, World!')", // Content of JavaScript file
			accountID:  "123456",
			campaignID: "100000",
			want:       "",
			wantErr:    true,
		},
	}

	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				got, err := CampaignGlobalCodeDirectory(tt.workingDir, tt.accountID, tt.campaignID, tt.code, true)
				if (err != nil) != tt.wantErr {
					t.Errorf("CampaignGlobalCodeDirectory() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("CampaignGlobalCodeDirectory() = %v, want %v", got, tt.want)
				}
			})

		}
	}
}

func TestVariationGlobalCodeDirectoryJS(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []TestStruct{
		{
			name:        "ExistingDirectory",
			workingDir:  currentDir,
			code:        "console.log('Hello, World!')", // Content of JavaScript file
			accountID:   mockAccountID,
			campaignID:  mockCampaignID,
			variationID: mockVariationID,
			want:        currentDir + "/.abtasty/" + mockAccountID + "/" + mockCampaignID + "/" + mockVariationID + "/variationGlobalCode.js",
			wantErr:     false,
		},
		{
			name:        "NonExistingDirectory",
			workingDir:  "/path/to/nonexistent/directory",
			code:        "console.log('Hello, World!')", // Content of JavaScript file
			accountID:   mockAccountID,
			campaignID:  mockCampaignID,
			variationID: mockVariationID,
			want:        "",
			wantErr:     true,
		},
	}

	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				got, err := VariationGlobalCodeDirectoryJS(tt.workingDir, tt.accountID, tt.campaignID, tt.variationID, tt.code, true)
				if (err != nil) != tt.wantErr {
					t.Errorf("VariationGlobalCodeDirectoryJS() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("VariationGlobalCodeDirectoryJS() = %v, want %v", got, tt.want)
				}
			})

		}
	}
}

func TestVariationGlobalCodeDirectoryCSS(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []TestStruct{
		{
			name:        "ExistingDirectory",
			workingDir:  currentDir,
			code:        ".id{ \"color\" : black}",
			accountID:   mockAccountID,
			campaignID:  mockCampaignID,
			variationID: mockVariationID,
			want:        currentDir + "/.abtasty/" + mockAccountID + "/" + mockCampaignID + "/" + mockVariationID + "/variationGlobalCode.css",
			wantErr:     false,
		},
		{
			name:        "NonExistingDirectory",
			workingDir:  "/path/to/nonexistent/directory",
			code:        ".id{ \"color\" : black}",
			accountID:   mockAccountID,
			campaignID:  mockCampaignID,
			variationID: mockVariationID,
			want:        "",
			wantErr:     true,
		},
	}

	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				got, err := VariationGlobalCodeDirectoryCSS(tt.workingDir, tt.accountID, tt.campaignID, tt.variationID, tt.code, true)
				if (err != nil) != tt.wantErr {
					t.Errorf("VariationGlobalCodeDirectoryCSS() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("VariationGlobalCodeDirectoryCSS() = %v, want %v", got, tt.want)
				}
			})

		}
	}
}

func TestModificationCodeDirectory(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []TestStruct{
		{
			name:           "ExistingDirectory",
			workingDir:     currentDir,
			code:           "console.log('Hello, World!')",
			accountID:      mockAccountID,
			campaignID:     mockCampaignID,
			variationID:    mockVariationID,
			modificationID: mockModificationID,
			selector:       mockSelector,
			want:           currentDir + "/.abtasty/" + mockAccountID + "/" + mockCampaignID + "/" + mockVariationID + "/" + mockModificationID + "/element.js",
			wantErr:        false,
		},
		{
			name:           "NonExistingDirectory",
			workingDir:     "/path/to/nonexistent/directory",
			code:           "console.log('Hello, World!')",
			accountID:      mockAccountID,
			campaignID:     mockCampaignID,
			variationID:    mockVariationID,
			modificationID: mockModificationID,
			selector:       mockSelector,
			want:           "",
			wantErr:        true,
		},
	}

	for i, tt := range tests {
		if i == 0 {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ModificationCodeDirectory(tt.workingDir, tt.accountID, tt.campaignID, tt.variationID, tt.modificationID, tt.selector, []byte(tt.code), true)
				if (err != nil) != tt.wantErr {
					t.Errorf("ModificationCodeDirectory() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("ModificationCodeDirectory() = %v, want %v", got, tt.want)
				}
			})

		}
	}
}

func TestAddHeaderSelectorComment(t *testing.T) {
	re := regexp.MustCompile(`/\*\s*Selector: (.+)*\s*\*/`)
	fileCode := AddHeaderSelectorComment("example selector", []byte("console.log('Hello World !')"), re)
	fileContent := []byte("/* Selector: example selector */\nconsole.log('Hello World !')")
	assert.Equal(t, fileContent, fileCode)
}
