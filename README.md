# GoPack

GoPack is a tool to pack multiple image files into as little image files as possible.

## Building from source

With Go 1.18 installed:

```sh
git clone https://github.com/fabtjar/gopack.git
cd gopack
go build -o gopack cmd/gopack/main.go
```

## Usage

Generate a directory of images if you don't have any:
```sh
go build -o generate pkg/generate/generate.go
./generate 50
```

Pack all the images within the `image/` directory:
```sh
./gopack
```
