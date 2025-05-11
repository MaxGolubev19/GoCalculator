package orchestrator

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/MaxGolubev19/GoCalculator/pkg/proto"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (o *Orchestrator) RunTaskServer() {
	addr := fmt.Sprintf("0.0.0.0:%s", o.config.GrpcPort)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", o.config.GrpcPort)

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, o)
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

func (o *Orchestrator) GetTask(ctx context.Context, _ *emptypb.Empty) (*pb.TaskResponse, error) {
	if len(o.tasks) == 0 {
		return nil, status.Errorf(codes.NotFound, "")
	}

	task := &o.tasks[0]
	o.tasks = o.tasks[1:]

	return &pb.TaskResponse{
		Task: task,
	}, nil
}

func (o *Orchestrator) SubmitTask(ctx context.Context, request *pb.TaskRequest) (*emptypb.Empty, error) {
	id := int(request.Id)

	if o.taskId <= id {
		return nil, status.Errorf(codes.NotFound, "")
	}

	if request.StatusCode != 200 {
		o.actions[id].IsError = true
		return nil, nil
	}

	o.actions[id].Value = request.Result
	o.actions[id].IsCalculated = true

	return nil, nil
}

func (o *Orchestrator) AddTask(action *schemas.Action) int {
	var operationTime int
	switch action.Operation {
	case schemas.AddOperation:
		operationTime = o.config.TimeAdditionMS
	case schemas.SubOperation:
		operationTime = o.config.TimeSubstractionMS
	case schemas.MulOperation:
		operationTime = o.config.TimeMultiplicationMS
	case schemas.DivOperation:
		operationTime = o.config.TimeDivisionsMS
	}

	o.muTasks.Lock()
	defer o.muTasks.Unlock()

	taskId := o.taskId
	o.taskId++
	o.actions[taskId] = action

	o.tasks = append(o.tasks, pb.Task{
		Id:            int32(taskId),
		Arg1:          action.Left.Value,
		Arg2:          action.Right.Value,
		Operation:     pb.Operation(action.Operation),
		OperationTime: int32(operationTime),
	})

	return taskId
}
