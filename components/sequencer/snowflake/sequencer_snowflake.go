package snowflake

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"mosn.io/pkg/log"
)

var (
	machineRoomIdBits   = 5                                                      // 机房id 五位数
	machineIdBits       = 5                                                      // 机器id
	sequencerBits       = 12                                                     // 每毫秒产生的id数量
	maxMachineRoomIdBit = int64(^(-1 << machineRoomIdBits))                      // 当前机房id最大值
	maxMachineIdBit     = int64(^(-1 << machineIdBits))                          // 当前机器id最大值
	moveMachineRoomBit  = sequencerBits + machineIdBits                          // 机房的左移动位
	moveMachineBit      = sequencerBits                                          // 机器的左移动位
	moveTimeStart       = sequencerBits + machineIdBits + machineRoomIdBits      // 时间戳的左移动位
	maxSequencerIdBit   = ^(-1 << sequencerBits)                                 // id 最大值
	machineRoomId       int64                                                    // 机房id
	machineId           int64                                                    // 机器id
	currentAQS                                                              = 0  // 并发控制
	lastTimestamp       int64                                               = -1 //上次id产生的时间，防止时间回拨
	once                sync.Once
	instance            *singleton
	lock                sync.Mutex
	logger              log.Logger
)

const (
	timeStart = 1655540588 // 系统开始时间
)

type singleton struct {
}

const (
	IP_PATTERN = "\\d{1,3}(\\.\\d{1,3}){3,5}$" //ip pattern
)

// GetInstanceWM 获取singleton对象 单例模式
func GetInstanceWM(machine, machineRoom int64) *singleton {
	if machine > maxMachineIdBit || machine < 0 {
		logger.Fatalf("the max lengths of the machine ID exceeds the maxMachineIdBit, maxMachineIdBit: %d ", maxMachineIdBit)
		return nil
	}
	if machineRoom > maxMachineRoomIdBit || machineRoom < 0 {
		logger.Fatalf("the max lengths of the machine room ID exceeds the maxMachineRoomIdBit, maxMachineIdBit: %d ", maxMachineIdBit)
		return nil
	}
	once.Do(func() {
		instance = &singleton{}
	})
	machineId = machine
	machineRoomId = machineRoom
	return instance
}

func (s *singleton) NextID() (int64, error) {
	lock.Lock()
	defer lock.Unlock()
	timestamp := time.Now().Unix() // 这儿注意下，golang是否有效率问题
	if timestamp < lastTimestamp {
		descTime := lastTimestamp - timestamp
		if descTime < 5 {
			// wait
			time.Sleep(time.Duration(descTime<<1) * time.Millisecond)
			timestamp = time.Now().Unix()
			if timestamp < lastTimestamp {
				return -1, fmt.Errorf("time moved backwards, refusing to generate id for %d milliseconds", descTime)
			}
		} else {
			return -1, fmt.Errorf("time moved backwards, refusing to generate id for %d milliseconds", descTime)
		}
	}
	if lastTimestamp == timestamp {
		// 在同一毫秒内
		currentAQS = (currentAQS + 1) & maxSequencerIdBit
		if currentAQS == 0 {
			timestamp = tilNextMillis(lastTimestamp)
		}
	} else {
		currentAQS = rand.Intn(2) + 1
	}
	lastTimestamp = timestamp
	nowTime := (timestamp - timeStart) << moveTimeStart
	machineRoom := machineRoomId << moveMachineRoomBit
	machine := machineId << moveMachineBit
	return nowTime | machineRoom | machine, nil
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().Unix()
	for {
		if timestamp <= lastTimestamp {
			break
		}
		timestamp = time.Now().Unix()
	}
	return timestamp
}
