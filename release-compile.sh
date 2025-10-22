#! /bin/bash

rm -rf *.tar.gz

printf "#\n#\n#\nCompiling trucki2prometheus for Go OS 'linux' & arch 'amd64'\n"

if ! GOOS=linux GOARCH=amd64 go build -a -v -o trucki2prometheus -ldflags "-X main.t2PromVersion=$(git rev-parse --short HEAD) -X main.buildDate=$(date +"%Y-%m-%dT%H:%M:%S")"; then
	printf "Failed to cross compile trucki2prometheus for Go OS 'linux' & arch 'amd64\n"
	exit 1
fi

if ! tar -cvf trucki2prometheus-amd64.tar.gz ./trucki2prometheus; then
	printf "Failed to create trucki2prometheus-amd64.tar.gz\n"
	exit 1
fi

if ! rm -f ./trucki2prometheus; then
	printf "Failed to remove amd64 truck2prometheus binary\n"
	exit 1
fi

printf "\n#\n#\n#\n\nCompiling trucki2prometheus for Go OS 'linux' & arch 'arm64'\n"

if ! GOOS=linux GOARCH=arm go build -a -v -o trucki2prometheus -ldflags "-X main.t2PromVersion=$(git rev-parse --short HEAD) -X main.buildDate=$(date +"%Y-%m-%dT%H:%M:%S")"; then
	printf "Failed to cross compile trucki2prometheus for Go OS 'linux' & arch 'arm\n"
	exit 1
fi

if ! tar -cvf trucki2prometheus-arm64.tar.gz ./trucki2prometheus; then
	printf "Failed to create trucki2prometheus-arm64.tar.gz\n"
	exit 1
fi

if ! rm -f ./trucki2prometheus; then
	printf "Failed to remove arm64 truck2prometheus binary\n"
	exit 1
fi


printf "\n#\n#\n#\nRelease complication succesful!\n"

ls -al | grep -e ".tar.gz"
