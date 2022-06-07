This repository parses [TSH](http://www.poslarchive.com/tsh/doc/all.html) data files and computes standings for players over many tournaments. The point system used for these standings is configurable.

#### Example run script:

DB_MIGRATIONS_PATH=file://../../migrations/ DB_PATH=mgi.db TOURNEY_SCHEMA_PATH=../../cfg/pts_mgi.csv ./tshparser


# protoc

To generate pb files, run this in the base directory:

`protoc --twirp_out=rpc --go_out=rpc --go_opt=paths=source_relative --twirp_opt=paths=source_relative ./proto/tshparser.proto`

Make sure you have done

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
`go install github.com/twitchtv/twirp/protoc-gen-twirp@latest`
