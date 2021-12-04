package hardware

import (
	"errors"
	"fmt"
	"os"

	"github.com/jxsl13/osfacts/common"
)

var (
	errUUIDNotFound = errors.New("partition uuid not found")
)

func getPartitionUUID(partName string) (string, error) {
	de, err := os.ReadDir("/dev/disk/by-uuid")
	if err != nil {
		return "", err
	}
	for _, uuid := range de {
		dev, err := common.RealPath("/dev/disk/by-uuid/" + uuid.Name())
		if err != nil {
			return "", fmt.Errorf("%w: %v", errUUIDNotFound, err)
		}
		if dev == ("/dev/" + partName) {
			return uuid.Name(), nil
		}
	}
	return "", errUUIDNotFound
}
