package web

type TrustedIp struct {
	id int `gorm:"column:id"`
	eip string `gorm:"column:eip"`
}

func (i TrustedIp) TableName() string {
	return "trustedip"
}

// 测试 ip在不在数据库内
func isIpTrusted(ip string) bool {
	var internalIp TrustedIp
	err := DB.Where("eip = ",ip).Take(&internalIp).Error
	if err != nil{
		return false
	}
	return true
}
