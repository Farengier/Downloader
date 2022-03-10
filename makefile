build:
	@echo +++ build +++
	@go build -o bin/Downloader src/main.go
	@GOOS=windows go build -o bin/Downloader.exe src/main.go
	@echo Done