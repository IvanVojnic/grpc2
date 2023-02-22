package main

import (
	pr "EFms2GRPS/proto"
	"context"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"time"
)

type Handler struct{}

type BodyReq struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type BodyReqNums struct {
	Nums []int `json:"nums"`
}

func main() {
	e := echo.New()
	handler := NewHandler()
	e.POST("/getSum", handler.getSum)
	e.POST("/check", handler.checkNum)
	e.Logger.Fatal(e.Start(":8080"))
}

func (h *Handler) getSum(ctx echo.Context) error {
	var req BodyReq
	ctx.Bind(&req)

	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pr.NewMs1Client(conn)
	res, errGRPC := c.Add(context.Background(), &pr.AddRequest{X: int32(req.X), Y: int32(req.Y)})
	if errGRPC != nil {
		log.Fatal(errGRPC)
	}
	return ctx.JSON(http.StatusOK, res)
}

func (h *Handler) checkNum(ctx echo.Context) error {
	var req BodyReqNums
	ctx.Bind(&req)
	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pr.NewMs1Client(conn)
	stream, err := c.IsEvan(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, errRecv := stream.Recv()
			if errRecv == io.EOF {
				close(waitc)
				return
			}
			if errRecv != nil {
				log.Fatalf("Failed to receive a num : %v", errRecv)
			}
			log.Printf("Got message %s", in.IsEvan)
			time.Sleep(time.Second * 2)
		}
	}()
	for _, num := range req.Nums {
		reqGRPC := pr.IsEvenNumRequest{Num: int32(num)}
		if errSend := stream.Send(&reqGRPC); errSend != nil {
			log.Fatalf("Failed to send a num: %v", errSend)
		}
	}
	stream.CloseSend()
	<-waitc
	return ctx.String(http.StatusOK, "look to console")
}

func NewHandler() *Handler {
	return &Handler{}
}
