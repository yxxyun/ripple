module github.com/yxxyun/ripple

go 1.13

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorilla/websocket v1.4.1
	github.com/rubblelabs/ripple v0.0.0-20190714134121-6dd7d15dd060
	github.com/willf/bitset v1.1.10
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
)

replace github.com/rubblelabs/ripple v0.0.0-20190714134121-6dd7d15dd060 => github.com/yxxyun/ripple v0.1.1
