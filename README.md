This repository parses [TSH](http://www.poslarchive.com/tsh/doc/all.html) data files and computes standings for players over many tournaments. The point system used for these standings is configurable.

A user interface is included that allows an operator to add/remove tournaments, add/remove player aliases, and view standings over a period of time. This UI is written in [Elm](https://elm-lang.org).

#### How to run:

```
docker-compose up
```

# protoc

To generate pb files, run this in the `proto` directory:

```
go generate
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

### Why Elm?

My package-lock.json file for several of my projects has literally tens of thousands of lines in it, with all sorts of different versions of the same package. These update constantly as all sorts of security issues are found, resulting in errors and incompatibilities. It is terrible. I want to try something new.