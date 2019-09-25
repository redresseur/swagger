package template

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestOutputAuthorityCodeWithTemplate(t *testing.T) {
	TestInterfaceComplete(t)
	assert.NoError(t, OutputAuthorityCodeWithTemplate(os.Stdout, "./author.tpl"))
	//assert.NoError(t, OutputAuthorityCode(os.Stdout))
}