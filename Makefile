BINARY=terraform-provider-jenkins
export COMPOSE_FILE=./integration/docker-compose.yml

default: build

# Builds the provider and adds it to your GOPATH/bin folder.
build:
	go build -o ${BINARY}