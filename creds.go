/****************************************************************************

  Credential generic code.

****************************************************************************/

package main

type credIndex map[string]Credential

// Credential interface
type Credential interface {

	// Describes credential to stdout, human-readable
	Describe()

	// Returns unique cred ID
	GetID() string

	// GetKey returns the key of the cred
	GetKey() string

	// Gets the ID string of the cred
	GetName() string

	// Gets the file (or bundle name) without the extension
	GetFilename() string

	// Output the cred as a CredentialInfo struct
	AsCredentialsInfo() *CredentialsInfo
}
