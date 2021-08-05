package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
)

type StandaloneRedisSequencer struct {
	client     *redis.Client
	metadata   utils.RedisMetadata
	biggerThan map[string]int64

	logger log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStandaloneRedisSequencer returns a new redis sequencer
func NewStandaloneRedisSequencer(logger log.ErrorLogger) *StandaloneRedisSequencer {
	s := &StandaloneRedisSequencer{
		logger: logger,
	}
	return s
}

/*
   1. exists and >= biggerThan, no operation required, return 0
   2. not exists or < biggthan, reset val, return 1
   3. lua script occur error, such as tonumer(string), return error
*/
const initScript = `
if  redis.call('exists', KEYS[1])==1 and tonumber(redis.call('get', KEYS[1])) >= tonumber(ARGV[1]) then
    return 0
else
     redis.call('set', KEYS[1],ARGV[1])
     return 1
end
`

func (s *StandaloneRedisSequencer) Init(config sequencer.Configuration) error {
	m, err := utils.ParseRedisMetadata(config.Properties)
	if err != nil {
		return err
	}
	//init
	s.metadata = m
	s.biggerThan = config.BiggerThan

	// construct client
	s.client = utils.NewRedisClient(m)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	//check biggerThan, initialize if not satisfied
	for k, needV := range s.biggerThan {
		if needV <= 0 {
			continue
		}

		eval := s.client.Eval(s.ctx, initScript, []string{k}, needV)
		err = eval.Err()
		//occur error,  such as value is string type
		if err != nil {
			return err
		}
		//As long as there is no error, the initialization is successful
		//It may be a reset value or it may be satisfied before
	}
	return nil
}

func (s *StandaloneRedisSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	incr := s.client.Incr(s.ctx, req.Key)

	err := incr.Err()
	if err != nil {
		return nil, err
	}

	return &sequencer.GetNextIdResponse{
		NextId: incr.Val(),
	}, nil
}

func (s *StandaloneRedisSequencer) GetSegment(req *sequencer.GetSegmentRequest) (bool, *sequencer.GetSegmentResponse, error) {

	// size=0 only check support
	if req.Size == 0 {
		return true, nil, nil
	}

	by := s.client.IncrBy(s.ctx, req.Key, int64(req.Size))
	err := by.Err()
	if err != nil {
		return true, nil, err
	}

	return true, &sequencer.GetSegmentResponse{
		From: by.Val() - int64(req.Size) + 1,
		To:   by.Val(),
	}, nil
}
func (s *StandaloneRedisSequencer) Close() error {
	s.cancel()
	return s.client.Close()
}
