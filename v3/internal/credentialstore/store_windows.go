package credentialstore

import (
	"errors"

	"github.com/danieljoos/wincred"
)

func Set(service, user, value string) error {
	credential := wincred.NewGenericCredential(service + ":" + user)
	credential.UserName = user
	credential.CredentialBlob = []byte(value)
	return credential.Write()
}

func Get(service, user string) (string, error) {
	credential, err := wincred.GetGenericCredential(service + ":" + user)
	if err != nil {
		return "", credentialError(err)
	}
	return string(credential.CredentialBlob), nil
}

func Delete(service, user string) error {
	credential := wincred.NewGenericCredential(service + ":" + user)
	return credentialError(credential.Delete())
}

func credentialError(err error) error {
	if errors.Is(err, wincred.ErrElementNotFound) {
		return ErrNotFound
	}
	return err
}
