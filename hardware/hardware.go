package hardware

// Assert that Hardware implements the interface
var _ = HardwareInterface(&Hardware{})

type HardwareInterface interface {
	Populate() (map[string]string, error)
}
