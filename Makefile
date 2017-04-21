.PHONY: lint build vet fmt install clean

define info
	echo -e '\n\e[33m> msg \e[39m\n'
endef

default: build

build: vet lint
	@$(info:msg=building gitlab-cli)
	@go build -i -o gitlab-cli ./main.go && echo "built"

vet:
	@$(info:msg=examine source code)
	@go tool vet -shadowstrict . ./api && echo "passes"

lint:
	@$(info:msg=linting source code)
	@golint -min_confidence 0.1 -set_exit_status . ./api && echo "passes"

fmt:
	@$(info:msg=formatting source code)
	@go fmt . ./api && echo "formatted"

clean:
	@$(info:msg=cleaning)
	@go clean -i . ./api && echo "cleaned"

install: clean
	@$(info:msg=installing)
	@go install && echo "installed"
