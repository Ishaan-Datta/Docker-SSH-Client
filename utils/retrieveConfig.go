package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type KerberosAuth struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Realm     string `json:"realm"`
	KdcServer string `json:"kdc"`
}

type SAMLAuth struct {
	IdpURL     string `json:"idp_url"`
	SpEntityID string `json:"sp_entity_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type OAuth2Auth struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	TokenURL     string   `json:"token_url"`
	Scopes       []string `json:"scopes"`
	RefreshToken string   `json:"refresh_token"`
}

type AuthConfig struct {
	Kerberos   *KerberosAuth   `json:"kerberos,omitempty"`
	SAML       *SAMLAuth       `json:"saml,omitempty"`
	OAuth2     *OAuth2Auth     `json:"oauth2,omitempty"`
}

type Config struct {
	Type      string     `json:"type"`
	Host      string     `json:"host"`
	CertPath  string     `json:"cert_path"`
	TLSVerify string     `json:"tls_verify"`
	Auth      AuthConfig `json:"auth_methods"`
}

func parseAuthMethod(configMap map[string]interface{}) (AuthConfig, error) {
	authMethod, ok := configMap["auth_method"].(string)
	if !ok {
		return AuthConfig{}, fmt.Errorf("auth_method not provided or invalid")
	}

	switch authMethod {
	case "Kerberos":
		return AuthConfig{
			Kerberos: &KerberosAuth{
				Username:  configMap["username"].(string),
				Password:  configMap["password"].(string),
				Realm:     configMap["realm"].(string),
				KdcServer: configMap["kdc"].(string),
			},
		}, nil
	case "SAML":
		return AuthConfig{
			SAML: &SAMLAuth{
				IdpURL:     configMap["idp_url"].(string),
				SpEntityID: configMap["sp_entity_id"].(string),
				Username:   configMap["username"].(string),
				Password:   configMap["password"].(string),
			},
		}, nil
	case "OAuth2":
		scopes := strings.Split(configMap["scopes"].(string), ",")
		return AuthConfig{
			OAuth2: &OAuth2Auth{
				ClientID:     configMap["client_id"].(string),
				ClientSecret: configMap["client_secret"].(string),
				TokenURL:     configMap["token_url"].(string),
				Scopes:       scopes,
				RefreshToken: configMap["refresh_token"].(string),
			},
		}, nil
	default:
		return AuthConfig{}, fmt.Errorf("unsupported auth_method: %s", authMethod)
	}
}

func parseConfigFile(configFilePath string) ([]Config, []Config, error) {
	// Read the JSON configuration file
	json_data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse JSON into a slice of generic maps
	var configMaps []map[string]interface{}
	err = json.Unmarshal([]byte(json_data), &configMaps)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing config file: %w", err)
	}

	var localConfig []Config
	var remoteConfigs []Config

	for _, configMap := range configMaps {
		var config Config
		config.Host = configMap["host"].(string)
		config.CertPath = configMap["cert_path"].(string)
		config.TLSVerify = configMap["tls_verify"].(string)

		if configMap["type"].(string) == "local" {
			config.Auth = AuthConfig{}
			localConfig = append(localConfig, config)
		} else if configMap["type"].(string) == "remote" {
			auth, err := parseAuthMethod(configMap)

			if err != nil {
				return nil, nil, err
			}

			config.Auth = auth
			remoteConfigs = append(remoteConfigs, config)
		}

	}

	return localConfig, remoteConfigs, nil
}

// validate results
func RetrieveLocalConfiguration(configFilePath string) (configurationList []Config, err error) {
	configurations, _, err := parseConfigFile(configFilePath)

	if err != nil {
		return []Config{}, err
	}

	if configurations == nil{
		return []Config{}, fmt.Errorf("no valid local configurations found")
	}

	if len(configurations) > 1 {
		return []Config{}, fmt.Errorf("more than one local configuration found")
	}

	return configurations, nil
}

// validate results (no empty fields)
func RetrieveRemoteConfiguration(configFilePath string) (configurationList []Config, err error) {
	_, configurations, err := parseConfigFile(configFilePath)

	if err != nil {
		return []Config{}, nil
	}

	return configurations, nil
}

func RetrieveRemoteConfigurationNames(configFilePath string) (configurationNames []string, err error) {
	var hostNames []string

	configs, err := RetrieveRemoteConfiguration(configFilePath)

	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		hostNames = append(hostNames, config.Host)
	}

	return hostNames, nil
}