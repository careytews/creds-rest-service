package main

// CredentialsInfo describes the details of a set of credentials, Web or VPN
type CredentialsInfo struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Start       string `json:"start,omitempty"`
	End         string `json:"end,omitempty"`
}

// ListCredentialsResponse describes the response to a listcreds request
type ListCredentialsResponse struct {
	Credentials []*CredentialsInfo `json:"credentials,omitempty"`
}

// CredentialsDownloadResponse describes the response to a cred download request
type CredentialsDownloadResponse struct {
	Name     string `json:"name,omitempty"`
	Filename string `json:"filename,omitempty"`
	Data     string `json:"data,omitempty"`
	Password string `json:"password,omitempty"`
}

// ErrorMessageResponse is sent if there's an error message to deliver,
// so we avoid using http.Error and the subsequent error in the console.
type ErrorMessageResponse struct {
	Message string `json:message,omitempty"`
}
