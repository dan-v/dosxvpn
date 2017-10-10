package genconfig

import (
	"bytes"
	"encoding/base64"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

func GenerateAppleConfig(ip, name, password, privateKey, caCert, serverCert string) (string, error) {
	tmpl := template.Must(template.New("mobileconfig").Parse(mobileConfigTemplate))
	tmplData := struct {
		IP                 string
		Name               string
		PrivateKey         string
		PrivateKeyPassword string
		CACert             string
		ServerCert         string
		UUID1              string
		UUID2              string
		UUID3              string
		UUID4              string
		UUID5              string
		UUID6              string
	}{
		IP:                 ip,
		Name:               name,
		PrivateKeyPassword: password,
		PrivateKey:         base64.StdEncoding.EncodeToString([]byte(privateKey)),
		CACert:             base64.StdEncoding.EncodeToString([]byte(caCert)),
		ServerCert:         base64.StdEncoding.EncodeToString([]byte(serverCert)),
		UUID1:              uuid.NewV4().String(),
		UUID2:              uuid.NewV4().String(),
		UUID3:              uuid.NewV4().String(),
		UUID4:              uuid.NewV4().String(),
		UUID5:              uuid.NewV4().String(),
		UUID6:              uuid.NewV4().String(),
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, tmplData)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
