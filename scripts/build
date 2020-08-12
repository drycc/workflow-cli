#!/usr/bin/env bash

build-tag(){
  CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o _dist/$1/drycc-$1-linux-386 drycc.go
  CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o _dist/$1/drycc-$1-darwin-386 drycc.go
  CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o _dist/$1/drycc-$1-windows-386 drycc.go
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o _dist/$1/drycc-$1-linux-amd64 drycc.go
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o _dist/$1/drycc-$1-darwin-amd64 drycc.go
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o _dist/$1/drycc-$1-windows-amd64 drycc.go
}

build-revision(){
  CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o _dist/$1/drycc-$1-linux-386 drycc.go
  CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o _dist/$1/drycc-$1-darwin-386 drycc.go
  CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o _dist/$1/drycc-$1-windows-386 drycc.go
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o _dist/$1/drycc-$1-linux-amd64 drycc.go
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o _dist/$1/drycc-$1-darwin-amd64 drycc.go
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o _dist/$1/drycc-$1-windows-amd64 drycc.go
}

echo "------------------$1 $2------------------"
"$1" "$2"
