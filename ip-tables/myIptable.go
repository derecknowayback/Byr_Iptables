package ip_tables

// iptables -t nat -A PREROUTING -d 10.11.0.0/16 -j DNAT --to-destination 192.168.0.0/16
// iptables -t nat -A POSTROUTING -s 192.168.0.0/16 -j SNAT --to-source 10.11.0.0/16

// 插入一条规则
