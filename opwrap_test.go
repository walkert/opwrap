package opwrap

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type envData struct {
	key string
	val string
}

func TestGetClient(t *testing.T) {
	tfile, err := ioutil.TempFile("", "tst")
	if err != nil {
		t.Fatalf("Unable to create temporary file: %v", err)
	}
	tname := tfile.Name()
	defer os.Remove(tname)
	tests := []struct {
		name      string
		env       []envData
		expectEnv bool
		wantErr   bool
		errString string
	}{
		{
			name: "BasicOPNoenv",
		},
		{
			name: "OneVarMissing",
			env: []envData{
				envData{"OPWRAP_CERT", "cert"},
				envData{"OPWRAP_CONF", "conf"},
			},
		},
		{
			name: "WithBadPort",
			env: []envData{
				envData{"OPWRAP_CERT", "cert"},
				envData{"OPWRAP_CONF", "conf"},
				envData{"OPWRAP_PORT", "port"},
			},
			expectEnv: true,
			wantErr:   true,
			errString: "unable to convert 'port' into an integer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the env first
			for _, key := range []string{"OPWRAP_PORT", "OPWRAP_CONF", "OPWRAP_CERT"} {
				os.Unsetenv(key)
			}
			for _, pair := range tt.env {
				os.Setenv(pair.key, pair.val)
			}
			envData, ok := getStashEnv()
			if !ok {
				if tt.expectEnv {
					t.Errorf("Expected to get back stash env vars but did not")
					return
				}
				return
			}
			_, err := getStashClient(envData)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Unexpected error: %v", err)
					return
				}
				if !strings.Contains(err.Error(), tt.errString) {
					t.Fatalf("Expected error string to contain '%s', but got: %v", tt.errString, err)
				}
			}
		})
	}
}
