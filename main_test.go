package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadProxyMappingsSuccess(t *testing.T) {
	t.Run("Test load proxy mappings success", func(t *testing.T) {
		proxyMappings, err := loadProxyMappings(strings.NewReader("proxy /webhooks/* http://localhost:9001"))
		assert.NoError(t, err)

		assert.NotEmpty(t, proxyMappings)

		assert.Equal(t, &ProxyMapping{RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}, proxyMappings[0])
	})
}
