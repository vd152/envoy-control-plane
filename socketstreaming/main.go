package main

import (
	"log"
	"net"
	"net/http"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var WS *websocket.Conn
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		if string(p) == "tcp changed" {
			log.Println("performing tcp change now")
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	WS, _ = upgrader.Upgrade(w, r, nil)

	log.Println("Client Connected")
	err := WS.WriteMessage(1, []byte("Hi Client!"))

	if err != nil {
		log.Println(err)
	}
	reader(WS)
}

type ALSServer struct {
}

func (a *ALSServer) StreamAccessLogs(logStream pb.AccessLogService_StreamAccessLogsServer) error {
	log.Println("Streaming access logs")
	for {
		log.Println("checking")
		data, err := logStream.Recv()
		if err != nil {
			return err
		}

		//err = WS.WriteMessage(1, []byte(data.String()))
		err = WS.WriteJSON(data)
		log.Println("done")
	}
}

func NewALSServer() *ALSServer {
	return &ALSServer{}
}
func startHTTP() {
	http.HandleFunc("/ws", wsEndpoint)
	http.ListenAndServe(":3000", nil)
}
func startTCP() {
	log.Println("Starting ALS Server")
	listener, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Failed to start listener on port 8080: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAccessLogServiceServer(grpcServer, NewALSServer())
	grpcServer.Serve(listener)
}
func main() {
	log.Println("Starting HTTP Server")
	go func() {
		startHTTP()
	}()

	startTCP()
}
