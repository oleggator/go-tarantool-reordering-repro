build/gorepro: build go.sum main.go
	go build -tags=go_tarantool_ssl_disable -o build/gorepro gorepro

build:
	mkdir -p build
