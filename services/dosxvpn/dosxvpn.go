package dosxvpn

type Service struct{}

func (s Service) UserData() string {
	return `
    - name: dosxvpn-sysctl.service
      enable: true
      command: start
      content: |
        [Unit]
        Description=Handles settings for sysctl

        [Service]
        Type=oneshot
        User=root
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.ip_forward=1
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.forwarding=1
        ExecStartPre=/usr/sbin/sysctl -w net.ipv6.conf.all.forwarding=1
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.accept_source_route=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.default.accept_source_route=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.accept_redirects=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.default.accept_redirects=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.secure_redirects=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.default.secure_redirects=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.icmp_ignore_bogus_error_responses=1
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.rp_filter=1
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.default.rp_filter=1
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.send_redirects=0
        ExecStartPre=/usr/sbin/sysctl -w net.ipv4.conf.all.send_redirects=0
        ExecStartPre=/usr/bin/echo 1 > /proc/sys/net/ipv4/route/flush
        ExecStartPre=/usr/bin/echo 1 > /proc/sys/net/ipv6/route/flush
        ExecStart=/usr/bin/echo
    - name: dosxvpn-update.service
      content: |
        [Unit]
        Description=Handles updates for dosxvpn

        [Service]
        Type=oneshot
        ExecStartPre=/usr/bin/docker pull dosxvpn/updater:latest
        ExecStart=/usr/bin/docker run --rm --privileged -v /var/run/docker.sock:/var/run/docker.sock dosxvpn/updater:latest
    - name: dosxvpn-update.timer
      enable: true
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
        ExecStartPre=/usr/bin/docker pull dosxvpn/strongswan:latest
        ExecStart=/usr/bin/docker run --name dosxvpn --privileged --net=host -v ipsec.d:/etc/ipsec.d -v strongswan.d:/etc/strongswan.d -v /lib/modules:/lib/modules -v /etc/localtime:/etc/localtime -e VPN_DOMAIN=$public_ipv4 dosxvpn/strongswan:latest
        ExecStop=/usr/bin/docker stop dosxvpn
`
}
