build:
	go build -ldflags "-s -w"  -trimpath -o build/ ./cmd/...

clean:
	rm -rf build cmd/geodemo/*.zip



.PHONY: build clean 