package opwrap

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/walkert/op"
	"github.com/walkert/stash/client"
)

const (
	envCert = "OPWRAP_CERT"
	envConf = "OPWRAP_CONF"
	envPort = "OPWRAP_PORT"
)

type stashEnv struct {
	name  string
	value string
}

type stashMap map[string]*stashEnv

var smap = stashMap{
	"cert": &stashEnv{name: envCert},
	"conf": &stashEnv{name: envConf},
	"port": &stashEnv{name: envPort},
}

func getStashEnv() (stashMap, bool) {
	var missing bool
	for _, env := range smap {
		env.value = os.Getenv(env.name)
		if env.value == "" {
			missing = true
		}
	}
	if missing {
		return nil, false
	}
	return smap, true
}

func getStashClient(envData stashMap) (*client.Client, error) {
	pval := envData["port"].value
	iport, err := strconv.Atoi(pval)
	if err != nil {
		return nil, fmt.Errorf("unable to convert '%s' into an integer", pval)
	}
	c, err := client.New(iport, envData["conf"].value, envData["cert"].value)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetOP returns a pointer to an op.Op object configured with a password (if possible)
func GetOP(opts ...op.Opt) (*op.Op, error) {
	// Check to see if the caller has cert/conf/port set in their environment.
	// If any are not set, just return a standard Op object.
	envData, ok := getStashEnv()
	if !ok {
		return op.New(opts...)
	}
	c, err := getStashClient(envData)
	if err != nil {
		return op.New(opts...)
	}
	// Attempt to read the vault password from the configured stash server and
	// include logic for entering the password if it isn't set or has expired.
	// If this isn't possible, just return a standard Op object.
	pass, err := c.GetPassword()
	if err != nil {
		if strings.Contains(err.Error(), "not set") {
			fmt.Fprintln(os.Stderr, "Password not currently set - please re-enter")
			c.SetPassword()
			pass, err = c.GetPassword()
			if err != nil {
				return op.New(opts...)
			}
		} else {
			return op.New(opts...)
		}
	}
	opts = append(opts, op.WithPassword(string(pass)))
	return op.New(opts...)
}
