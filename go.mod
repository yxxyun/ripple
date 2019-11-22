module websockets

go 1.13

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorilla/websocket v1.4.1
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/rubblelabs/ripple v0.0.0-20190714134121-6dd7d15dd060
	github.com/willf/bitset v1.1.10
	golang.org/x/crypto v0.0.0-20191112222119-e1110fd1c708 // indirect
)

replace github.com/rubblelabs/ripple v0.0.0-20190714134121-6dd7d15dd060 => ./
