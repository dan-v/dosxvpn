package dosxvpn

type Service struct{}

func (s Service) UserData() string {
	return `
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
        After=docker.service,dummy-interface.service

        [Service]
        User=core
        Restart=always
        TimeoutStartSec=0
        KillMode=none
        EnvironmentFile=/etc/environment
        ExecStartPre=-/usr/bin/docker kill dosxvpn
        ExecStartPre=-/usr/bin/docker rm dosxvpn
        ExecStartPre=/usr/bin/docker pull dosxvpn/strongswan
        ExecStart=/usr/bin/docker run --name dosxvpn --privileged --net=host -v ipsec.d:/etc/ipsec.d -v strongswan.d:/etc/strongswan.d -v /lib/modules:/lib/modules -v /etc/localtime:/etc/localtime -e VPN_DNS=1.1.1.1 -e VPN_DOMAIN=$public_ipv4 dosxvpn/strongswan
        ExecStop=/usr/bin/docker stop dosxvpn
`
}
