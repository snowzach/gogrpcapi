EXECUTABLE := api
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
MIGRATIONDIR := store/postgres/migrations
MIGRATIONS :=  $(wildcard ${MIGRATIONDIR}/*.sql)
TOOLS := ${GOPATH}/bin/go-bindata \
	${GOPATH}/bin/mockery \
	${GOPATH}/src/github.com/golang/protobuf/proto \
	${GOPATH}/bin/protoc-gen-go \
	${GOPATH}/bin/protoc-gen-grpc-gateway \
	${GOPATH}/bin/protoc-gen-swagger
export PROTOBUF_INCLUDES = -I. -I/usr/include -I$(shell go list -e -f '{{.Dir}}' .) -I$(shell go list -e -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway/runtime)/../third_party/googleapis
PROTOS := ./gogrpcapi/thing.pb.go \
	./server/rpc/thing.pb.gw.go \
	./server/rpc/version.pb.gw.go

.PHONY: default
default: ${EXECUTABLE}

# This is all the tools required to compile, test and handle protobufs
tools: ${TOOLS}

${GOPATH}/bin/go-bindata:
	GO111MODULE=off go get -u github.com/go-bindata/go-bindata/...

${GOPATH}/bin/mockery:
	go get github.com/vektra/mockery/cmd/mockery

${GOPATH}/src/github.com/golang/protobuf/proto:
	go get github.com/golang/protobuf/proto

${GOPATH}/bin/protoc-gen-go:
	go get github.com/golang/protobuf/protoc-gen-go

${GOPATH}/bin/protoc-gen-grpc-gateway:
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

${GOPATH}/bin/protoc-gen-swagger:
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

# Handle all grpc endpoint protobufs
%.pb.gw.go: %.proto
	protoc ${PROTOBUF_INCLUDES} --go_out=paths=source_relative,plugins=grpc:. --grpc-gateway_out=paths=source_relative,logtostderr=true:. --swagger_out=logtostderr=true:. $*.proto

# Handle any non-specific protobufs
%.pb.go: %.proto
	protoc ${PROTOBUF_INCLUDES} --go_out=paths=source_relative,plugins=grpc:. $*.proto

${MIGRATIONDIR}/bindata.go: ${MIGRATIONS}
	# Building bindata
	go-bindata -o ${MIGRATIONDIR}/bindata.go -prefix ${MIGRATIONDIR} -pkg migrations ${MIGRATIONDIR}/*.sql

.PHONY: mocks
mocks: tools
	mockery -dir ./gogrpcapi -name ThingStore

.PHONY: ${EXECUTABLE}
${EXECUTABLE}: tools ${PROTOS} ${MIGRATIONDIR}/bindata.go
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}

.PHONY: test
test: tools ${PROTOS} ${MIGRATIONDIR}/bindata.go mocks
	go test -cover ./...

.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI

.PHONY: relocate
relocate:
        @test ${TARGET} || ( echo ">> TARGET is not set. Use: make relocate TARGET=<target>"; exit 1 )
        $(eval ESCAPED_PACKAGENAME := $(shell echo "${PACKAGENAME}" | sed -e 's/[\/&]/\\&/g'))
        $(eval ESCAPED_TARGET := $(shell echo "${TARGET}" | sed -e 's/[\/&]/\\&/g'))
        # Renaming package ${PACKAGENAME} to ${TARGET}
        @grep -rlI '${PACKAGENAME}' * | xargs -i@ sed -i 's/${ESCAPED_PACKAGENAME}/${ESCAPED_TARGET}/g' @
        # Complete...
        # NOTE: This does not update the git config nor will it update any imports of the root directory of this project.
