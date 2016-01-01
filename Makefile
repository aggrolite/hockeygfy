os=freebsd
arch=386

all:
	GOOS=$(os) GOARCH=$(arch) go build
