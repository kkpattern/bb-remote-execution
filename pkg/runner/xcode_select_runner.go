package runner

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	if xcodeVersionOverride, ok := request.EnvironmentVariables["XCODE_VERSION_OVERRIDE"]; ok {
		buildVIndex := strings.LastIndex(xcodeVersionOverride, ".")
		xcodeVersionShort := xcodeVersionOverride[:buildVIndex]
		developerDir := "/Applications/Xcode-"+xcodeVersionShort+".app/Contents/Developer/"
		if _, err := os.Stat(developerDir); !os.IsNotExist(err) {
			request.EnvironmentVariables["DEVELOPER_DIR"] = developerDir
			if appleSdkPlatform, ok := request.EnvironmentVariables["APPLE_SDK_PLATFORM"]; ok {
				appleSdkVersionOverride, found := request.EnvironmentVariables["APPLE_SDK_VERSION_OVERRIDE"]
				if !found {
					appleSdkVersionOverride = ""
				}
				sdkRoot := fmt.Sprintf("%sPlatforms/%s.platform/Developer/SDKs/%s%s.sdk", developerDir, appleSdkPlatform, appleSdkPlatform, appleSdkVersionOverride)
				if _, err := os.Stat(sdkRoot); !os.IsNotExist(err) {
					request.EnvironmentVariables["SDKROOT"] = sdkRoot
				}
			}
		}
	}
	response, err := r.base.Run(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
