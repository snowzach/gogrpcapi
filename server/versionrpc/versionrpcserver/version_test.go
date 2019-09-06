package versionrpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snowzach/gogrpcapi/conf"
)

func TestVersionGet(t *testing.T) {

	var vs versionRPCServer
	response, err := vs.Version(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, conf.GitVersion, response.Version)

}
