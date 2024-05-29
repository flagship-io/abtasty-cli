package config

import "github.com/flagship-io/abtasty-cli/cmd/version"

const (
	OutputFormat         = "table"
	GrantType            = "client_credentials"
	Expiration           = 43200
	Scope                = "*"
	ClientID             = "clientID"
	ClientSecret         = "clientSecret"
	Token                = "token"
	AccountID            = "accountID"
	AccountEnvironmentID = "accountEnvironmentID"
)

var DefaultUserAgent = "abtasty-cli/" + version.Version
