package pihole

type Service struct{}

func (s Service) UserData() string {
	return `
    - name: pihole-etc-host.service
      command: start
      content: |
        [Unit]
        Description=pihole /etc/hosts entry
        ConditionFirstBoot=true

        [Service]
        User=root
        Type=oneshot
        ExecStart=/bin/sh -c "echo 1.1.1.2         pi.hole >> /etc/hosts"
    - name: pihole.service
      command: start
      content: |
        [Unit]
        Description=pihole
        After=docker.service,dummy-interface.service

        [Service]
        User=core
        Restart=always
        TimeoutStartSec=0
        KillMode=none
        EnvironmentFile=/etc/environment
        ExecStartPre=-/usr/bin/docker kill pihole
        ExecStartPre=-/usr/bin/docker rm pihole
        ExecStartPre=/usr/bin/docker pull diginc/pi-hole:latest
        ExecStart=/usr/bin/docker run --name pihole --net=host -e DNS1=1.1.1.1 -e ServerIP=1.1.1.2 -e ServerIPv6=fd9d:bc11:4020:: -e WEBPASSWORD=dosxvpn -v pihole-etc:/etc/pihole -v pihole-dnsmasq.d:/etc/dnsmasq.d diginc/pi-hole:latest
        ExecStop=/usr/bin/docker stop pihole`
}
