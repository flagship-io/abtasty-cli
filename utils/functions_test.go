package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecuteCommand(t *testing.T) {
	mockCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "Hello, World!")
		},
	}

	output, err := ExecuteCommand(mockCmd)
	if err != nil {
		t.Errorf("ExecuteCommand() error = %v", err)
	}

	expectedOutput := "Hello, World!\n"
	if output != expectedOutput {
		t.Errorf("ExecuteCommand() returned unexpected output. Got: %s, Want: %s", output, expectedOutput)
	}
}

func TestGetFeatureExperimentationHost(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     "https://staging-api.flagship.io",
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     "https://api.flagship.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetFeatureExperimentationHost()

			if got != tt.want {
				t.Errorf("GetFeatureExperimentationHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebExperimentationHost(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     "https://staging-api.abtasty.com/api",
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     "https://api.abtasty.com/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetWebExperimentationHost()

			if got != tt.want {
				t.Errorf("GetWebExperimentationHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebExperimentationBackEndHost(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     "https://staging-api.abtasty.com/backend",
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     "https://api.abtasty.com/backend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetWebExperimentationBackEndHost()

			if got != tt.want {
				t.Errorf("GetWebExperimentationBackEndHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostFeatureExperimentationAuth(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     "https://staging-auth.flagship.io",
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     "https://auth.flagship.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetHostFeatureExperimentationAuth()

			if got != tt.want {
				t.Errorf("GetHostFeatureExperimentationAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*

func GetHostWebExperimentationAuth() string {
	if os.Getenv("ABT_STAGING") == "true" {
		return "https://staging-api-auth.abtasty.com"
	}

	return "https://api-auth.abtasty.com"
}

*/

func TestGetHostWebExperimentationAuth(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     "https://staging-api-auth.abtasty.com",
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     "https://api-auth.abtasty.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetHostWebExperimentationAuth()

			if got != tt.want {
				t.Errorf("GetHostWebExperimentationAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func GetWebExperimentationBrowserAuth(clientId, clientSecret string) string {
	if os.Getenv("ABT_STAGING") == "true" {
		return fmt.Sprintf(`https://staging-auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientId, clientSecret)
	}

	return fmt.Sprintf(`https://auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientId, clientSecret)
}

*/

func TestGetWebExperimentationBrowserAuth(t *testing.T) {
	clientID := "ClientID"
	clientSecret := "ClientSecret"

	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     fmt.Sprintf(`https://staging-auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientID, clientSecret),
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     fmt.Sprintf(`https://auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientID, clientSecret),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetWebExperimentationBrowserAuth(clientID, clientSecret)

			if got != tt.want {
				t.Errorf("GetWebExperimentationBrowserAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWebExperimentationBrowserAuthSuccess(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "Staging environment",
			envValue: "STAGING",
			want:     "https://staging-auth.abtasty.com/authorization-granted",
		},
		{
			name:     "Production environment",
			envValue: "",
			want:     "https://auth.abtasty.com/authorization-granted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ABT_ENV", tt.envValue)
			defer os.Unsetenv("ABT_ENV")

			got := GetWebExperimentationBrowserAuthSuccess()

			if got != tt.want {
				t.Errorf("GetWebExperimentationBrowserAuthSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
