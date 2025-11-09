package flag_sender

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestPlugin1 = `package main

import "fmt"

func main() {
	fmt.Println("test plugin 1")
}`

const TestPlugin2 = `package main

import (
	"context"
	"flag_sender/pkg/plugins"
)

type Client struct {
	url   string
	token string
}

var NewClient plugins.NewClientFunc = func(url, token string) plugins.IClient {
	return &Client{
		url:   url,
		token: token,
	}
}

func (c *Client) SendFlags(ctx context.Context, flags []string) (map[string]*plugins.FlagResult, error) {
	return nil, nil
}`

func setupPlugin(t *testing.T, plugin string) string {
	tmpDir, err := os.MkdirTemp("/tmp", "plugindir_")
	if err != nil {
		t.Fatalf("error making temp dir %s", err.Error())
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	pluginGoPath := path.Join(tmpDir, "main.go")
	err = os.WriteFile(pluginGoPath, []byte(plugin), 0700)
	if err != nil {
		t.Fatalf("error writing plugin")
	}

	pluginBinaryPath := path.Join(tmpDir, "plugin")
	err = exec.CommandContext(t.Context(), "go", "build", "-buildmode=plugin", "-o", pluginBinaryPath, pluginGoPath).Run()
	if err != nil {
		t.Fatalf("error building plugin: %s", err.Error())
	}

	return pluginBinaryPath
}

// TODO: fix happy path test
// func TestPlugin_HappyPath(t *testing.T) {
// 	pluginPath := setupPlugin(t, TestPlugin2)

// 	plugin, err := loadPlugin(pluginPath, "test-jury-addr", "test-token")
// 	require.NotNil(t, plugin)
// 	require.NoError(t, err)
// }

func TestError(t *testing.T) {
	testcases := []struct {
		name     string
		testCode string
	}{
		{
			name:     "code without constructor & interface for plugin",
			testCode: TestPlugin1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pluginPath := setupPlugin(t, tc.testCode)
			plugin, err := loadPlugin(pluginPath, "test-jury-addr", "test-token")
			assert.Nil(t, plugin)
			require.Error(t, err)
		})
	}
}
