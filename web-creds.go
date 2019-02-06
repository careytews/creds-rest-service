/****************************************************************************

  Web credentials-specific code.

****************************************************************************/

package main

import (
	"encoding/json"
	"log"
	"strings"
)

// WebCredential defintion
type WebCredential struct {
	ID          string `json:"id,omitempty"`
	User        string `json:"user,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Key         string `json:"key,omitempty"`
	Start       string `json:"start,omitempty"`
	End         string `json:"end,omitempty"`
	Bundle      string `json:"bundle,omitempty"`
	Password    string `json:"password,omitempty"`
}

// GetID gets the ID of the credential
func (c WebCredential) GetID() string { return c.ID }

// Describe causes the credential to describe itself
func (c WebCredential) Describe() {
	b, err := json.Marshal(c)
	if err != nil {
		log.Println("describe failed: " + err.Error())
	} else {
		log.Println(string(b))
	}
}

// GetName gets the CN of the cred and the eventual name of the downloaded file
func (c WebCredential) GetName() string {
	return c.Name
}

// GetKey gets the key of the cred
func (c WebCredential) GetKey() string {
	return c.Key
}

// GetFilename of the bundle without the extension
func (c WebCredential) GetFilename() string {
	return strings.Split(c.Bundle, ".")[0]
}

// AsCredentialsInfo returns the web credential in CredentialsInfo format
func (c WebCredential) AsCredentialsInfo() *CredentialsInfo {
	cred := &CredentialsInfo{
		ID:          c.ID,
		Type:        c.Type,
		Name:        c.Name,
		Description: c.Description,
		Start:       c.Start,
		End:         c.End,
	}

	return cred
}
