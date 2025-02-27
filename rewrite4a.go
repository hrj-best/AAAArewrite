package rewrite4a

import (
	"context"
	"encoding/binary"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Rewrite4A 插件结构
type Rewrite4A struct {
	Next plugin.Handler
}

// ServeDNS 处理 DNS 查询
func (r4a Rewrite4A) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	// 解析客户端 IPv4
	clientIP := net.ParseIP(state.IP()).To4()
	if clientIP == nil {
		return dns.RcodeServerFailure, nil // 如果获取失败，直接返回
	}

	// 确保 r 不为空，并且包含 AAAA 记录
	if r == nil || len(r.Answer) == 0 {
		return dns.RcodeServerFailure, nil // 如果获取失败，直接返回
	}

	// 遍历 Answer 部分，找到 AAAA 记录
	for _, ans := range r.Answer {
		if aaaaRecord, ok := ans.(*dns.AAAA); ok {
			// 原始 IPv6 地址
			originalIPv6 := aaaaRecord.AAAA

			// 保留前 64 位
			newIPv6 := make([]byte, 16)
			copy(newIPv6[:8], originalIPv6[:8])

			// 转换 IPv4 为 32bit 并填充
			binary.BigEndian.PutUint32(newIPv6[8:12], binary.BigEndian.Uint32(clientIP))

			// 解析查询的域名，并转换为 ASCII 截断 32bit
			domain := strings.TrimSuffix(state.Name(), ".") // 去掉末尾的点
			asciiDomain := []byte(domain)

			// 填充后 32bit（最多保留 4 个字符）
			copy(newIPv6[12:], asciiDomain)

			// 修改 AAAA 记录
			aaaaRecord.AAAA = net.IP(newIPv6)
		}
	}

	// 直接返回修改后的响应
	w.WriteMsg(r)
	return dns.RcodeSuccess, nil
}

// Name 返回插件名称
func (r4a Rewrite4A) Name() string {
	return "rewrite4a"
}
