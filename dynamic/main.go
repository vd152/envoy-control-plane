package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	envoy_config_tap_v3_pb "github.com/envoyproxy/go-control-plane/envoy/data/tap/v3"
	"github.com/golang/protobuf/jsonpb"
)

func init() {
}

func main() {
	flag.Parse()

	c := `
config_id: my_tap_id
tap_config:
  match_config:
   any_match: true
  output_config:
    streaming: true
    max_buffered_rx_bytes: 5000
    max_buffered_tx_bytes: 5000		
    sinks:
    - format: JSON_BODY_AS_BYTES
      streaming_admin: {}`

	body := []byte(c)

	resp, err := http.Post("http://localhost:9901/tap", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Status: [%s]\n", resp.Status)
	fmt.Println()
	var rb []byte
	reader := bufio.NewReader(resp.Body)
	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading streamed bytes %v", err)
		}
		fmt.Printf("%s", string(line))
		rb = append(rb, line...)
		var tw envoy_config_tap_v3_pb.TraceWrapper
		rdr := bytes.NewReader(rb)
		err = jsonpb.Unmarshal(rdr, &tw)
		if err == nil {
			bt := tw.GetHttpBufferedTrace()
			pbody := bt.Response.Body.GetAsBytes()
			log.Printf("Message %s\n", string(pbody))
			rb = []byte("")
		}

	}

}
