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
coreos:
  update:
    reboot-strategy: reboot
  locksmith:
    window-start: 10:00
    window-length: 1h
  units:
    - name: "etcd2.service"
      command: "start"
    - name: dummy-interface.service
      command: start
      content: |
        [Unit]
        Description=Creates a dummy local interface

        [Service]
        User=root
        Type=oneshot
        ExecStart=/bin/sh -c "modprobe dummy; ip link set dummy0 up; ifconfig dummy0 1.1.1.1/32; echo 1.1.1.1         pi.hole >> /etc/hosts"
`
}
