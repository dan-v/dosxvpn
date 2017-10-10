package genconfig

const androidConfigTemplate = `
{
    "uuid": "{{.UUID}}",
    "name": "{{.Name}}",
    "type": "ikev2-cert",
    "remote": {
		"addr": "{{.IP}}",
		"cert": "{{.CACert}}"
	},
	"split-tunneling": {
		"block-ipv4": true,
		"block-ipv6": true
    },
    "local": {
        "p12": "{{.PrivateKey}}"
    },
    "mtu": 1280
}`
