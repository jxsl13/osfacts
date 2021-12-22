package hardware

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/jxsl13/osfacts/common"
)

var (
	originalMemoryFactKeys = []string{
		"Buffers", "Cached", "SwapCached", // additional new facts
	}

	memoryFactKeys = append(originalMemoryFactKeys, []string{
		"MemTotal", "SwapTotal", "MemFree", "SwapFree", // old facts
	}...)

	cpuKeys = []string{
		// model name is for Intel arch, Processor (mind the uppercase P)
		// works for some ARM devices, like the Sheevaplug.
		// 'ncpus active' is SPARC attribute
		"model name", "Processor", "vendor_id", "cpu", "Vendor", "processor",
	}

	// regex used against findmnt output to detect bind mounts
	bindMountRe = regexp.MustCompile(`.*]`)
	// regex used against mtab content to find entries that are bind mounts
	mtabBindMountRe = regexp.MustCompile(`.*bind.*"`)

	// regex used for replacing octal escape sequences
	octalEscapeRe = regexp.MustCompile(`\\[0-9]{3}`)
)

type Hardware struct {
	platform       string
	processor      string
	processorCores string
	mounts         []string
	devices        []string
}

func (h *Hardware) Platform() string {
	return "Linux"
}
func (h *Hardware) GetProcessor() string {
	return ""
}
func (h *Hardware) GetProcessorCores() string {
	return ""
}
func (h *Hardware) GetMounts() []string {
	return nil
}
func (h *Hardware) GetDevices() []string {
	return nil
}

