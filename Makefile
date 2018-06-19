account-proto:
	cd ./account && protoc account.proto --go_out=plugins=grpc:pb

catalog-proto:
	cd ./catalog && protoc catalog.proto --go_out=plugins=grpc:pb

order-proto:
	cd ./order && protoc order.proto --go_out=plugins=grpc:pb

graphql-gen:
	cd ./graphql/graph && gqlgen -schema ../schema.graphql