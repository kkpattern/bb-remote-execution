package runner

import (
	"context"
	"fmt"

	runner_pb "github.com/buildbarn/bb-remote-execution/pkg/proto/runner"

	"google.golang.org/protobuf/types/known/emptypb"
)

type xcodeSelectRunner struct {
	base                       runner_pb.RunnerServer
}

// NewXcodeSelectRunner creates a decorator of RunnerServer
// that will set DEVELOPER_DIR and SDKROOT env for requests.
func NewXcodeSelectRunner(base runner_pb.RunnerServer) runner_pb.RunnerServer {
	return &xcodeSelectRunner{
		base: base,
	}
}

func (r *xcodeSelectRunner) CheckReadiness(ctx context.Context, request *emptypb.Empty) (*emptypb.Empty, error) {
	return r.base.CheckReadiness(ctx, request)
}

func (r *xcodeSelectRunner) Run(ctx context.Context, request *runner_pb.RunRequest) (*runner_pb.RunResponse, error) {
	for name, value := range request.EnvironmentVariables {
		fmt.Print(name+"="+value)
	}
	response, err := r.base.Run(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
