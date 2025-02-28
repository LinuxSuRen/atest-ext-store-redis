package pkg

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"

	"github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/version"
)

type remoteserver struct {
	remote.UnimplementedLoaderServer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer() remote.LoaderServer {
	return &remoteserver{}
}

const keyPrefix = "atest-"

func getKey(name string) string {
	return keyPrefix + name
}

func (s *remoteserver) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	return
}
func (s *remoteserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	return
}
func (s *remoteserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	return
}
func (s *remoteserver) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	return
}
func (s *remoteserver) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	return
}
func (s *remoteserver) ListTestCases(ctx context.Context, suite *remote.TestSuite) (reply *server.TestCases, err error) {
	return
}
func (s *remoteserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	return
}
func (s *remoteserver) GetTestCase(ctx context.Context, input *server.TestCase) (reply *server.TestCase, err error) {
	return
}
func (s *remoteserver) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	return
}
func (s *remoteserver) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	return
}
func (s *remoteserver) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	reply = &server.ExtensionStatus{
		Version: version.GetVersion(),
		Ready:   true,
	}
	return
}
func (s *remoteserver) PProf(ctx context.Context, in *server.PProfRequest) (data *server.PProfData, err error) {
	log.Println("pprof", in.Name)

	data = &server.PProfData{
		Data: extension.LoadPProf(in.Name),
	}
	return
}

func (s *remoteserver) Query(ctx context.Context, query *server.DataQuery) (result *server.DataQueryResult, err error) {
	var cli *redis.Client
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}

	prefix := query.Key
	if prefix == "" {
		err = fmt.Errorf("prefix is required")
		return
	}

	var keys []string
	cursor := uint64(0)
	const count = 1000

	result = &server.DataQueryResult{}

	for {
		// Perform SCAN to find keys matching the prefix
		var partialKeys []string
		var err error
		partialKeys, cursor, err = cli.Scan(ctx, cursor, prefix, count).Result()
		if err != nil {
			return
		}

		// Collect found keys
		keys = append(keys, partialKeys...)

		for _, key := range keys {
			if val := cli.Get(ctx, key); val != nil {
				result.Data = append(result.Data, &server.Pair{
					Key:   key,
					Value: val.Val(),
				})
			}
		}

		// Break loop when all keys are collected
		if cursor == 0 {
			break
		}
	}
	return
}
