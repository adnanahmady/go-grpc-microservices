# go-grpc-microservices

A simplified demo of three Go microservices communicating via gRPC,
implemented in ***monorepo*** code storage concept.

**Caution**: the main focus of this project is the gRPC and
communication between the tree microservices. 

# Index

* [Commands](#commands)
    * [Bring up the application](#bring-up-the-application)
    * [Bring down the application](#release-the-docker-resources)
* [Project Structure](#project-structure)

## Commands

This project uses docker to manage the dependencies,
so you can use the make commands to have the project
up and running as it should be.

### Bring up the application

you can up the docker environment using this command.
by running this command a built version will be created
for the `user-service` and you can connect to it though
the port `50051` using grpc.

```shell
make up
```

### Release the docker resources

You can bring down the running project docker environment
by this command.

```shell
make down
```

## Project Structure

You can take a brief view of current project strucutre in here.

```shell
.
├── cmd
│   └── userservice
│       └── main.go
├── config
│   ├── loader.go
│   └── structs.go
├── config.yaml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── providers.go
│   ├── user
│   │   └── server.go
│   ├── wire_gen.go
│   └── wire.go
├── LICENSE
├── logs # application log files, separated by day
│   └── app-2025-11-24.log
├── Makefile # app command helpers
├── pkg
│   ├── app
│   │   └── info.go
│   ├── applog
│   │   └── logger.go
│   ├── proto
│   │   ├── services_grpc.pb.go
│   │   ├── services.pb.go
│   │   └── services.proto
│   └── request
│       ├── context.go
│       └── middlewares.go
└── README.md

12 directories, 22 files
```
