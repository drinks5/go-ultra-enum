test:
	go build -ldflags "-X main.version=`git describe --tags`" main.go
	go generate
	go test .

release:
	rm -rf dist
	goreleaser
