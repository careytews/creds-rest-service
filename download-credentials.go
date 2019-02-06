package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
)

func downloadCredentials(w http.ResponseWriter, r *http.Request) {

	const functionName = "downloadCredentials"
	const errorMessage = "Download failed"

	vars := mux.Vars(r)
	token := getToken(r)
	id := vars["id"]
	parts := strings.Split(id, ":")

	err := validateDownloadID(parts)
	if err != nil {
		logAndReturnError(w, functionName, "validateDownloadID", errorMessage, err)
		return
	}

	credType := parts[0]
	credName := parts[1]

	tokenInfo, err := getTokenInfo(token)
	if err != nil {
		logAndReturnError(w, functionName, "getTokenInfo", errorMessage, err)
		return
	}

	email := getEmail(tokenInfo)

	creds, err := getIndex(email)
	if err != nil {
		logAndReturnError(w, functionName, "getIndex", errorMessage, err)
		return
	}

	// TODO: Replace this evil hack with something better
	if credType == "vpn" {
		id = strings.Split(id, ".")[0]
		id = strings.Replace(id, "-uk", "", 1)
		id = strings.Replace(id, "-us", "", 1)
	}

	// getCredentialEntry
	cred := creds[id]
	if cred == nil {
		logAndReturnError(w, functionName, "getCredentialEntry", errorMessage, fmt.Errorf("entry does not exist"))
		return
	}

	data, err := getCredential(email, credName, cred.GetKey())
	if err != nil {
		logAndReturnError(w, functionName, "getCredential", errorMessage, err)
		return
	}

	password := []byte{}

	// getPassword
	if credType == "web" {
		password, err = getCredential(email, cred.GetFilename()+".pass", cred.GetKey())
		if err != nil {
			logAndReturnError(w, functionName, "getPassword", errorMessage, err)
			return
		}
	}

	response := CredentialsDownloadResponse{
		Name:     cred.GetName(),
		Filename: credName,
		Data:     string(data),
		Password: string(password),
	}
	json.NewEncoder(w).Encode(response)
}

// Fetches a decrypted credential payload from the store. The payload is Base-64 encoded.
func getCredential(email string, filename string, key string) ([]byte, error) {

	// Download the object
	object, err := downloadCredential(email + "/" + filename)
	if err != nil {
		return []byte{}, err
	}

	// Get the object as a buffer
	buffer := bytes.NewBuffer(object)

	// De-binhex the file
	rawkey := make([]byte, len(key)/2)
	_, err = hex.Decode(rawkey, []byte(key))
	if err != nil {
		return []byte{}, err
	}

	// Decrypt the file
	key7, err := ckmsDecrypt(email, rawkey)
	if err != nil {
		return []byte{}, err
	}

	// Undo the hexbin encoding on ciphertext.
	ciph := make([]byte, len(buffer.String())/2)
	_, err = hex.Decode(ciph, []byte(buffer.String()))
	if err != nil {
		return []byte{}, err
	}

	// Decryot using AES
	cred, err := decrypt(ciph, key7)
	if err != nil {
		return []byte{}, err
	}

	// This object is so we just decode the content field.
	var obj struct {
		Content string `json:"content"`
	}

	// Decode JSON
	err = json.Unmarshal(cred, &obj)
	if err != nil {
		return []byte{}, err
	}

	return []byte(obj.Content), nil
}

func downloadCredential(object string) ([]byte, error) {

	// Get a Google Storage client.
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return []byte{}, err
	}

	// Construct pathname to cred.
	bucket := getBucket()

	// Fetch the object
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return []byte{}, err
	}

	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func validateDownloadID(id []string) error {
	var err error
	length := len(id)
	if length != 2 {
		err = fmt.Errorf("id has %d part(s), expected 2", length)
	} else if strings.Compare(id[0], "vpn") != 0 && strings.Compare(id[0], "web") != 0 {
		err = fmt.Errorf("invalid ID %s", id[0])
	}
	return err
}
