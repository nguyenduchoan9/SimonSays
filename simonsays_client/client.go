package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	symonsayspb "github.com/nguyenduchoan9/simonsays/simonsayspb"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().Unix())

	// Dial server
	fmt.Println("Simon says hello client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	defer cc.Close()

	c := symonsayspb.NewSimonSaysClient(cc)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c symonsayspb.SimonSaysClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.Game(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*symonsayspb.Request{
		{
			Event: &symonsayspb.Request_Join{
				Join: &symonsayspb.Request_Player{
					Id: "Player 1:",
				},
			},
		},
	}
	waitc := make(chan struct{})
	startGame := make(chan struct{})
	// Two players is joining a game
	go func(gameChannel chan struct{}) {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(2000 * time.Millisecond)
		}
		close(gameChannel)
	}(startGame)
	// Two players start inputting the Color
	go func(channelGame chan struct{}) {
		<-channelGame
		for i := 1; i <= 20; i++ {
			// generate random nummber and send it to stream
			rnd := int32(rand.Intn(i))
			input := rnd % int32(len(symonsayspb.Color_value))
			req := symonsayspb.Request{
				Event: &symonsayspb.Request_Press{
					Press: toColor(input),
				},
			}
			fmt.Printf("Sending message: %v\n", &req)
			if err := stream.Send(&req); err != nil {
				log.Fatalf("can not send %v", err)
			}
			time.Sleep(time.Millisecond * 2000)
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}(startGame)
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: turn:%v, light:%v\n", res.GetTurn(), res.GetLightUp())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func toColor(number int32) symonsayspb.Color {
	switch symonsayspb.Color_name[number] {
	case symonsayspb.Color_BLUE.String():
		return symonsayspb.Color_BLUE
	case symonsayspb.Color_GREEN.String():
		return symonsayspb.Color_GREEN
	case symonsayspb.Color_RED.String():
		return symonsayspb.Color_RED
	case symonsayspb.Color_YELLOW.String():
		return symonsayspb.Color_YELLOW
	default:
		return symonsayspb.Color_BLUE
	}
}
