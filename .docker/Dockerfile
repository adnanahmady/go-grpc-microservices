FROM golang:1.24-bullseye

RUN apt-get update && apt-get install -y sudo \
    && useradd -m docker \
    && usermod -aG sudo docker \
    && echo "docker ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

ENV PB_REL "https://github.com/protocolbuffers/protobuf/releases" 
RUN sudo apt-get update && sudo apt-get install -y zip unzip curl \
    && curl -LO ${PB_REL}/download/v30.2/protoc-30.2-linux-x86_64.zip \
    && sudo unzip protoc-30.2-linux-x86_64.zip -d /usr/local \
    && rm protoc-30.2-linux-x86_64.zip

USER docker

RUN GOFLAGS="" go install github.com/google/wire/cmd/wire@latest \
    && GOFLAGS="" go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && GOFLAGS="" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH ${PATH}:$(go env GOPATH)/bin

