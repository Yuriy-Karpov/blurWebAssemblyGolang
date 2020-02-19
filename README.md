# blurWebAssemblyGolang
this is an experiment to test the speed of an WebAssembly in a golang

Use this command to compile golang in wasm
```
GOARCH=wasm GOOS=js go build -o main.wasm main.go
```

read more about the compilation in https://github.com/golang/go/wiki/WebAssembly 
