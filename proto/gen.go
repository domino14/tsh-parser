package proto

//go:generate protoc --twirp_out=../rpc --go_out=../rpc --go_opt=paths=source_relative --twirp_opt=paths=source_relative ./tshparser.proto ./auth.proto
