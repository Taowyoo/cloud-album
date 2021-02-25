build:
	@go build -o cloud_album cli.go || { echo "Build failed"; exit 1; }
	@echo "Build successfully"
