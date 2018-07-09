#!/usr/bin/env fish

env GOOS=linux GOARCH=amd64 go build -v -o klein.amd64.linux github.com/kamaln7/klein/
env GOOS=darwin GOARCH=amd64 go build -v -o klein.amd64.darwin github.com/kamaln7/klein/
env GOOS=windows GOARCH=amd64 go build -v -o klein.amd64.windows.exe github.com/kamaln7/klein/
gsha256sum *amd64* > sha256sum.txt
