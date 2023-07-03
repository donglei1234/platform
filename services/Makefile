.PHONY: swagger auth iap matchmaking notification room storage leaderboard guild

auth:
	protoc -I. -I./auth/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./auth/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./auth/generated/grpc/go/ \
	./auth/api/*

iap:
	protoc -I. -I./iap/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./iap/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./iap/generated/grpc/go/ \
	./iap/api/*

matchmaking:
	protoc -I. -I./matchmaking/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./matchmaking/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./matchmaking/generated/grpc/go/ \
	./matchmaking/api/*

notification:
	protoc -I. -I./notification/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./notification/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./notification/generated/grpc/go/ \
	./notification/api/*

room:
	protoc -I. -I./room/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./room/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./room/generated/grpc/go/ \
	./room/api/*

storage:
	protoc -I. -I./storage/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./storage/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./storage/generated/grpc/go/ \
	./storage/api/*

leaderboard:
	protoc -I. -I./leaderboard/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./leaderboard/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./leaderboard/generated/grpc/go/ \
	./leaderboard/api/*

guild:
	protoc -I. -I./guild/api -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:./guild/generated/grpc/go/ \
	--grpc-gateway_out=logtostderr=true:./guild/generated/grpc/go/ \
	./guild/api/*

all:
	make auth
	make iap
	make matchmaking
	make notification
	make room
	make storage
	make guild
	make leaderboard

swagger:
	protoc -I. -I./auth/api -I./iap/api -I./matchmaking/api -I./notification/api -I./room/api -I./storage/api \
	-I./leaderboard/api -I./guild/api \
	-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--swagger_out=logtostderr=true,allow_merge:./swagger \
	./*/api/*.proto