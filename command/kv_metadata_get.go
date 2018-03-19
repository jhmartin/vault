package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVMetadataGetCommand)(nil)
var _ cli.CommandAutocomplete = (*KVMetadataGetCommand)(nil)

type KVMetadataGetCommand struct {
	*BaseCommand
}

func (c *KVMetadataGetCommand) Synopsis() string {
	return "Retrieves key metadata from the KV store"
}

func (c *KVMetadataGetCommand) Help() string {
	helpText := `
Usage: vault kv metadata get [options] KEY

  Retrieves the metadata from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned.

      $ vault kv metadata get secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVMetadataGetCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *KVMetadataGetCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVMetadataGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataGetCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	path := sanitizePath(args[0])
	path, err = addPrefixToVKVPath(path, "metadata")
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := kvReadRequest(client, path, nil)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading %s: %s", path, err))
		return 2
	}
	if secret == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
