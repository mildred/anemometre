help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build for Raspberry Pi
	GOARCH=arm GOARM=5 go build .

fetch.go: node_modules/whatwg-fetch/fetch.js
	echo "package main;" >$@
	echo 'const html_fetch string = `' >>$@
	cat $< | sed 's/`/` + "`" + `/g' >>$@
	echo '`' >>$@

