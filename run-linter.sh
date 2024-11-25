if ! which golangci-lint &> /dev/null
then
  echo "golangci-lint not found, trying to install it..."
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
fi
golangci-lint run --fix

echo "Lints checked and possible lints fixed successfully"