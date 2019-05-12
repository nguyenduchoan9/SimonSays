package main

import (
	"fmt"
	"io"
	"log"
	"net"

	symonsayspb "github.com/nguyenduchoan9/simonsays/simonsayspb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Simon says hello")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	symonsayspb.RegisterSimonSaysServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Game(stream symonsayspb.SimonSays_GameServer) error {
	fmt.Println("A player is challenging you...")

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Exit")
			return nil
		}
		if err != nil {
			fmt.Printf("The %v occurred", err)
			continue
		}

		// continue if data received from stream
		var sendErr error
		if join := req.GetJoin(); join != nil {
			sendErr = stream.Send(&symonsayspb.Response{
				Event: &symonsayspb.Response_Turn{
					Turn: symonsayspb.Response_BEGIN,
				},
			})
		} else {
			sendErr = stream.Send(&symonsayspb.Response{
				Event: &symonsayspb.Response_LightUp{
					LightUp: req.GetPress(),
				},
			})
		}
		req.GetPress()
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}
