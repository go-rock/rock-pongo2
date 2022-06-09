package rockpongo2

import (
	"io"

	"github.com/flosch/pongo2"
	"github.com/go-rock/rock"
)

type ViewEngine struct {
	Config ViewConfig
}

type ViewConfig struct {
	ViewDir   string
	Extension string
	Box       *pongo2.TemplateSet
}

func Default() *ViewEngine {
	config := ViewConfig{
		ViewDir:   "./views2/",
		Extension: ".html",
	}
	return New(config)
}

func New(config ViewConfig) *ViewEngine {
	return &ViewEngine{
		Config: config,
	}
}

func (e *ViewEngine) SetViewDir(viewDir string) {
	e.Config.ViewDir = viewDir
}

func (e *ViewEngine) GetViewDir() string {
	return e.Config.ViewDir
}

func (e *ViewEngine) Name() string {
	return "pg"
}

func (e *ViewEngine) Ext() string {
	return e.Config.Extension
}

func (r *ViewEngine) ExecuteWriter(writer io.Writer, filename string, bindingData interface{}) error {
	data := bindingData
	filename = rock.EnsureTemplateName(filename, r)
	template, err := r.loadFile(filename)
	if err != nil {
		return err
	}
	content := convertContext(data)
	err = template.ExecuteWriter(content, writer)
	if err != nil {
		return err
	}

	return err
}

func convertContext(templateData interface{}) pongo2.Context {
	if templateData == nil {
		return nil
	}

	if contextData, isPongoContext := templateData.(pongo2.Context); isPongoContext {
		return contextData
	}

	if contextData, isContextViewData := templateData.(rock.M); isContextViewData {
		return pongo2.Context(contextData)
	}

	return templateData.(map[string]interface{})
}

func (r *ViewEngine) loadFile(name string) (*pongo2.Template, error) {
	box := r.Config.Box
	var err error
	var template *pongo2.Template
	if box != nil {
		template, err = box.FromCache(name)
		if err != nil {
			return nil, err
		}
	} else {
		template, err = pongo2.FromFile(r.Config.ViewDir + name)
		if err != nil {
			return nil, err
		}
	}
	return template, nil
}
