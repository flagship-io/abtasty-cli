package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type HTTPListResponseFE[T any] struct {
	Items             []T `json:"items"`
	CurrentItemsCount int `json:"current_items_count"`
	CurrentPage       int `json:"current_page"`
	TotalCount        int `json:"total_count"`
	ItemsPerPage      int `json:"items_per_page"`
	LastPage          int `json:"last_page"`
}

type HTTPListResponseWE[T any] struct {
	Data       []T        `json:"_data"`
	Pagination Pagination `json:"_pagination"`
}

type Pagination struct {
	Total      int `json:"_total"`
	Pages      int `json:"_pages"`
	Page       int `json:"_page"`
	MaxPerPage int `json:"_max_per_page"`
}

func ExecuteCommand(cmd *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occurred: %s", err)

	}

	return buf.String(), err
}

func CheckSingleFlag(bool1, bool2 bool) bool {
	count := 0
	if bool1 {
		count++
	}
	if bool2 {
		count++
	}

	return count == 1
}

func GetFeatureExperimentationHost() string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return "https://staging-api.flagship.io"
	}

	return "https://api.flagship.io"
}

func GetWebExperimentationHost() string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return "https://staging-api.abtasty.com/api"
	}

	return "https://api.abtasty.com/api"
}

func GetWebExperimentationBackEndHost() string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return "https://staging-api.abtasty.com/backend"
	}

	return "https://api.abtasty.com/backend"
}

func GetHostFeatureExperimentationAuth() string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return "https://staging-auth.flagship.io"
	}

	return "https://auth.flagship.io"
}

func GetHostWebExperimentationAuth() string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return "https://staging-api-auth.abtasty.com"
	}

	return "https://api-auth.abtasty.com"
}

func GetWebExperimentationBrowserAuth(clientId, clientSecret string) string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return fmt.Sprintf(`https://staging-auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientId, clientSecret)
	}

	return fmt.Sprintf(`https://auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientId, clientSecret)
}

func GetWebExperimentationBrowserAuthSuccess() string {
	if os.Getenv("ABT_ENV") == "STAGING" {
		return "https://staging-auth.abtasty.com/authorization-granted"
	}

	return "https://auth.abtasty.com/authorization-granted"
}

func DefaultGlobalCodeWorkingDir() string {
	wdDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error occurred: %s", err)
	}

	return wdDir
}
