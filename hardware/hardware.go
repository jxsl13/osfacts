package hardware

import "github.com/jxsl13/osfacts/interfaces"

// Assert that Hardware implements the interface
var _ interfaces.Collector = &Hardware{}