func (h *Hardware) Populate() (map[string]string, error) {
	m, err := h.getMemoryFacts()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (h *Hardware) getMemoryFacts() (map[string]string, error) {
	memoryFacts := make(map[string]string, 12)

	lines, err := common.GetFileLines("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	memoryStats := make(map[string]int, 10)
	for _, line := range lines {
		data, err := common.SplitLineIntoAtLeast(line, ":", 2)
		if err != nil {
			continue
		}
		key := data[0]

		if common.InSlice(key, originalMemoryFactKeys) {
			valParts := strings.Split(strings.TrimSpace(data[1]), " ")
			if len(valParts) == 0 {
				continue
			}
			value := valParts[0]
			intVal, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			memoryFacts[fmt.Sprintf("%s_mb", strings.ToLower(key))] = strconv.Itoa(intVal / 1024)
		}

		if common.InSlice(key, memoryFactKeys) {
			valParts := strings.Split(strings.TrimSpace(data[1]), " ")
			if len(valParts) == 0 {
				continue
			}
			value := valParts[0]
			intVal, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			memoryStats[strings.ToLower(key)] = intVal / 1024
		}
	}

	if common.IntMapContainsKeys(memoryStats, "memtotal", "memfree") {
		memoryStats["real:used"] = memoryStats["memtotal"] - memoryStats["memfree"]
	}

	if common.IntMapContainsKeys(memoryStats, "cached", "memfree", "buffers") {
		memoryStats["nocache:free"] = memoryStats["cached"] + memoryStats["memfree"] + memoryStats["buffers"]
	}

	if common.IntMapContainsKeys(memoryStats, "memtotal", "nocache:free") {
		memoryStats["nocache:used"] = memoryStats["memtotal"] - memoryStats["nocache:free"]
	}

	if common.IntMapContainsKeys(memoryStats, "swaptotal", "swapfree") {
		memoryStats["swap:used"] = memoryStats["swaptotal"] - memoryStats["swapfree"]
	}

	// TODO: think of a better solution for nesting of facts
	memoryFacts["memory_mb.real.total"] = strconv.Itoa(memoryStats["memtotal"])
	memoryFacts["memory_mb.real.used"] = strconv.Itoa(memoryStats["real:used"])
	memoryFacts["memory_mb.real.free"] = strconv.Itoa(memoryStats["memfree"])

	memoryFacts["memory_mb.nocache.free"] = strconv.Itoa(memoryStats["nocache:free"])
	memoryFacts["memory_mb.nocache.used"] = strconv.Itoa(memoryStats["nocache:used"])

	memoryFacts["memory_mb.swap.total"] = strconv.Itoa(memoryStats["swaptotal"])
	memoryFacts["memory_mb.swap.free"] = strconv.Itoa(memoryStats["swapfree"])
	memoryFacts["memory_mb.swap.used"] = strconv.Itoa(memoryStats["swap:used"])
	memoryFacts["memory_mb.swap.cached"] = strconv.Itoa(memoryStats["swapcached"])

	return memoryFacts, nil
}

func (h *Hardware) getCpuFacts() (map[string]string, error) {
	cpuFacts := make(map[string]string, 5)

	var (
		i                   = 0
		vendorIdOccurrence  = 0
		modelNameOccurrence = 0
		processorOccurence  = 0

		physId  = ""
		coreId  = ""
		sockets = make(map[string]int)
		cores   = make(map[string]int)

		xen         = false
		xenParaVirt = false
	)

	if common.Exists("/proc/xen") {
		xen = true
	} else {
		line, err := common.GetFileFirstLine("/sys/hypervisor/type")
		if err == nil && strings.TrimSpace(line) == "xen" {
			xen = true
		}
	}

	processors := make([]string, 0, 1)

	lines, err := common.GetFileLines("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		data, err := common.SplitLineIntoAtLeast(line, ":", 2)
		if err != nil {
			continue
		}
		key := strings.TrimSpace(data[0])
		value := strings.TrimSpace(data[1])

		if xen && key == "flags" && !strings.Contains(value, "vme") {
			// Check for vme cpu flag, Xen paravirt does not expose this.
			// Need to detect Xen paravirt because it exposes cpuinfo
			// differently than Xen HVM or KVM and causes reporting of
			// only a single cpu core.
			xenParaVirt = true
		}

		if common.InSlice(key, cpuKeys) {
			processors = append(processors, value)
			switch key {
			case "vendor_id":
				vendorIdOccurrence++
			case "model name":
				modelNameOccurrence++
			case "processor":
				processorOccurence++
			}
			i++
		} else if key == "physical id" {
			physId = value
			if _, found := sockets[physId]; !found {
				sockets[physId] = 1
			}
		} else if key == "core id" {
			coreId = value
			if _, found := cores[coreId]; !found {
				cores[coreId] = 1
			}
		} else if key == "cpu cores" {
			iVal, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			sockets[physId] = iVal
		} else if key == "siblings" {
			iVal, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			cores[coreId] = iVal
		} else if key == "# processors" {
			cpuFacts["processor_cores"] = value
		} else if key == "ncpus active" {
			iVal, err := strconv.Atoi(value)
			if err != nil {
				continue
			}
			i = iVal
		}
	}

	// Skip for platforms without vendor_id/model_name in cpuinfo (e.g ppc64le)
	if vendorIdOccurrence > 0 && vendorIdOccurrence == modelNameOccurrence {
		i = vendorIdOccurrence
	}

	// aarch seemingly not supported by Go
	if common.StartsWithOneOf(runtime.GOARCH, "arm", "aarch", "ppc") {
		i = processorOccurence
	}

	if !common.StartsWithOneOf(runtime.GOARCH, "s390") {

		if xenParaVirt {
			strI := strconv.Itoa(i)
			cpuFacts["processor_count"] = strI
			cpuFacts["processor_cores"] = strI
			cpuFacts["processor_threads_per_core"] = strI
			cpuFacts["processor_vcpus"] = strI
		} else {

			var (
				processorCount          = 0
				processorCores          = 0
				processorThreadsPerCore = 0.0
			)

			if len(sockets) > 0 {
				processorCount = len(sockets)
			} else {
				processorCount = i
			}
			cpuFacts["processor_count"] = strconv.Itoa(processorCount)

			processorCores = common.FirstIntItemWithDefault(sockets, 1)
			cpuFacts["processor_cores"] = strconv.Itoa(processorCores)

			coreValues := common.IntMapToIntValues(cores)

			if len(coreValues) > 0 {
				processorThreadsPerCore = float64(coreValues[0]) / float64(processorCores)
			} else {
				processorThreadsPerCore = float64(1 / processorCores)
			}
			cpuFacts["processor_threads_per_core"] = common.Float64ToString(processorThreadsPerCore)
			cpuFacts["processor_vcpus"] = common.Float64ToString(processorThreadsPerCore * float64(processorCores) * float64(processorCount))

			// if the number of processors available to the module's
			// thread cannot be determined, the processor count
			// reported by /proc will be the default:

			nproc := runtime.NumCPU()
			cpuFacts["processor_nproc"] = strconv.Itoa(nproc)

		}
	}

	return cpuFacts, nil
}
