
name: Go-legacy-CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    - name: Install Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}
    - name: Checkout code
      uses: actions/checkout@v4
    - name: Test on 386
      run: |
        GOARCH=386 go test