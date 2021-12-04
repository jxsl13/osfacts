package hardware

type Hardware struct {
	platform       string
	processor      string
	processorCores string
	mounts         []string
	devices        []string
}

func (h *Hardware) Platform() string {
	return ""
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

func (h *Hardware) Populate() error {

}
