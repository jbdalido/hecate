default: linux-amd64

cleanup:
	rm -Rf build/

build: cleanup
	mkdir -p build/

linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/hecate hecate/hecate.go

osx: build
	go build -o build/hecate hecate/hecate.go