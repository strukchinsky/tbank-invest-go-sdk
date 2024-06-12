# tbank-invest-go-sdk

Simple repository, that contains [investAPI](https://github.com/RussianInvestments/investAPI) as git subtree and generated go code from proto files.

Currently we are using v1.x of API, so files are compiled in current directory. After next major update `Makefile` should be updated to generate files in `v2` directory.

`investAPI` can be updated with `make update` and proto files can be compiled with `make compile`.
