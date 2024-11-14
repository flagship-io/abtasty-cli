package utils

import (
	"fmt"
	"log"
	"os"
)

func GetFeatureExperimentationHost() string {
	if os.Getenv("ABT_STAGING") == "true" {
		return "https://staging-api.flagship.io"
	}

	return "https://api.flagship.io"
}

func GetWebExperimentationHost() string {
	if os.Getenv("ABT_STAGING") == "true" {
		return "https://staging-api.abtasty.com/api"
	}

	return "https://api.abtasty.com/api"
}

func GetWebExperimentationBackEndHost() string {
	if os.Getenv("ABT_STAGING") == "true" {
		return "https://staging-api.abtasty.com/backend"
	}

	return "https://api.abtasty.com/backend"
}

func GetHostFeatureExperimentationAuth() string {
	if os.Getenv("ABT_STAGING") == "true" {
		return "https://staging-auth.flagship.io"
	}

	return "https://auth.flagship.io"
}

func GetHostWebExperimentationAuth() string {
	if os.Getenv("ABT_STAGING") == "true" {
		return "https://staging-api-auth.abtasty.com"
	}

	return "https://api-auth.abtasty.com"
}

func GetWebExperimentationBrowserAuth(clientId, clientSecret string) string {
	if os.Getenv("ABT_STAGING") == "true" {
		return fmt.Sprintf(`https://staging-auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientId, clientSecret)
	}

	return fmt.Sprintf(`https://auth.abtasty.com/authorize?client_id=%s&client_secret=%s&redirect_uri=http://localhost:8010/auth/callback`, clientId, clientSecret)
}

func GetWebExperimentationBrowserAuthSuccess() string {
	if os.Getenv("ABT_STAGING") == "true" {
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

const FEATURE_EXPERIMENTATION = "fe"
const WEB_EXPERIMENTATION = "we"
const HOME_CLI = ".cli"
const HIT_ANALYTICS_URL = "https://events.flagship.io/analytics"
