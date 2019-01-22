if ($args[0] -eq "build") {
  go build -a -installsuffix cgo -ldflags "-s -X github.com/drycc/workflow-cli/version.BuildVersion=$(git rev-parse --short HEAD)" -o drycc.exe .
} elseif ($args[0] -eq "test") {
  go test --cover --race -v $(glide novendor)
} elseif ($args[0] -eq "bootstrap") {
  glide install -u
} else {
  echo "Unknown command: '$args'"
  exit 1
}

exit $LASTEXITCODE
