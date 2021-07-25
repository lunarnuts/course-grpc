PROTOC_VERSION=3.17.3
PROTOC_ARCH=linux-x86_64
PROTOC_DIR=https://github.com/protocolbuffers/protobuf/releases/download
BINDIR=~/bin
TMPDIR=/tmp/protoc

proto-gen:
	echo $(PATH)
	~/bin/protoc --proto_path=./common/proto \
--go_opt=paths=source_relative \
--go-grpc_opt=paths=source_relative \
--go_out=common/transport \
--go-grpc_out=common/transport \
messages.proto

download-protoc:
	curl -LO $(PROTOC_DIR)/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_ARCH).zip

unzip-protoc:
	unzip protoc-$(PROTOC_VERSION)-$(PROTOC_ARCH).zip -d $(TMPDIR) && cp -a $(TMPDIR)/bin/protoc  $(BINDIR)

install-go-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26

install-go-protoc-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

install-dep: download-protoc unzip-protoc install-go-protoc install-go-protoc-grpc

