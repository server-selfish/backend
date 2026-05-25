package pkg

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"html/template"

	docker_infra "github.com/server-selfish/backend/internal/infra/docker"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

func ParseTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParseTemplateFromEmbed(name string, data interface{}) (string, error) {
	content, err := docker_infra.GetTemplate(name)
	if err != nil {
		return "", err
	}

	return ParseTemplate(string(content), data)
}

func parseRSAPrivateKeyFromPEM(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, defined_error.ErrInvalidPrivateKey
	}

	if pk, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return pk, nil
	}

	privAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if pk, ok := privAny.(*rsa.PrivateKey); ok {
			return pk, nil
		}
	}

	return nil, defined_error.ErrFailedParseRSAPrivateKey
}
