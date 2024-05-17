package version

import (
	"testing"

	"github.com/flagship-io/abtasty-cli/utils"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	output, _ := utils.ExecuteCommand(VersionCmd)
	assert.Contains(t, output, "ABTasty CLI version:")
}
