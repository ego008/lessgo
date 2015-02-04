package pagelet

import (
	"path/filepath"
	"strings"
)

type ConfigBase struct {
	HttpPort          int    // e.g. 8080
	HttpAddr          string // e.g. "", "127.0.0.1", "unix:/tmp/app.socket (if the port is zero)"
	UrlBasePath       string
	Module            []ConfigModule
	LocaleCookieKey   string
	SessionCookieKey  string
	InstanceId        string
	LessIdsServiceUrl string
}

type ConfigModule struct {
	Name      string
	ViewPaths []string
	Routes    []Route
}

func (c *ConfigBase) moduleInit(module string) {

	for _, v := range c.Module {

		if v.Name == module {
			return
		}
	}

	m := ConfigModule{
		Name:      module,
		ViewPaths: []string{},
		Routes:    []Route{},
	}

	c.Module = append(c.Module, m)
}

// RouteStaticAppend
func (c *ConfigBase) RouteStaticAppend(module, path, pathto string) {

	c.moduleInit(module)

	for i, v := range c.Module {

		if v.Name != module {
			continue
		}

		route := Route{
			Type: "static",
			Path: strings.Trim(path, "/"),
			Tree: []string{pathto},
		}
		//println(route)
		c.Module[i].Routes = append(v.Routes, route)

		break
	}
}

func (c *ConfigBase) RouteAppend(module, path string, args ...map[string]string) {

	c.moduleInit(module)

	path = strings.Trim(path, "/")
	tree := strings.Split(path, "/")
	if len(tree) < 1 {
		return
	}

	route := Route{
		Type:    "std",
		Path:    path,
		Tree:    tree,
		TreeLen: len(tree),
	}
	if len(args) == 1 {
		route.Params = args[0]
	}

	for i, v := range c.Module {

		if v.Name != module {
			continue
		}

		c.Module[i].Routes = append(v.Routes, route)

		break
	}
}

func (c *ConfigBase) ViewPath(module, path string) {

	path = filepath.Clean(path)

	c.moduleInit(module)

	for i, v := range c.Module {

		if v.Name != module {
			continue
		}

		gotPath := false
		for _, v2 := range v.ViewPaths {

			if v2 != path {
				continue
			}

			gotPath = true
		}

		if !gotPath {
			c.Module[i].ViewPaths = append(v.ViewPaths, path)
		}

		break
	}
}

func (c *ConfigBase) ViewFuncRegistry(name string, fn interface{}) {

	if _, ok := templateFuncs[name]; ok {
		return
	}

	templateFuncs[name] = fn
}

func (c *ConfigBase) I18n(file string) {
	i18nLoadMessages(file)
}
