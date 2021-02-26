all:
	go-assets-builder html public -o assets.go
	go build