build-Func:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bootstrap main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.
