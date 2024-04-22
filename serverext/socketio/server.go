package socketio

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	port              int
	wsMaxConnNum      int64
	disConnectChan    chan *Client
	kickOfflineChan   chan *Client
	errChan           chan *Client
	clients           *clientMaps
	clientPool        sync.Pool
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	handshakeTimeout  time.Duration
	//hubServer         *Server
	validate *validator.Validate
}

func NewServer(opts ...Option) (*Server, error) {
	var config configs
	for _, o := range opts {
		o(&config)
	}
	v := validator.New()
	return &Server{
		port:             config.port,
		wsMaxConnNum:     config.maxConnNum,
		handshakeTimeout: config.handshakeTimeout,
		clientPool: sync.Pool{
			New: func() interface{} {
				return new(Client)
			},
		},
		disConnectChan:  make(chan *Client, 1000),
		kickOfflineChan: make(chan *Client, 1000),
		//kickHandlerChan: make(chan *kickHandler, 1000),
		validate: v,
		clients:  NewClientMaps(),
		//Compressor:      NewGzipCompressor(),
		//Encoder:         NewGobEncoder(),
	}, nil
}

func (s *Server) Run() {
	router := gin.New()
	options := &engineio.Options{
		PingTimeout:        0,
		PingInterval:       0,
		Transports:         nil,
		SessionIDGenerator: nil,
		RequestChecker: func(request *http.Request) (http.Header, error) {
			return nil, nil

		},
		ConnInitor: func(request *http.Request, conn engineio.Conn) {

		},
	}
	sioServer := socketio.NewServer(options)
	sioServer.OnConnect("/", func(conn socketio.Conn) error {
		//id, _ := strconv.ParseInt(conn.ID(), 10, 64)
		//queryValues, _ := url.ParseQuery(conn.URL().RawQuery)
		client := s.clientPool.Get().(*Client)
		client.SetValue(conn, 0, 0, 0)
		conn.SetContext(client)
		oldClient := s.clients.AddClient(client)
		if oldClient != nil {
			s.kickOfflineChan <- oldClient
		}
		return nil
	})

	sioServer.OnError("/", func(conn socketio.Conn, e error) {
		client := conn.Context().(*Client)
		s.errChan <- client
		log.Println("meet error:", e)
	})
	sioServer.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		client := conn.Context().(*Client)
		s.disConnectChan <- client
	})
	router.GET("/socket.io/*any", gin.WrapH(sioServer))
	router.POST("/socket.io/*any", gin.WrapH(sioServer))
	if err := router.Run(":" + strconv.Itoa(s.port)); err != nil {
		logx.Must(err)
	}

}
