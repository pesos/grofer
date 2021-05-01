package container

import (
	"testing"
)

func TestGetContainerMetrics(t *testing.T) {

	cid := "48e255249c6a"
	_, err := GetContainerMetrics(cid)
	if err != nil {
		return
	}

}
