package juju

import (
	"launchpad.net/juju/go/environs"
	"launchpad.net/juju/go/state"
	"regexp"
)

var (
	ValidService = regexp.MustCompile("^[a-z][a-z0-9]*(-[a-z0-9]*[a-z][a-z0-9]*)*$")
	ValidUnit    = regexp.MustCompile("^[a-z][a-z0-9]*(-[a-z0-9]*[a-z][a-z0-9]*)*/[0-9]+$")
)

// Conn holds a connection to a juju.
type Conn struct {
	Environ environs.Environ
	state   *state.State
}

// NewConn returns a Conn pointing at the environName environment, or the
// default environment if not specified.
func NewConn(environName string) (*Conn, error) {
	environs, err := environs.ReadEnvirons("")
	if err != nil {
		return nil, err
	}
	environ, err := environs.Open(environName)
	if err != nil {
		return nil, err
	}
	return &Conn{Environ: environ}, nil
}

// Bootstrap initializes the Conn's environment and makes it ready to deploy
// services.
func (c *Conn) Bootstrap(uploadTools bool) error {
	return c.Environ.Bootstrap(uploadTools)
}

// Destroy destroys the Conn's environment and all its instances.
func (c *Conn) Destroy() error {
	return c.Environ.Destroy(nil)
}

// State returns the conn's State.
func (c *Conn) State() (*state.State, error) {
	if c.state == nil {
		info, err := c.Environ.StateInfo()
		if err != nil {
			return nil, err
		}
		st, err := state.Open(info)
		if err != nil {
			return nil, err
		}
		c.state = st
	}
	return c.state, nil
}
