package hardware

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jxsl13/osfacts/common"
)

type Hardware struct {
	sysctl map[string]string
}

func (h *Hardware) init() error {
	sysctl, err := common.GetSysctl("hw", "machdep", "kern")
	if err != nil {
		return err
	}
	h.sysctl = sysctl
	return nil
}

func (h *Hardware) Populate() (map[string]string, error) {
	err := h.init()
	if err != nil {
		return nil, err
	}

	hardwareFacts := make(map[string]string, 10)

	macFacts, err := h.getMacFacts()
	if err != nil {
		return nil, err
	}
	cpuFacts, err := h.getCpuFacts()
	if err != nil {
		return nil, err
	}
	memoryFacts, err := h.getMemoryFacts()
	if err != nil {
		return nil, err
	}
	uptimeFacts, err := h.getUptimeFacts()
	if err != nil {
		return nil, err
	}

	hardwareFacts = common.UpdateMap(hardwareFacts, macFacts)
	hardwareFacts = common.UpdateMap(hardwareFacts, cpuFacts)
	hardwareFacts = common.UpdateMap(hardwareFacts, memoryFacts)
	hardwareFacts = common.UpdateMap(hardwareFacts, uptimeFacts)

	return hardwareFacts, nil
}

func (h *Hardware) getSystemProfile() (map[string]string, error) {
	cmd, err := common.GetCommand("/usr/sbin/system_profiler", "SPHardwareDataType")
	if err != nil {
		return nil, err
	}

	lines, err := cmd.OutputLines()
	if err != nil {
		return nil, err
	}
	systemProfile := make(map[string]string)
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		switch len(parts) {
		case 0, 1:
			continue
		default:
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			systemProfile[key] = value
		}
	}
	return systemProfile, nil
}

func (h *Hardware) getMacFacts() (map[string]string, error) {
	macFacts := make(map[string]string, 4)
	cmd, err := common.GetCommand("sysctl", "hw.model")
	if err != nil {
		return nil, err
	}

	out, err := cmd.OutputLines()
	if err != nil {
		return nil, err
	}
	lineParts := strings.Split(out[len(out)-1], " ")
	switch len(lineParts) {
	case 0, 1:
	default:
		macFacts["model"] = lineParts[1]
		macFacts["product_name"] = lineParts[1]
		macFacts["osversion"] = h.sysctl["kern.osversion"]
		macFacts["osrevision"] = h.sysctl["kern.osrevision"]
	}
	return macFacts, nil
}

func (h *Hardware) getCpuFacts() (map[string]string, error) {
	cpuFacts := make(map[string]string, 3)

	if _, found := h.sysctl["machdep.cpu.brand_string"]; found {
		cpuFacts["processor"] = h.sysctl["machdep.cpu.brand_string"]
		cpuFacts["processor_cores"] = h.sysctl["machdep.cpu.core_count"]
	} else {
		// PowerPC
		systemProfile, err := h.getSystemProfile()
		if err != nil {
			return nil, err
		}
		cpuFacts["processor"] = fmt.Sprintf("%s @ %s", systemProfile["Processor Name"], systemProfile["Processor Speed"])
		cpuFacts["processor_cores"] = h.sysctl["hw.physicalcpu"]
	}

	if value, found := h.sysctl["hw.logicalcpu"]; found {
		cpuFacts["processor_vcpus"] = value
	} else if value, found := h.sysctl["hw.ncpu"]; found {
		cpuFacts["processor_vcpus"] = value
	} else {
		cpuFacts["processor_vcpus"] = ""
	}

	return cpuFacts, nil
}

func (h *Hardware) getMemoryFacts() (map[string]string, error) {
	defaultMemTotalMb, err := strconv.ParseInt(h.sysctl["hw.memsize"], 10, 64)
	if err != nil {
		return nil, err
	}
	memoryFacts := map[string]int64{
		"memtotal_mb": defaultMemTotalMb / 1024 / 1024,
		"memfree_mb":  0,
	}

	cmd, err := common.GetCommand("vm_stat")
	if err != nil {
		return nil, err
	}

	lines, err := cmd.OutputLines()
	if err != nil {
		return nil, err
	}

	memoryStats := make([][]string, 0, len(lines))
	for _, line := range lines {
		kvTuple := strings.SplitN(strings.TrimRight(line, "."), ":", 2)
		memoryStats = append(memoryStats, kvTuple)
	}

	memoryStatsMap := make(map[string]int64, len(memoryStats))

	for _, kvTuple := range memoryStats {
		switch len(kvTuple) {
		case 0, 1:
			continue
		default:
			key := kvTuple[0]
			value, err := strconv.ParseInt(strings.TrimSpace(kvTuple[1]), 10, 64)
			if err != nil {
				// string conversion failed
				continue
			}
			// set
			memoryStatsMap[key] = value
		}
	}

	totalUsed := int64(0)

	const (
		pageSize       = 4096
		pagesWiredDown = "Pages wired down"
		pagesActive    = "Pages active"
		pagesInactive  = "Pages inactive"
	)

	if value, found := memoryStatsMap[pagesWiredDown]; found {
		totalUsed += value * pageSize
	}
	if value, found := memoryStatsMap[pagesActive]; found {
		totalUsed += value * pageSize
	}
	if value, found := memoryStatsMap[pagesInactive]; found {
		totalUsed += value * pageSize
	}

	usedMb := totalUsed / 1024 / 1024
	memoryFacts["memfree_mb"] = memoryFacts["memtotal_mb"] - usedMb

	return common.Int64ToStringMap(memoryFacts), nil
}

func (h *Hardware) getUptimeFacts() (map[string]string, error) {
	// https://github.com/ansible/ansible/blob/47448f14583520ef0fc62108cea3f92323179810/lib/ansible/module_utils/facts/hardware/darwin.py#L134
	sysctl, err := common.GetSysctlCmd("-b", "kern.boottime")
	if err != nil {
		return nil, err
	}

	b, err := sysctl.Output()
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	seconds := int64(0)
	err = binary.Read(r, binary.LittleEndian, &seconds)
	if err != nil {
		return nil, err
	}
	start := time.Unix(int64(seconds), 0)

	return map[string]string{
		"uptime_seconds": strconv.FormatUint(uint64(time.Since(start).Seconds()), 10),
	}, nil
}
