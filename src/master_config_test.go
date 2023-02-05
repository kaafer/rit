package src

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationMasterConstructorConfig(t *testing.T) {
	a := assert.New(t)
	master := NewMaster()
	a.True(master.hubCounter == 0)
	a.True(master.clientCounter == 0)
}
