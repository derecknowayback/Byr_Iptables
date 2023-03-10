package iptables

import (
	"github.com/coreos/go-iptables/iptables"
	"strings"
)

const (
	TABLE     = "nat"
	DnatChain = "PREROUTING"
	SnatChain = "POSTROUTING"
)

// iptables -t nat -A PREROUTING -d 10.11.0.0/16 -j DNAT --to-destination 192.168.0.0/16
// iptables -t nat -A POSTROUTING -s 192.168.0.0/16 -j SNAT --to-source 10.11.0.0/16

var Iptables *iptables.IPTables

type BadIpFormat struct {
	errMsg string
}

func (err BadIpFormat) Error() string{
	return err.errMsg
}

// InitIptables 初始化iptables
func InitIptables() error {
	ip4t, err := iptables.New() // ipv4版本, timeout 0
	Iptables = ip4t

	return err
}


func Dnat(ip1, ip2 string) error {
	// iptables -t nat -A PREROUTING -d 10.11.0.0/16 -j DNAT --to-destination 192.168.0.0/16
	// 需要我们去除 Ip2后缀
	ip2AndMask := strings.Split(ip2,"/")
	err := Iptables.Append(TABLE, DnatChain, getDnatArgs(ip1, ip2AndMask[0])...)
	return err
}

func deleteDnat(table, chain,ip1, ip2 string  ) error {
	err := Iptables.DeleteIfExists(table, chain, getDnatArgs(ip1, ip2)...)
	return err
}

// -d 10.11.0.0/16 -j DNAT --to-destination 192.168.0.0/16
func getDnatArgs(ip1, ip2 string) []string {
	args := []string{
		"-d",
		ip1,
		"-j",
		"DNAT",
		"--to-destination",
		ip2,
	}
	return args
}

func Snat(ip1, ip2 string) error {
	// iptables -t nat -A POSTROUTING -s 192.168.0.0/16 -j SNAT --to-source 10.11.0.0/16
	// 去除ip2后缀
	ip2AndMask := strings.Split(ip2,"/")
	err := Iptables.Append(TABLE, SnatChain, getSnatArgs(ip1, ip2AndMask[0])...)
	return err
}

func deleteSnat(table, chain,ip1, ip2 string) error {
	err := Iptables.DeleteIfExists(table, chain, getSnatArgs(ip1, ip2)...)
	return err
}

// -s 192.168.0.0/16 -j SNAT --to-source 10.11.0.0/16
func getSnatArgs(ip1, ip2 string) []string {
	args := []string{
		"-s",
		ip1,
		"-j",
		"SNAT",
		"--to-source",
		ip2,
	}
	return args
}



