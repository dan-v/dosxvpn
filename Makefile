all: osx

osx:
	GOOS=darwin GOARCH=amd64 go build -o ./build/osx/x86-64/dosxvpn .
	cd build/osx/x86-64 && zip -r ./dosxvpn-cli.zip ./dosxvpn 
	cd platypus && ./build.sh
	cd build/osx/x86-64 && zip -r ./dosxvpn-app.zip ./dosxvpn.app

clean:
	rm -rf build

.PHONY: osx
