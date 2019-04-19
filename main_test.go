package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var proxyMappingTests = []struct {
	name              string
	reader            io.Reader
	expectedError     error
	expectedInterface *ProxyMapping
}{
	{
		name:              "Test load proxy mapping success",
		reader:            strings.NewReader("proxy /webhooks/* http://localhost:9001"),
		expectedError:     nil,
		expectedInterface: &ProxyMapping{RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"},
	},
}

func TestLoadProxyMapping(t *testing.T) {
	for _, tt := range proxyMappingTests {
		t.Run(tt.name, func(t *testing.T) {
			proxyMappings, err := loadProxyMappings(tt.reader)
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedInterface != nil {
				assert.NotEmpty(t, proxyMappings)
				assert.Equal(t, tt.expectedInterface, proxyMappings[0])
			}
		})
	}
}
