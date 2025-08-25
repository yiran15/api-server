package oauth_test

import (
	"net"
	"testing"
)

func TestDns(t *testing.T) {
	domain := "keycloak.qqlx.net"

	ips, err := net.LookupHost(domain)
	if err != nil {
		t.Fatalf("DNS 解析失败: %v", err)
	}

	for _, ip := range ips {
		t.Logf("解析结果: %s -> %s", domain, ip)
	}
}
