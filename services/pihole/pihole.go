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
        ExecStart=/bin/sh -c "echo 1.1.1.1         pi.hole >> /etc/hosts"
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
        ExecStartPre=/usr/bin/docker pull diginc/pi-hole:alpine
        ExecStart=/usr/bin/docker run --name pihole --net=host -e ServerIP=1.1.1.1 -e WEBPASSWORD=dosxvpn diginc/pi-hole:alpine
        ExecStop=/usr/bin/docker stop pihole`
}
