package flakeidgenerator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NUMBER_ID_GENERATED = 10000

var payload = IdPayload{
	Timestamp: 1678160200841,
	MachineId: 12,
	Sequence:  5,
}

var generator = NewIdFlakeGenerator(
	IdFlakeGeneratorSetting{
		StartTime:       time.Time{},
		MachineIdGetter: func() uint8 { return 11 },
	},
)

func removeDuplicateItem(arr []uint64) []uint64 {
	allValue := make(map[uint64]bool)
	retval := []uint64{}
	for _, item := range arr {
		if _, value := allValue[item]; !value {
			allValue[item] = true
			retval = append(retval, item)
		}
	}
	return retval
}

func routineGenerateId(c chan uint64, i int) {
	id, err := generator.NextId()
	if err == nil {
		c <- id
	} else {
		println("generate error ", err.Error())
		c <- 0
	}

	if i == NUMBER_ID_GENERATED-1 {
		time.AfterFunc(200*time.Millisecond, func() { close(c) })
	}
}

func TestGenerateId(t *testing.T) {
	id, err := idFromPayload(payload)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, id, payload.Timestamp)
}

func TestReverseId(t *testing.T) {
	id, err := idFromPayload(payload)
	if err != nil {
		t.Errorf("generate error %v", err.Error())
		return
	}

	newPayload := idToPayload(id)
	assert.Equal(t, payload.Timestamp, newPayload.Timestamp)
	assert.Equal(t, payload.MachineId, newPayload.MachineId)
	assert.Equal(t, payload.Sequence, newPayload.Sequence)
}

func TestGenerateIdWithoutDuplicate(t *testing.T) {
	ids := []uint64{}

	for i := 0; i < NUMBER_ID_GENERATED; i++ {
		id, err := generator.NextId()
		if err == nil {
			ids = append(ids, id)
		}
	}

	assert.Equal(t, NUMBER_ID_GENERATED, len(ids))
	result := removeDuplicateItem(ids)
	assert.Equal(t, NUMBER_ID_GENERATED, len(result))
}

func TestConcurrentIdGenerate(t *testing.T) {
	c := make(chan uint64)
	ids := []uint64{}
	for i := 0; i < NUMBER_ID_GENERATED; i++ {
		go routineGenerateId(c, i)
	}

	for i := range c {
		if i > 0 {
			ids = append(ids, i)
		}
	}

	assert.Equal(t, NUMBER_ID_GENERATED, len(ids))
	result := removeDuplicateItem(ids)
	assert.Equal(t, NUMBER_ID_GENERATED, len(result))
}
