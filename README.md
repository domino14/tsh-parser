This repository parses [TSH](http://www.poslarchive.com/tsh/doc/all.html) data files and computes standings for players over many tournaments. The point system used for these standings is configurable.

#### Example run script:

```
DB_MIGRATIONS_PATH=file://../../migrations/ DB_PATH=mgi.db TOURNEY_SCHEMA_PATH=../../cfg/pts_mgi.csv SECRET_KEY=foobar ./tshparser
```

# protoc

To generate pb files, run this in the base directory:

```
protoc --twirp_out=rpc --go_out=rpc --go_opt=paths=source_relative --twirp_opt=paths=source_relative ./proto/tshparser.proto ./proto/auth.proto
```

Make sure you have done

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
```
and that you have the `protoc` compiler.

### Why use protobuf/Twirp? Isn't it bloated, annoying, etc, especially for such a small project?

See this:

https://github.com/domino14/tsh-parser/blob/f703b8efe7bbd443f53af379ff784c259ba6e00b/pkg/server/server.go#L38-L272

(my attempt at writing an http server to handle this). BTW there are several bugs in there.

Basically, you get a great API for free. The only thing I slightly dislike about Twirp is that it requires POSTs for everything. But why be so dogmatic about REST?