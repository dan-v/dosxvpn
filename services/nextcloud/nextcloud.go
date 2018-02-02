package nextcloud

type Service struct{}

func (s Service) UserData() string {
	return `
    - name: nextcloud-db.service
      command: start
      content: |
        [Unit]
        Description=Nextcloud Database
        After=docker.service,dummy-interface.service

        [Service]
        User=core
        Restart=always
        TimeoutStartSec=0
        KillMode=none
        EnvironmentFile=/etc/environment
        ExecStartPre=-/usr/bin/docker kill db_nextcloud
        ExecStartPre=-/usr/bin/docker rm db_nextcloud
        ExecStartPre=/usr/bin/docker pull mariadb:latest
        ExecStart=/usr/bin/docker run --name db_nextcloud -v db-nextcloud-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=nextcloud -e MYSQL_DATABASE=nextcloud -e MYSQL_USER=nextcloud -e MYSQL_PASSWORD=nextcloud mariadb:latest
        ExecStop=/usr/bin/docker stop db_nextcloud
    - name: nextcloud.service
      command: start
      content: |
        [Unit]
        Description=Nextcloud
        After=docker.service,dummy-interface.service,nextcloud-db.service

        [Service]
        User=core
        Restart=always
        TimeoutStartSec=0
        KillMode=none
        EnvironmentFile=/etc/environment
        ExecStartPre=-/usr/bin/docker kill nextcloud
        ExecStartPre=-/usr/bin/docker rm nextcloud
        ExecStartPre=/usr/bin/docker pull nextcloud:latest
        ExecStart=/usr/bin/docker run --name nextcloud -v nextcloud:/var/www/html --link db_nextcloud:db_nextcloud -e MYSQL_DATABASE=nextcloud -e MYSQL_USER=nextcloud -e MYSQL_PASSWORD=nextcloud -e MYSQL_HOST=db_nextcloud -e NEXTCLOUD_ADMIN_USER=admin -e NEXTCLOUD_ADMIN_PASSWORD=dosxvpn -p 8080:80 nextcloud:latest
        ExecStop=/usr/bin/docker stop nextcloud`
}
