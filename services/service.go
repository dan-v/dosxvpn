package services

import (
	"bytes"
	"text/template"
)

type Service interface {
	UserData() string
}

func GenerateCloudConfig(authorizedKey []byte, services []Service) (string, error) {
	var userData string
	for _, service := range services {
		userData += service.UserData()
	}
	t, err := template.New("userdata").Parse(userData)
	if err != nil {
		return "", err
	}

	type userDataParams struct {
		SSHAuthorizedKey string
	}
	params := userDataParams{
		SSHAuthorizedKey: string(authorizedKey),
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
