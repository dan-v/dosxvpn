all: osx

osx:
	GOOS=darwin GOARCH=amd64 go build -o ./build/osx/x86-64/dosxvpn ./cmd/dosxvpn
	cd platypus && ./build.sh
	cd build/osx/x86-64 && zip -r ./dosxvpn.zip ./dosxvpn.app

clean:
	rm -rf build

.PHONY: osx
