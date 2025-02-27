package rewrite4a

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

// 初始化插件
func init() { plugin.Register("rewrite4a", setup) }

// setup 设置插件
func setup(c *caddy.Controller) error {
	r4a := Rewrite4A{}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		r4a.Next = next
		return r4a
	})

	return nil
}
