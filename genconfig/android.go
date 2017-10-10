package genconfig

import (
	"bytes"
	"encoding/base64"
	"text/template"

	uuid "github.com/satori/go.uuid"
)

func GenerateAndroidConfig(ip, name, privateKey, caCert string) (string, error) {
	tmpl := template.Must(template.New("androidconfig").Parse(androidConfigTemplate))
	tmplData := struct {
		IP         string
		Name       string
		PrivateKey string
		CACert     string
		UUID       string
	}{
		IP:         ip,
		Name:       name,
		PrivateKey: base64.StdEncoding.EncodeToString([]byte(privateKey)),
		CACert:     base64.StdEncoding.EncodeToString([]byte(caCert)),
		UUID:       uuid.NewV4().String(),
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, tmplData)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
