package converter

import "errors"

// Plugin can be used to extends functionality beyond what
// is offered by commonmark.
type Plugin interface {
	// The public name of the plugin, e.g. "strikethrough"
	Name() string

	// Init is called to initialize the plugin. It can be used to
	// *validate* the arguments and *register* the rules.
	Init(conv *Converter) error
}

// WithPlugins can be used to add additional functionality to the converter.
func WithPlugins(plugins ...Plugin) converterOption {
	return func(c *Converter) error {
		for _, plugin := range plugins {
			pluginName := plugin.Name()
			if pluginName == "" {
				return errors.New("the plugin has no name")
			}
			c.registeredPlugins = append(c.registeredPlugins, pluginName)

			err := plugin.Init(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
