package hardware

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	hardwareSingleton *Hardware
	hardwareOnce      = sync.Once{}
)

func GetTestHardware(t *testing.T) *Hardware {
	if hardwareSingleton != nil {
		return hardwareSingleton
	}
	hardwareOnce.Do(func() {
		h := &Hardware{}
		err := h.init()
		if err != nil {
			t.Fatal(err)
		}
		hardwareSingleton = h

	})
	return hardwareSingleton
}

func TestHardware_getSystemProfile(t *testing.T) {
	h := GetTestHardware(t)
	got, err := h.getSystemProfile()
	if err != nil {
		t.Fatalf("Hardware.getSystemProfile() error = %v", err)
		return
	}

	if len(got) == 0 {
		t.Fatalf("Hardware.getSystemProfile() len == 0")
	}
}

func TestHardware_getMacFacts(t *testing.T) {
	h := GetTestHardware(t)
	assert := assert.New(t)

	macFacts, err := h.getMacFacts()
	if err != nil {
		t.Fatalf("Hardware.getMacFacts() error = %v", err)
		return
	}

	assert.NotEmpty(macFacts["model"])
	assert.NotEmpty(macFacts["product_name"])
	assert.NotEmpty(macFacts["osversion"])
	assert.NotEmpty(macFacts["osrevision"])
}

func TestHardware_getCpuFacts(t *testing.T) {
	h := GetTestHardware(t)
	assert := assert.New(t)
	cpuFacts, err := h.getCpuFacts()
	if err != nil {
		t.Fatalf("Hardware.getCpuFacts() error = %v", err)
		return
	}

	assert.NotEmpty(cpuFacts["processor"])
	assert.NotEmpty(cpuFacts["processor_cores"])
	assert.NotEmpty(cpuFacts["processor_vcpus"])
}

func TestHardware_getMemoryFacts(t *testing.T) {
	h := GetTestHardware(t)

	got, err := h.getMemoryFacts()
	if err != nil {
		t.Fatalf("Hardware.getMemoryFacts() error = %v", err)
		return
	}

	if len(got) < 2 {
		t.Fatalf("Hardware.getMemoryFacts() error = %v : %v", err, got)
	}
}

func TestHardware_getUptimeFacts(t *testing.T) {
	h := GetTestHardware(t)
	assert := assert.New(t)
	got, err := h.getUptimeFacts()
	if err != nil {
		t.Fatalf("Hardware.getUptimeFacts() error = %v", err)
		return
	}
	assert.NotEmpty(got)

	for _, v := range got {
		i, err := strconv.ParseInt(v, 10, 64)
		assert.NoError(err)
		assert.Less(i, int64(60*60*24*365))
	}
}

func TestHardware_Populate(t *testing.T) {
	h := GetTestHardware(t)
	assert := assert.New(t)
	facts, err := h.Populate()
	if err != nil {
		t.Fatalf("Hardware.Populate() error = %v", err)
		return
	}

	assert.NotEmpty(facts["processor_cores"])
	assert.NotEmpty(facts["processor_vcpus"])
	assert.NotEmpty(facts["memfree_mb"])
	assert.NotEmpty(facts["memtotal_mb"])
	assert.NotEmpty(facts["uptime_seconds"])
}
