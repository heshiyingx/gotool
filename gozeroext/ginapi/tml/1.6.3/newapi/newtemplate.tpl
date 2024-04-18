syntax = "v1"

type Request {
  Name string `path:"name,options=you|me"`
}

type Response {
  Message string `json:"message"`
}

service {{.name}}-api {
  @handler {{.handler}}Handler
  get /from/:name(Request) returns (Response)
}

// goctl rpc protoc  {{.name}}.proto --go_out=. --go-grpc_out=. --zrpc_out=.
