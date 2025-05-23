package pkg

import (
	"context"
	"fmt"
	"github.com/linuxsuren/api-testing/pkg/testing"
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
	var client *redis.Client
	if client, err = s.getClient(ctx); err != nil {
		return
	}

	suites = &remote.TestSuites{
		Data: []*remote.TestSuite{},
	}

	for _, key := range client.Keys(ctx, keyPrefix+"*").Val() {
		var data []byte
		if data, err = client.Get(ctx, key).Bytes(); err == nil {
			var testSuite *testing.TestSuite
			if testSuite, err = testing.Parse(data); err == nil {
				suites.Data = append(suites.Data, remote.ConvertToGRPCTestSuite(testSuite))
			}
		}
	}
	return
}
func (s *remoteserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}

	var client *redis.Client
	if client, err = s.getClient(ctx); err != nil {
		return
	}

	var data []byte
	if data, err = testing.ToYAML(remote.ConvertToNormalTestSuite(testSuite)); err == nil {
		if err = client.Set(ctx, getKey(testSuite.Name), data, 0).Err(); err == nil {
			reply = &server.Empty{}
		}
	}
	return
}
func (s *remoteserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var client *redis.Client
	if client, err = s.getClient(ctx); err != nil {
		return
	}

	var data []byte
	if data, err = client.Get(ctx, getKey(suite.Name)).Bytes(); err == nil {
		var testSuite *testing.TestSuite
		if testSuite, err = testing.Parse(data); err == nil {
			reply = remote.ConvertToGRPCTestSuite(testSuite)
		}
	}
	return
}
func (s *remoteserver) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	return
}
func (s *remoteserver) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	var client *redis.Client
	if client, err = s.getClient(ctx); err != nil {
		return
	}

	if err = client.Del(ctx, getKey(suite.Name)).Err(); err == nil {
		reply = &server.Empty{}
	}
	return
}
func (s *remoteserver) ListTestCases(ctx context.Context, suite *remote.TestSuite) (reply *server.TestCases, err error) {
	return
}
func (s *remoteserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	var client *redis.Client
	if client, err = s.getClient(ctx); err != nil {
		return
	}
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
	}

	var client *redis.Client
	if client, err = s.getClient(ctx); err == nil {
		if _, err = client.Ping(ctx).Result(); err != nil {
			reply.Ready = false
			reply.Message = fmt.Sprintf("redis server is not ready, error: %s", err.Error())
		} else {
			reply.Ready = true
		}
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
