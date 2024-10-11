// Package utils
/**
 * @author ErSan
 * @email  mlt131220@163.com
 * @date   2024/7/27 17:16
 * @description snowflake 算法实现
 */
package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	// 2024年7月27号20:00 时刻的毫秒级时间戳，如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	startTime int64 = 1722081600000
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Snowflake, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Snowflake{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Snowflake) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := int64((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}

var (
	instance *Snowflake
	once     sync.Once
)

// GetInstance 单例模式-懒汉模式
func GetInstance() *Snowflake {
	once.Do(func() {
		instance, _ = NewWorker(1)
	})

	return instance
}

// GetSnowflakeId 获取ID
func GetSnowflakeId() int64 {
	return GetInstance().GetId()
}
