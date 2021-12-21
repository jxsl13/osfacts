package hardware

import (
	"sync"
	"testing"
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
		hardwareSingleton = h

	})
	return hardwareSingleton
}

func Test_getMemoryFacts(t *testing.T) {
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

func TestHardware_getCpuFacts(t *testing.T) {
	h := GetTestHardware(t)

	got, err := h.getCpuFacts()
	if err != nil {
		t.Fatalf("Hardware.getCpuFacts() error = %v", err)
		return
	}

	if len(got) < 2 {
		t.Fatalf("Hardware.getCpuFacts() error = %v : %v", err, got)
	}
}
