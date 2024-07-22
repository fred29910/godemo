build:
	go build  -trimpath -o build/ ./cmd/...

clean:
	rm -rf build



.PHONY: build clean 