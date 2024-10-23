build:
	GOOS=js GOARCH=wasm go build -ldflags "-s -w" -trimpath -o build/main.wasm ./cmd/wasmdemo

clean:
	rm -rf build

.PHONY: build clean