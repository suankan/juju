// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelgeneration

import (
	"github.com/juju/errors"
	"gopkg.in/juju/names.v2"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/apiserver/params"
	coremodel "github.com/juju/juju/core/model"
)

// Client provides methods that the Juju client command uses to interact
// with models stored in the Juju Server.
type Client struct {
	base.ClientFacade
	facade base.FacadeCaller
}

// NewClient creates a new `Client` based on an existing authenticated API
// connection.
func NewClient(st base.APICallCloser) *Client {
	frontend, backend := base.NewClientFacade(st, "ModelGeneration")
	return &Client{ClientFacade: frontend, facade: backend}
}

// AddGeneration adds a model generation to the config.
func (c *Client) AddGeneration(model names.ModelTag) error {
	var result params.ErrorResult
	arg := params.Entity{Tag: model.String()}
	err := c.facade.FacadeCall("AddGeneration", arg, &result)
	if err != nil {
		return errors.Trace(err)
	}
	if result.Error != nil {
		return errors.Trace(result.Error)
	}
	return nil
}

// CancelGeneration adds a model generation to the config.
func (c *Client) CancelGeneration() error {
	var result params.ErrorResult
	err := c.facade.FacadeCall("CancelGeneration", nil, &result)
	if err != nil {
		return errors.Trace(err)
	}
	return result.Error
}

// SwitchGeneration adds a model generation to the config.
func (c *Client) SwitchGeneration(model names.ModelTag, version string) error {
	var result params.ErrorResult
	arg := params.GenerationVersionArg{Model: params.Entity{Tag: model.String()}}
	switch version {
	case "current":
		arg.Version = coremodel.GenerationCurrent
	case "next":
		arg.Version = coremodel.GenerationNext
	default:
		return errors.Trace(errors.New("version must be 'next' or 'current'"))
	}
	err := c.facade.FacadeCall("SwitchGeneration", arg, &result)
	if err != nil {
		return errors.Trace(err)
	}
	if result.Error != nil {
		return errors.Trace(result.Error)
	}
	return nil
}

// AdvanceGeneration advances a unit and/or applications to the 'next' generation.
func (c *Client) AdvanceGeneration(model names.ModelTag, entities []string) error {
	var results params.ErrorResults
	arg := params.AdvanceGenerationArg{Model: params.Entity{Tag: model.String()}}
	if len(entities) == 0 {
		return errors.Trace(errors.New("No units or applications to advance"))
	}
	for _, entity := range entities {
		switch {
		case names.IsValidApplication(entity):
			arg.Entities = append(arg.Entities,
				params.Entity{names.NewApplicationTag(entity).String()})
		case names.IsValidUnit(entity):
			arg.Entities = append(arg.Entities,
				params.Entity{names.NewUnitTag(entity).String()})
		default:
			return errors.Trace(errors.New("Must be application or unit"))
		}
	}
	err := c.facade.FacadeCall("AdvanceGeneration", arg, &results)
	if err != nil {
		return errors.Trace(err)
	}
	return results.Combine()
}