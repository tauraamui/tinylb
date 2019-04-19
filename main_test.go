package main

import (
	"errors"
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

	{
		name:              "Test load proxy mapping fail fail missing request uri field",
		reader:            strings.NewReader("proxy"),
		expectedError:     errors.New("config line 1, missing request uri field for proxy mapping"),
		expectedInterface: nil,
	},

	{
		name:              "Test load proxy mapping fail fail missing targert url field",
		reader:            strings.NewReader("proxy /webhooks/*"),
		expectedError:     errors.New("config line 1, missing target url field for proxy mapping"),
		expectedInterface: nil,
	},

	{
		name:              "Test load proxy mapping fail fail blank line",
		reader:            strings.NewReader(""),
		expectedError:     nil,
		expectedInterface: nil,
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
