output =
ifeq ($(OS),Windows_NT)
	output = cloud_album.exe
else
	output = cloud_album
endif

build:
	@go build -o ${output} cli.go || { echo "Build failed"; exit 1; }
	@echo "Build successfully"

