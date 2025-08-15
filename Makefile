build:
	set GOARCH=wasm&& set GOOS=js&& go build -o web/app.wasm
	go build -o tmp.exe

run: build
	tmp.exe
