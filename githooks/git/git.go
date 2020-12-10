package git

import (
	"os"
	cm "rycus86/githooks/common"
	strs "rycus86/githooks/strings"
	"strings"
)

// ConfigScope Defines the scope of a config file, such as local, global or system.
type ConfigScope string

// Available ConfigScope's
const (
	LocalScope  ConfigScope = "--local"
	GlobalScope ConfigScope = "--global"
	System      ConfigScope = "--system"
	Traverse    ConfigScope = ""
)

// Context defines the context to execute it commands
type Context struct {
	cm.CmdContext
}

// CtxC creates a git command execution context with
// working dir `cwd`.
func CtxC(cwd string) *Context {
	return &Context{cm.CmdContext{BaseCmd: "git", Cwd: cwd}}
}

// CtxCSanitized creates a git command execution context with
// working dir `cwd` and sanitized environement.
func CtxCSanitized(cwd string) *Context {
	return (&Context{cm.CmdContext{BaseCmd: "git", Cwd: cwd, Env: SanitizeEnv(os.Environ())}})
}

// Ctx creates a git command execution context
// with current working dir.
func Ctx() *Context {
	return CtxC("")
}

// CtxSanitized creates a git command execution context
// with current working dir and sanitized environement.
func CtxSanitized() *Context {
	return CtxCSanitized("")
}

// GetConfig gets a Git configuration values.
func (c *Context) GetConfig(key string, scope ConfigScope) string {
	var out string
	var err error
	if scope != Traverse {
		out, err = c.Get("config", string(scope), key)
	} else {
		out, err = c.Get("config", key)
	}
	if err == nil {
		return out
	}
	return ""
}

// getConfigWithArgs gets a Git configuration values.
func (c *Context) getConfigWithArgs(key string, scope ConfigScope, args ...string) string {
	var out string
	var err error
	if scope != Traverse {
		out, err = c.Get(append(append([]string{"config"}, args...), string(scope), key)...)
	} else {
		out, err = c.Get(append(append([]string{"config"}, args...), key)...)
	}
	if err == nil {
		return out
	}
	return ""
}

// GetConfigAll gets a all Git configuration values.
func (c *Context) GetConfigAll(key string, scope ConfigScope, args ...string) []string {
	return strs.SplitLines(c.getConfigWithArgs(key, scope, "--get-all"))
}

// GetConfigAllU gets a all Git configuration values unsplitted.
func (c *Context) GetConfigAllU(key string, scope ConfigScope, args ...string) string {
	return c.getConfigWithArgs(key, scope, "--get-all")
}

// SetConfig sets a Git configuration values.
func (c *Context) SetConfig(key string, value interface{}, scope ConfigScope) error {
	v := strs.Fmt("%v", value)

	if scope != Traverse {
		return c.Check("config", string(scope), key, v)
	}
	return c.Check("config", key, v)
}

// IsConfigSet tells if a git config is set.
func (c *Context) IsConfigSet(key string, scope ConfigScope) bool {
	var err error
	if scope != Traverse {
		err = c.Check("config", string(scope), key)
	} else {
		err = c.Check("config", key)
	}
	return err == nil
}

func SanitizeEnv(env []string) []string {
	return strs.Filter(env, func(s string) bool {
		return !strings.Contains(s, "GIT_DIR") &&
			!strings.Contains(s, "GIT_WORK_TREE")
	})
}
