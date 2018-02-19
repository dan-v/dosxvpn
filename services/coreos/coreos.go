package coreos

type Service struct{}

func (s Service) UserData() string {
	return `#cloud-config
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
  - path: /var/lib/iptables/rules-save
    permissions: 0644
    owner: root:root
    content: |
      *nat
      :PREROUTING ACCEPT [0:0]
      :POSTROUTING ACCEPT [0:0]
      -A POSTROUTING -s 192.168.99.0/24 -m policy --pol none --dir out -j MASQUERADE
      COMMIT
      *filter
      :INPUT DROP [0:0]
      :FORWARD DROP [0:0]
      :OUTPUT ACCEPT [0:0]
      -A INPUT -i lo -j ACCEPT
      -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
      -A INPUT -p esp -j ACCEPT
      -A INPUT -p ah -j ACCEPT
      -A INPUT -p ipencap -m policy --dir in --pol ipsec --proto esp -j ACCEPT
      -A INPUT -p icmp --icmp-type echo-request -m hashlimit --hashlimit-upto 5/s --hashlimit-mode srcip --hashlimit-srcmask 32 --hashlimit-name icmp-echo-drop -j ACCEPT
      -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --set --name SSH
      -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --update --seconds 60 --hitcount 10 --rttl --name SSH -j DROP
      -A INPUT -p tcp --dport 22 -m state --state NEW -j ACCEPT
      -A INPUT -p udp -m multiport --dports 500,4500 -j ACCEPT
      -A INPUT -p tcp --destination-port 443 -j REJECT --reject-with tcp-reset
      -A INPUT -p udp --destination-port 80 -j REJECT --reject-with icmp-port-unreachable
      -A INPUT -p udp --destination-port 443 -j REJECT --reject-with icmp-port-unreachable
      -A INPUT -d 1.1.1.2 -p udp -j ACCEPT
      -A INPUT -d 1.1.1.2 -p tcp -j ACCEPT
      -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
      -A FORWARD -m conntrack --ctstate NEW -s 192.168.99.0/24 -m policy --pol ipsec --dir in -j ACCEPT
      COMMIT
  - path: /var/lib/ip6tables/rules-save
    permissions: 0644
    owner: root:root
    content: |
      *nat
      :PREROUTING ACCEPT [0:0]
      :POSTROUTING ACCEPT [0:0]
      -A POSTROUTING -s fd9d:bc11:4020::/48 -m policy --pol none --dir out -j MASQUERADE
      COMMIT
      *filter
      :INPUT DROP [0:0]
      :FORWARD DROP [0:0]
      :OUTPUT ACCEPT [0:0]
      :ICMPV6-CHECK - [0:0]
      :ICMPV6-CHECK-LOG - [0:0]
      -A INPUT -i lo -j ACCEPT
      -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
      -A INPUT -p esp -j ACCEPT
      -A INPUT -m ah -j ACCEPT
      -A INPUT -p icmpv6 --icmpv6-type echo-request -m hashlimit --hashlimit-upto 5/s --hashlimit-mode srcip --hashlimit-srcmask 32 --hashlimit-name icmp-echo-drop -j ACCEPT
      -A INPUT -p icmpv6 --icmpv6-type router-advertisement -m hl --hl-eq 255 -j ACCEPT
      -A INPUT -p icmpv6 --icmpv6-type neighbor-solicitation -m hl --hl-eq 255 -j ACCEPT
      -A INPUT -p icmpv6 --icmpv6-type neighbor-advertisement -m hl --hl-eq 255 -j ACCEPT
      -A INPUT -p icmpv6 --icmpv6-type redirect -m hl --hl-eq 255 -j ACCEPT
      -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --set --name SSH
      -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --update --seconds 60 --hitcount 10 --rttl --name SSH -j DROP
      -A INPUT -p tcp --dport 22 -m state --state NEW -j ACCEPT
      -A INPUT -p udp -m multiport --dports 500,4500 -j ACCEPT
      -A INPUT -p tcp --destination-port 443 -j REJECT --reject-with tcp-reset
      -A INPUT -p udp --destination-port 80 -j REJECT --reject-with icmp6-port-unreachable
      -A INPUT -p udp --destination-port 443 -j REJECT --reject-with icmp6-port-unreachable
      -A INPUT -d fd9d:bc11:4020::/48 -p udp -j ACCEPT
      -A INPUT -d fd9d:bc11:4020::/48 -p tcp -j ACCEPT
      -A FORWARD -j ICMPV6-CHECK
      -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
      -A FORWARD -m conntrack --ctstate NEW -s fd9d:bc11:4020::/48 -m policy --pol ipsec --dir in -j ACCEPT
      -A ICMPV6-CHECK -p icmpv6 -m hl ! --hl-eq 255 --icmpv6-type router-solicitation -j ICMPV6-CHECK-LOG
      -A ICMPV6-CHECK -p icmpv6 -m hl ! --hl-eq 255 --icmpv6-type router-advertisement -j ICMPV6-CHECK-LOG
      -A ICMPV6-CHECK -p icmpv6 -m hl ! --hl-eq 255 --icmpv6-type neighbor-solicitation -j ICMPV6-CHECK-LOG
      -A ICMPV6-CHECK -p icmpv6 -m hl ! --hl-eq 255 --icmpv6-type neighbor-advertisement -j ICMPV6-CHECK-LOG
      -A ICMPV6-CHECK-LOG -j LOG --log-prefix "ICMPV6-CHECK-LOG DROP "
      -A ICMPV6-CHECK-LOG -j DROP
      COMMIT

coreos:
  update:
    reboot-strategy: reboot
  locksmith:
    window-start: 10:00
    window-length: 1h
  units:
    - name: iptables-restore.service
      enable: true
      command: start
    - name: ip6tables-restore.service
      enable: true
      command: start
    - name: dummy-interface.service
      command: start
      enable: true
      content: |
        [Unit]
        Description=Creates a dummy local interface

        [Service]
        User=root
        Type=oneshot
        ExecStartPre=/bin/sh -c "modprobe dummy"
        ExecStartPre=-/bin/sh -c "ip link add dummy0 type dummy"
        ExecStartPre=/bin/sh -c "ip link set dummy0 up"
        ExecStartPre=-/bin/sh -c "ifconfig dummy0 inet6 add fd9d:bc11:4020::/48"
        ExecStartPre=-/bin/sh -c "ifconfig dummy0 1.1.1.2/32"
        ExecStartPre=-/bin/sh -c "ifconfig dummy0 inet6 add fd9d:bc11:4020::/48"
        ExecStartPre=-/bin/sh -c "ifconfig dummy0 1.1.1.2/32"
        ExecStart=/bin/sh -c "echo"
`
}
