/****************************************************************************

  VPN credential-specific code.

****************************************************************************/

package main

import (
	"encoding/json"
	"log"
)

// VpnCredential defintion
type VpnCredential struct {
	ID          string `json:"id,omitempty"`
	User        string `json:"user,omitempty"`
	Type        string `json:"type,omitempty"`
	Device      string `json:"device,omitempty"`
	Description string `json:"description,omitempty"`
	Key         string `json:"key,omitempty"`
	Start       string `json:"start,omitempty"`
	End         string `json:"end,omitempty"`
	DeviceType  string `json:"device_type,omitempty"`
	Uk          string `json:"uk,omitempty"`
	Us          string `json:"us,omitempty"`
}

// GetID gets the ID of the credential
func (c VpnCredential) GetID() string { return c.ID }

// Describe causes the credential to describe itself
func (c VpnCredential) Describe() {
	b, err := json.Marshal(c)
	if err != nil {
		log.Println("describe failed: " + err.Error())
	} else {
		log.Println(string(b))
	}
}

// GetName gets the device name, which is also the name of the file
func (c VpnCredential) GetName() string {
	return c.Device
}

// GetKey gets the key of the cred
func (c VpnCredential) GetKey() string {
	return c.Key
}

// GetFilename gets the device name and leaves off everything after the "-"
func (c VpnCredential) GetFilename() string {
	return c.Device
}

// AsCredentialsInfo returns the VPN credential in CredentialsInfo format
func (c VpnCredential) AsCredentialsInfo() *CredentialsInfo {
	cred := &CredentialsInfo{
		ID:          c.ID,
		Type:        c.Type,
		Name:        c.Device,
		Description: c.Description,
		Start:       c.Start,
		End:         c.End,
	}

	return cred
}
