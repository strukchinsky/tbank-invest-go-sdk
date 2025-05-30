# tbank-invest-go-sdk [![Go Report Card](https://goreportcard.com/badge/github.com/strukchinsky/tbank-invest-go-sdk)](https://goreportcard.com/report/github.com/strukchinsky/tbank-invest-go-sdk)

Simple repository, that contains [investAPI](https://github.com/RussianInvestments/investAPI) as git subtree,  generated go code from proto files and SDK for easier API interactions.

Currently we are using v1.x of API, so files are compiled in current directory. After next major update `Makefile` should be updated to generate files in `v2` directory.

`investAPI` can be updated with `make update` and proto files can be generated with `make generate`.

For extended utilities see `sdk` folder.
