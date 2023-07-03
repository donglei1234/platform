package fxsvcapp

import (
	"errors"
)

var (
	ErrTcpGrpcMuxMismatch       = errors.New("tcp grpc server and connection mux network do not match")
	ErrUnixGrpcMuxMismatch      = errors.New("unix grpc server and connection mux network do not match")
	ErrTcpHttpMuxMismatch       = errors.New("tcp http server and connection mux network do not match")
	ErrUnixHttpMuxMismatch      = errors.New("unix http server and connection mux network do not match")
	ErrMissingDocumentStoreURL  = errors.New("please set the DOCUMENT_STORE_URL environment variable")
	ErrMissingMemoryStoreURL    = errors.New("please set the MEMORY_STORE_URL environment variable")
	ErrInvalidDocumentStoreURL  = errors.New("invalid DOCUMENT_STORE_URL environment variable")
	ErrInvalidMemoryStoreURL    = errors.New("invalid MEMORY_STORE_URL environment variable")
	ErrMissingTracerServiceName = errors.New("missing tracer service name")
	ErrUnsupportedTracer        = errors.New("unsupported tracer type")
)
