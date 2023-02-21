package main

import (
	pr "EFms2GRPS/proto"
	"context"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type Handler struct{}

type BodyReq struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func main() {
	e := echo.New()
	handler := NewHandler()
	e.POST("/getSum", handler.getSum)
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

func NewHandler() *Handler {
	return &Handler{}
}
