package main

import (
	"fmt"
	"io/ioutil"

	tap "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
	"github.com/golang/protobuf/proto"
)

func main() {
	path := "track_0.pb"
	// file, _ := os.Open(path)
	file, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	// defer file.Close()

	wrapper := tap.TraceWrapper{}
	proto.Unmarshal(file, &wrapper)
	trace := wrapper.GetSocketBufferedTrace()
	localAddress := trace.Connection.LocalAddress.GetSocketAddress().Address
	localPort := trace.Connection.LocalAddress.GetSocketAddress().PortSpecifier
	remoteAddress := trace.Connection.RemoteAddress.GetSocketAddress().Address
	remotePort := trace.Connection.RemoteAddress.GetSocketAddress().PortSpecifier
	fmt.Println("localAddress: ", localAddress)
	fmt.Println("localPort: ", localPort)
	fmt.Println("remoteAddress: ", remoteAddress)
	fmt.Println("remotePort: ", remotePort)

	// binary.Read(file, binary.LittleEndian, &wrapper)
	//fmt.Println(wrapper)
}
