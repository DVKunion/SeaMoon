.PHONY: all
all: generate


generate:
	protoc --proto_path=./ --go_out=./ \
	--go_opt=Mtunnel.proto=./ \
	--go_opt=Mgost.proto=./gost \
	--go-grpc_out=./ \
	--go-grpc_opt=Mtunnel.proto=./ \
	--go-grpc_opt=Mgost.proto=./gost \
	./*.proto
