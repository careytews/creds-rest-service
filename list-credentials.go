package main

import (
	"encoding/json"
	"net/http"
)

func listCredentials(w http.ResponseWriter, r *http.Request) {

	const functionName = "listCredentials"
	const errorMessage = "Failed to retrieve index"

	// Extract the token from the request
	token := getToken(r)

	// Get the token info for the token
	tokenInfo, err := getTokenInfo(token)
	if err != nil {
		logAndReturnError(w, functionName, "getTokenInfo", errorMessage, err)
		return
	}

	// Get the email
	email := getEmail(tokenInfo)

	// Get the creds index
	creds, err := getIndex(email)
	if err != nil {
		logAndReturnError(w, functionName, "getIndex", errorMessage, err)
		return
	}

	// Success! We have an index, and we can return it
	credentials := make([]*CredentialsInfo, 0)
	for _, cred := range creds {
		credentials = append(credentials, cred.(Credential).AsCredentialsInfo())
	}

	response := ListCredentialsResponse{Credentials: credentials}

	json.NewEncoder(w).Encode(response)
}
