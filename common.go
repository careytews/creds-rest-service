package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
)

// Fetches a user's account name.  This is the first email address listed
// in the EmailAddresses list in the profile.
func getEmail(tokenInfo *AccessTokenInfo) string {
	return tokenInfo.Email
}

// Fetches the credential index for a user, returned as a list of credentials.
func getIndex(user string) (credIndex, error) {

	// Get a Google Storage client.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// Construct pathname to index.
	bucket := getBucket()
	object := user + "/INDEX"

	// Fetch the object.
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer rc.Close()

	// Scan the payload, parsing credential objects.
	creds := make(credIndex, 0)
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {

		// First pass of credential index objects, just to get the type
		// field.
		var ent map[string]string
		err := json.Unmarshal([]byte(scanner.Text()), &ent)
		if err != nil {
			return nil, err
		}

		if ent["type"] == "vpn" {

			// Parse VPN credential index entry.
			cred := &VpnCredential{}
			err := json.Unmarshal([]byte(scanner.Text()), cred)
			if err != nil {
				return nil, err
			}

			// Construct a unique ID from the device name.
			cred.ID = "vpn:" + cred.Device
			cred.User = user
			creds[cred.ID] = cred

		} else if ent["type"] == "web" {

			// Parse web credential index entry.
			cred := &WebCredential{}
			err := json.Unmarshal([]byte(scanner.Text()), cred)
			if err != nil {
				return nil, err
			}

			// Construct a unique ID from the bundle name.
			cred.ID = "web:" + cred.Bundle
			cred.User = user
			creds[cred.ID] = cred

		} else {

			// Ignore if not VPN or web type.
			msg := "Can't understand type " + ent["type"] + " continuing."
			log.Println(msg)
		}
	}

	// Error if scanner failed with errors.
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return creds list.
	return creds, nil

}

func getBucket() string {
	bucket := os.Getenv("CREDENTIALS_BUCKET")
	if bucket == "" {
		bucket = "trust-networks-credentials"
	}
	return bucket
}

func getProject() string {
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if project == "" {
		project = "trust-networks"
	}
	return project
}

func logAndReturnError(w http.ResponseWriter, caller string, call string, msg string, err error) {
	log.Printf("ERROR: %s: %s - %s\n", caller, call, err.Error())
	response := ErrorMessageResponse{Message: msg}
	json.NewEncoder(w).Encode(response)
}

// This can be rewritten when we start to use the token from a header
func getToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}
