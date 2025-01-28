package utils

// parse json config based on assigned type, pull credentials and fire gRPC request, encrypt?

// "golang.org/x/crypto/ssh"
// "github.com/coreos/go-oidc"
// "golang.org/x/oauth2"

// // Kerberos authentication method
// func kerberosAuth() (ssh.AuthMethod, error) {
// 	// Example: using kinit to obtain a Kerberos ticket
// 	cmd := exec.Command("kinit", "user@REALM")
// 	if err := cmd.Run(); err != nil {
// 		return nil, err
// 	}

// 	// Use the Kerberos ticket for authentication
// 	return ssh.PublicKeysCallback(sshAgent()), nil
// }

// // OAuth2 authentication method
// func oauth2Auth() (ssh.AuthMethod, error) {
// 	// Configure OAuth2 client
// 	ctx := context.Background()
// 	provider, err := oidc.NewProvider(ctx, "https://your-oauth2-provider.com")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get provider: %v", err)
// 	}

// 	config := oauth2.Config{
// 		ClientID:     "your-client-id",
// 		ClientSecret: "your-client-secret",
// 		Endpoint:     provider.Endpoint(),
// 		RedirectURL:  "http://localhost:8080/callback",
// 		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
// 	}

// 	// Exchange code for token
// 	authCodeURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
// 	fmt.Printf("Visit the following URL to authenticate: %v\n", authCodeURL)

// 	// Here you would implement a way to capture the authorization code from the redirect URL,
// 	// for simplicity, we assume it's already captured and stored in authCode.
// 	var authCode string
// 	fmt.Println("Enter the authorization code:")
// 	fmt.Scan(&authCode)

// 	token, err := config.Exchange(ctx, authCode)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to exchange token: %v", err)
// 	}

// 	// Example token-based SSH authentication
// 	return ssh.Password(token.AccessToken), nil
// }

// // SAML authentication method (Mock)
// func samlAuth() (ssh.AuthMethod, error) {
// 	// Placeholder for SAML authentication

// 	// Step 1: Initiate SAML authentication request to the Identity Provider (IdP)
// 	// This would typically involve redirecting the user to the IdP for authentication

// 	// Step 2: Handle the SAML response and extract the assertion
// 	// After authentication, the IdP will send back a SAML assertion (usually via POST)

// 	// Step 3: Validate the SAML assertion and extract the necessary information (e.g., user identity)

// 	// Step 4: Use the extracted information for SSH authentication
// 	// For simplicity, we mock this by returning a dummy SSH password method

// 	fmt.Println("SAML authentication is not fully implemented yet. Using mock credentials.")
// 	return ssh.Password("mock-password"), nil
// }

// // Helper function for ssh-agent
// func sshAgent() ssh.Signer {
// 	// Load the SSH agent
// 	// This is a mock; in a real implementation, you would load the private key from the SSH agent
// 	return nil
// }

// import (
// 	"context"
// 	"fmt"
// 	"os/exec"

// 	"golang.org/x/crypto/ssh"

// 	"github.com/coreos/go-oidc"
// 	"golang.org/x/oauth2"
// )

// // Kerberos authentication method
// func kerberosAuth() (ssh.AuthMethod, error) {
// 	// Example: using kinit to obtain a Kerberos ticket
// 	cmd := exec.Command("kinit", "user@REALM")
// 	if err := cmd.Run(); err != nil {
// 		return nil, err
// 	}

// 	// Use the Kerberos ticket for authentication
// 	return ssh.PublicKeysCallback(sshAgent()), nil
// }

// // OAuth2 authentication method
// func oauth2Auth() (ssh.AuthMethod, error) {
// 	// Configure OAuth2 client
// 	ctx := context.Background()
// 	provider, err := oidc.NewProvider(ctx, "https://your-oauth2-provider.com")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get provider: %v", err)
// 	}

// 	config := oauth2.Config{
// 		ClientID:     "your-client-id",
// 		ClientSecret: "your-client-secret",
// 		Endpoint:     provider.Endpoint(),
// 		RedirectURL:  "http://localhost:8080/callback",
// 		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
// 	}

// 	// Exchange code for token
// 	authCodeURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
// 	fmt.Printf("Visit the following URL to authenticate: %v\n", authCodeURL)

// 	// Here you would implement a way to capture the authorization code from the redirect URL,
// 	// for simplicity, we assume it's already captured and stored in authCode.
// 	var authCode string
// 	fmt.Println("Enter the authorization code:")
// 	fmt.Scan(&authCode)

// 	token, err := config.Exchange(ctx, authCode)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to exchange token: %v", err)
// 	}

// 	// Example token-based SSH authentication
// 	return ssh.Password(token.AccessToken), nil
// }

// // SAML authentication method (Mock)
// func samlAuth() (ssh.AuthMethod, error) {
// 	// Placeholder for SAML authentication

// 	// Step 1: Initiate SAML authentication request to the Identity Provider (IdP)
// 	// This would typically involve redirecting the user to the IdP for authentication

// 	// Step 2: Handle the SAML response and extract the assertion
// 	// After authentication, the IdP will send back a SAML assertion (usually via POST)

// 	// Step 3: Validate the SAML assertion and extract the necessary information (e.g., user identity)

// 	// Step 4: Use the extracted information for SSH authentication
// 	// For simplicity, we mock this by returning a dummy SSH password method

// 	fmt.Println("SAML authentication is not fully implemented yet. Using mock credentials.")
// 	return ssh.Password("mock-password"), nil
// }

// // Helper function for ssh-agent
// func sshAgent() ssh.Signer {
// 	// Load the SSH agent
// 	// This is a mock; in a real implementation, you would load the private key from the SSH agent
// 	return nil
// }