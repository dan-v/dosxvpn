package dosxvpn

import (
	"bytes"
	"text/template"
)

const baseUserData = `#cloud-config
ssh_authorized_keys:
  - {{.SSHAuthorizedKey}}
write_files:
  - path: /etc/ssh/sshd_config
    permissions: 0600
    owner: root:root
    content: |
      # Use most defaults for sshd configuration.
      UsePrivilegeSeparation sandbox
      Subsystem sftp internal-sftp

      PermitRootLogin no
      AllowUsers core
      PasswordAuthentication no
      ChallengeResponseAuthentication no

coreos:
  update:
    reboot-strategy: reboot
  locksmith:
    window-start: 10:00
    window-length: 1h
  units:
    - name: "etcd2.service"
      command: "start"
    - name: dosxvpn-update.service
      content: |
        [Unit]
        Description=Handles updates for dosxvpn

        [Service]
        Type=oneshot
        ExecStartPre=/usr/bin/docker pull dosxvpn/strongswan-updater
        ExecStart=/usr/bin/docker run --rm --privileged -v /var/run/docker.sock:/var/run/docker.sock dosxvpn/strongswan-updater
    - name: dosxvpn-update.timer
      command: start
      content: |
        [Unit]
        Description=Run dosxvpn-update on schedule

        [Timer]
        OnCalendar=*-*-* 0/12:00:00
    - name: dosxvpn.service
      command: start
      content: |
        [Unit]
        Description=dosxvpn
        After=docker.service

        [Service]
        User=core
        Restart=always
        TimeoutStartSec=0
        KillMode=none
        EnvironmentFile=/etc/environment
        ExecStartPre=-/usr/bin/docker kill dosxvpn
        ExecStartPre=-/usr/bin/docker rm dosxvpn
        ExecStartPre=/usr/bin/docker pull dosxvpn/strongswan
        ExecStart=/usr/bin/docker run --name dosxvpn --privileged -p 500:500/udp -p 4500:4500/udp -v ipsec.d:/etc/ipsec.d -v strongswan.d:/etc/strongswan.d -v /lib/modules:/lib/modules -v /etc/localtime:/etc/localtime -e VPN_DOMAIN=$public_ipv4 dosxvpn/strongswan
        ExecStop=/usr/bin/docker stop dosxvpn
`

type userDataParams struct {
	SSHAuthorizedKey string
}

func buildUserData(opt *options, keypair *sshKeyPair) (string, error) {
	t, err := template.New("userdata").Parse(baseUserData)
	if err != nil {
		return "", err
	}

	params := userDataParams{
		SSHAuthorizedKey: string(keypair.authorizedKey),
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
