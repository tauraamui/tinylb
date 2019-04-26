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
	expectedInterface []*ProxyMapping
}{
	{
		name:              "Test load proxy mapping success",
		reader:            strings.NewReader("proxy /webhooks/* http://localhost:9001"),
		expectedError:     nil,
		expectedInterface: []*ProxyMapping{&ProxyMapping{RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}},
	},

	{
		name:              "Test load multiple proxy mappings success",
		reader:            strings.NewReader("proxy /webhooks/* http://localhost:9001\nproxy /webhooks2/* http://localhost:9002"),
		expectedError:     nil,
		expectedInterface: []*ProxyMapping{&ProxyMapping{RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}, &ProxyMapping{RequestURI: "/webhooks2/*", TargetURL: "http://localhost:9002"}},
	},

	{
		name:              "Test load proxy mapping with domain context success",
		reader:            strings.NewReader("tacusci.com proxy /webhooks/* http://localhost:9001"),
		expectedError:     nil,
		expectedInterface: []*ProxyMapping{&ProxyMapping{DomainContext: "tacusci.com", RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}},
	},

	{
		name:              "Test load multiple proxy mappings with domain context success",
		reader:            strings.NewReader("tacusci.com proxy /webhooks/* http://localhost:9001\nplace.com proxy /cheese-cake http://localhost:9001"),
		expectedError:     nil,
		expectedInterface: []*ProxyMapping{&ProxyMapping{DomainContext: "tacusci.com", RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}, &ProxyMapping{DomainContext: "place.com", RequestURI: "/cheese-cake", TargetURL: "http://localhost:9001"}},
	},

	{
		name:              "Test load proxy mapping with domain context as localhost success",
		reader:            strings.NewReader("localhost proxy /webhooks/* http://localhost:9001"),
		expectedError:     nil,
		expectedInterface: []*ProxyMapping{&ProxyMapping{DomainContext: "localhost", RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}},
	},

	{
		name:              "Test load proxy mapping with domain context as localhost with port 8000 success",
		reader:            strings.NewReader("localhost:8000 proxy /webhooks/* http://localhost:9001"),
		expectedError:     nil,
		expectedInterface: []*ProxyMapping{&ProxyMapping{DomainContext: "localhost:8000", RequestURI: "/webhooks/*", TargetURL: "http://localhost:9001"}},
	},

	{
		name:              "Test load proxy mapping fail missing request uri field",
		reader:            strings.NewReader("proxy"),
		expectedError:     errors.New("config line 1, missing request uri field for proxy mapping"),
		expectedInterface: nil,
	},

	{
		name:              "Test load proxy mapping fail missing targert url field",
		reader:            strings.NewReader("proxy /webhooks/*"),
		expectedError:     errors.New("config line 1, missing target url field for proxy mapping"),
		expectedInterface: nil,
	},

	{
		name:              "Test load proxy mapping fail unknown command directive",
		reader:            strings.NewReader("nonesense /webhooks/* http://localhost:9001"),
		expectedError:     errors.New("config line 1, unknown directive nonesense"),
		expectedInterface: nil,
	},

	{
		name:              "Test load proxy mapping fail blank line",
		reader:            strings.NewReader(""),
		expectedError:     nil,
		expectedInterface: nil,
	},

	{
		name:              "Test load proxy mapping fail fail nil io.Reader",
		reader:            nil,
		expectedError:     errors.New("io.Reader instance is a nil pointer"),
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
				assert.Equal(t, tt.expectedInterface, proxyMappings)
			}
		})
	}
}
