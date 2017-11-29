package sshclient

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	KeyPair *SSHKeyPair
}

type SSHKeyPair struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     ssh.PublicKey
	PrivateKeyPEM []byte
	AuthorizedKey []byte
}

func New() (*Client, error) {
	ssh := &Client{}
	err := ssh.generateKeyPair()
	if err != nil {
		return nil, err
	}
	return ssh, nil
}

func (s *Client) openSession(user, host string) (*ssh.Session, error) {
	signer, err := ssh.NewSignerFromKey(s.KeyPair.PrivateKey)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Client) GetFileFromContainer(user, host, containerName, file string) (string, error) {
	command := fmt.Sprintf("docker exec %s cat %s", containerName, file)
	return s.Run(user, host, command)
}

func (s *Client) Run(user, host, command string) (string, error) {
	session, err := s.openSession(user, host)
	if err != nil {
		return "", err
	}

	rOut, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}
	rErr, err := session.StderrPipe()
	if err != nil {
		return "", err
	}
	combined := io.MultiReader(rOut, rErr)

	err = session.Start(command)
	if err != nil {
		return "", err
	}

	var lines []string
	scanner := bufio.NewScanner(combined)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	session.Close()

	output := strings.Join(lines, "\n")
	output = strings.TrimSpace(output)
	return output, nil
}

func (s *Client) generateKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	var pembuf bytes.Buffer
	err = pem.Encode(&pembuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

	s.KeyPair = &SSHKeyPair{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		PrivateKeyPEM: pembuf.Bytes(),
		AuthorizedKey: authorizedKey,
	}
	return nil
}
