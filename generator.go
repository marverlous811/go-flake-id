package flakeidgenerator

import (
	"errors"
	"sync"
	"time"
)

const (
	TIMESTAMP_LEN  = 48
	SEQUENCE_LEN   = 7
	MACHINE_ID_LEN = 8
	NANOSEC        = 1e-9
	MILISEC        = 1e-3
)

type IdPayload struct {
	Timestamp uint64
	MachineId uint8
	Sequence  uint8
}

type IdFlakeGeneratorSetting struct {
	StartTime       time.Time
	MachineIdGetter func() uint8
}

type IdFlakeGenerator struct {
	mutex       *sync.Mutex
	startTime   uint64
	elapsedTime uint64
	sequence    uint8
	machineId   uint8
}

func NewIdFlakeGenerator(st IdFlakeGeneratorSetting) *IdFlakeGenerator {
	sf := new(IdFlakeGenerator)
	sf.mutex = new(sync.Mutex)
	sf.sequence = 0
	sf.elapsedTime = 0

	if st.StartTime.After(time.Now()) {
		return nil
	} else if st.StartTime.IsZero() {
		sf.startTime = toFlaketimestamp(time.Date(2022, 6, 23, 0, 0, 0, 0, time.UTC))
	} else {
		sf.startTime = toFlaketimestamp(st.StartTime)
	}

	if st.MachineIdGetter == nil {
		sf.machineId = 1
	} else {
		sf.machineId = st.MachineIdGetter()
	}
	return sf
}

func (gen *IdFlakeGenerator) NextId() (uint64, error) {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()

	current := getElapsedTime(gen.startTime)
	if gen.elapsedTime < current {
		gen.elapsedTime = current
		gen.sequence = 0
	} else {
		seqMask := uint8((1 << SEQUENCE_LEN) - 1)
		gen.sequence = (gen.sequence + 1) & seqMask
		if gen.sequence == 0 {
			gen.elapsedTime++
			overtime := gen.elapsedTime - current
			time.Sleep(getSleepTime(overtime))
		}
	}

	payload := IdPayload{
		Timestamp: gen.elapsedTime,
		MachineId: gen.machineId,
		Sequence:  gen.sequence,
	}

	return idFromPayload(payload)
}

func getElapsedTime(startTime uint64) uint64 {
	return toFlaketimestamp(time.Now()) - startTime
}

func toFlaketimestamp(time time.Time) uint64 {
	return uint64(time.UTC().UnixMilli())
}

func getSleepTime(overtime uint64) time.Duration {
	return time.Duration(overtime*(MILISEC/NANOSEC)) - time.Duration(time.Now().UTC().UnixNano()%(MILISEC/NANOSEC))
}

func idFromPayload(payload IdPayload) (uint64, error) {
	if payload.Timestamp > (1<<(TIMESTAMP_LEN+SEQUENCE_LEN+MACHINE_ID_LEN+1))-1 {
		return 0, errors.New("OUT OF MEMORY")
	}

	return (0 << (TIMESTAMP_LEN + SEQUENCE_LEN + MACHINE_ID_LEN)) |
		(payload.Timestamp << (SEQUENCE_LEN + MACHINE_ID_LEN)) |
		uint64(payload.MachineId)<<SEQUENCE_LEN |
		uint64(payload.Sequence), nil
}

func idToPayload(id uint64) IdPayload {
	return IdPayload{
		Timestamp: getTimestamp(id),
		MachineId: getMachineId(id),
		Sequence:  getSequence(id),
	}
}

func getTimestamp(id uint64) uint64 {
	return id >> (SEQUENCE_LEN + MACHINE_ID_LEN)
}

func getMachineId(id uint64) uint8 {
	mask := ((1 << MACHINE_ID_LEN) - 1) << SEQUENCE_LEN
	return uint8((id & uint64(mask)) >> SEQUENCE_LEN)
}

func getSequence(id uint64) uint8 {
	mask := (1 << SEQUENCE_LEN) - 1
	return uint8((id & uint64(mask)))
}
