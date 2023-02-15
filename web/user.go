package web

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserPrivilege byte

var DB *gorm.DB

// 权限校验类
const (
	ADMIN UserPrivilege = iota
	NORMAL
)

type User struct {
	Id        int           `json:"id" gorm:"column:id;PRIMARY_KEY" "`
	UserName  string        `json:"userName" form:"userName" gorm:"column:username"`
	PassWord  string        `json:"passWord" form:"passWord" gorm:"column:password"`
	Privilege UserPrivilege `json:"privilege" form:"privilege"  gorm:"column:privilege"`
}

func (u User) TableName() string {
	return "users"
}

func NewUser(name, pass string, privilege UserPrivilege) *User {
	user := &User{
		0,
		name,
		pass,
		privilege,
	}
	return user
}

// NewAdmin 创建管理员
func NewAdmin(name, pass string) *User {
	return NewUser(name, pass, ADMIN)
}

// NewNormal 创建普通用户
func NewNormal(name, pass string) *User {
	return NewUser(name, pass, NORMAL)
}

func (u User) getPrivilege() UserPrivilege {
	return u.Privilege
}

func (u User) setPrivilege(p UserPrivilege) {
	u.Privilege = p
}

// 以下是数据库有关函数

// DSNFormat example: "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
const DSNFormat = "%s:%s@%s(%s)/%s?%s"

var (
	DSNUser   string = "root"
	DSNPass   string = "Cjp20030116"
	DSNProc   string = "tcp"
	DSNAddr   string = "127.0.0.1:3306"
	DSNDbName string = "BYR_Iptables"
	DSNDbArg  string = "charset=utf8mb4&parseTime=True&loc=Local"
)

// InitDB 初始化数据库
func InitDB() error {
	dsn := fmt.Sprintf(DSNFormat, DSNUser, DSNPass, DSNProc, DSNAddr, DSNDbName, DSNDbArg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB = db
	return err
}

// insertUser 返回 user的id
func insertUser(user *User) int {
	var err error
	user.PassWord, err = encodePass(user.PassWord)
	if err != nil {
		return -1
	}
	DB.Create(&user)
	return user.Id
}

// checkUser 检查用户是否存在，并检查密码是否正确
func checkUser(user User) User {
	internalUser := getUserByName(user.UserName)
	if internalUser == (User{}) {
		return User{}
	}
	if !comparePass(user.PassWord, internalUser.PassWord) {
		return User{}
	}
	return internalUser
}

// getUserByName 按照名字查找user, 如果没找到就返回空user
func getUserByName(name string) User {
	var user User
	err := DB.Where("username = ?", name).Take(&user).Error
	if err != nil {
		return User{}
	}
	return user
}

// getUserById 按照id查找user, 如果没找到就返回空user
func getUserById(id int) User {
	var user User
	err := DB.Where("id = ?", id).Take(&user).Error
	if err != nil {
		return User{}
	}
	return user
}


// 下面是一些辅助函数
// encodePass 用bcrypt加密下密码
func encodePass(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 0)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// comparePass 解密密码
// str1 是前台传来的， str2 是数据库中的
func comparePass(str1, str2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(str2), []byte(str1))
	return err == nil
}
