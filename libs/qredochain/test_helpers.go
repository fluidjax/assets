package qredochain

import (
	"testing"

	"github.com/qredo/assets/libs/logger"
	"github.com/stretchr/testify/assert"
)

func StartTestConnectionNode(t *testing.T) *NodeConnector {

	logger, err := logger.NewLogger("text", "info")
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, logger, "logger should not be nil")

	nc, err := NewNodeConnector("127.0.0.1:26657", "NODEID", logger)
	assert.NotNil(t, nc, "tmConnector should not be nil")
	assert.Nil(t, err, "Error should be nil")
	return nc
}
