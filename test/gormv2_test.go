package test

import (
	"fmt"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/gorm_v2"
	_ "ginskeleton/bootstrap"
	"sync"
	"testing"
	"time"
)

type tb_users struct {
	Id        uint   `json:"id"  gorm:"primaryKey" `
	UserName  string `json:"user_name"`
	Age       uint8  `json:"age"`
	Addr      string `json:"addr"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (*tb_users) TableName() string {
	return "tb_users"
}

type tb_role struct {
	Id          uint   `json:"id"  gorm:"primaryKey" `
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (*tb_role) TableName() string {
	return "tb_role"
}

type tb_user_log struct {
	Id        int `gorm:"primaryKey" `
	UserId    int
	Ip        string
	LoginTime string
	Remark    string
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (*tb_user_log) TableName() string {
	return "tb_user_log"
}

type tb_code_list struct {
	Code           string
	Name           string
	CompanyName    string
	Concepts       string
	ConceptsDetail string
	Province       string
	City           string
	Status         uint8
	Remark         string
	CreatedAt      string
	UpdatedAt      string
}

func (*tb_code_list) TableName() string {
	return "tb_code_list"
}

func TestGormSelect(t *testing.T) {

	var users []tb_users
	var roles []tb_role

	result := variable.GormDbMysql.Select("id", "user_name", "phone", "email", "remark").Where("user_name  like ?", "%test%").Find(&users)
	if result.Error != nil {
		t.Errorf("单元测试失败，错误明细:%s\n", result.Error.Error())
	}
	fmt.Printf("tb_users表数据：%v\n", users)

	result = variable.GormDbMysql.Where("name  like ?", "%test%").Find(&roles)
	if result.Error != nil {
		t.Errorf("单元测试失败，错误明细:%s\n", result.Error.Error())
	}
	fmt.Printf("tb_roles表数据：%v\n", roles)
}

func TestGormInsert(t *testing.T) {
	var usrLog = &tb_user_log{
		UserId:    3,
		Ip:        "192.168.1.110",
		LoginTime: time.Now().Format("2006-01-02 15:04:05"),
		Remark:    "备注信息1028",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	result := variable.GormDbMysql.Create(usrLog)
	if result.RowsAffected < 0 {
		t.Error("新增失败，错误详情：", result.Error.Error())
	}

	result = variable.GormDbMysql.Select("user_id", "ip", "remark").Create(usrLog)
	if result.RowsAffected < 0 {
		t.Error("新增失败，错误详情：", result.Error.Error())
	}
}

func TestGormUpdate(t *testing.T) {
	var usrLog = tb_user_log{
		Id:        13,
		UserId:    3,
		Ip:        "127.0.0.1",
		LoginTime: "2008-08-08 08:08:08",
		Remark:    "整个结构体对应的字段全部更新",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	result := variable.GormDbMysql.Save(&usrLog)
	if result.RowsAffected < 0 {
		t.Error("update失败，错误详情：", result.Error.Error())
	}

	var relaValue = map[string]interface{}{
		"user_id":    66,
		"ip":         "192.168.6.66",
		"login_time": time.Now().Format("2006-01-02 15:04:05"),
		"remark":     "指定字段更新，备注信息",
	}

	result = variable.GormDbMysql.Model(&usrLog).Select("user_id", "ip", "login_time", "remark").Where("id=?", 13).Updates(relaValue)
	if result.RowsAffected < 0 {
		t.Error("update失败，错误详情：", result.Error.Error())
	}
}

func TestGormDelete(t *testing.T) {

	var key_primary_struct = tb_role{
		Id: 4,
	}

	result := variable.GormDbMysql.Delete(key_primary_struct)
	if result.RowsAffected < 0 {
		t.Error("delete失败，错误详情：", result.Error.Error())
	}

	result = variable.GormDbMysql.Delete(&tb_role{}, 5)
	if result.RowsAffected < 0 {
		t.Error("delete失败，错误详情：", result.Error.Error())
	}
}

func TestRawSql(t *testing.T) {

	var receive []tb_user_log
	variable.GormDbMysql.Raw("select * from   tb_user_log  where id>?", 0).Scan(&receive)
	fmt.Printf("%v\n", receive)

	result := variable.GormDbMysql.Exec("update tb_user_log  set  remark=?  where   id=?", "gorm原生sql执行修改操作", 17)
	if result.RowsAffected < 0 {
		t.Error("原生sql执行失败，错误详情：", result.Error.Error())
	}
}

func TestBatchInsertSql(t *testing.T) {

	sql := `
INSERT  INTO tb_auth_post_mount_has_menu_button(fr_auth_post_mount_has_menu_id,fr_auth_button_cn_en_id)
SELECT 91,4 FROM  DUAL  WHERE   NOT EXISTS(SELECT 1 FROM tb_auth_post_mount_has_menu_button a  WHERE  a.fr_auth_post_mount_has_menu_id=91 AND a.fr_auth_button_cn_en_id=4)
	`
	for i := 0; i < 100; i++ {
		result := variable.GormDbMysql.Exec(sql)
		if result.Error != nil {
			t.Error("原生sql执行失败，错误详情：", result.Error.Error())
		}
	}

}

func TestUseTime(t *testing.T) {

	var receives []tb_code_list
	var time1 = time.Now()
	for i := 0; i < 100; i++ {
		receives = make([]tb_code_list, 0)
		variable.GormDbMysql.Model(tb_code_list{}).Select("code", "name", "company_name", "concepts_detail", "province", "city", "remark", "status", "created_at", "updated_at").Where("id<=?", 3500).Find(&receives)

	}
	fmt.Printf("gorm数据遍历完毕：最后一次条数：%d\n", len(receives))

	fmt.Printf("本次耗时（毫秒）：%d\n", time.Now().Sub(time1).Milliseconds())

}

func TestCocurrent(t *testing.T) {

	var wg sync.WaitGroup

	var conNum = make(chan uint16, 128)
	wg.Add(1000)
	time1 := time.Now()
	for i := 1; i <= 1000; i++ {
		conNum <- 1
		go func() {
			defer func() {
				<-conNum
				wg.Done()
			}()
			var received []tb_code_list
			variable.GormDbMysql.Table("tb_code_list").Select("code", "name", "company_name", "province", "city", "remark", "status", "created_at", "updated_at").Where("id<=?", 3500).Find(&received)

		}()
	}
	wg.Wait()
	fmt.Printf("耗时（ms）:%d\n", time.Now().Sub(time1).Milliseconds())

}

func TestCustomeParamsConnMysql(t *testing.T) {

	type DataList struct {
		Id            int
		Username      string
		Last_login_ip string
		Status        int
	}

	confPrams := gorm_v2.ConfigParams{
		Write: struct {
			Host     string
			DataBase string
			Port     int
			Prefix   string
			User     string
			Pass     string
			Charset  string
		}{Host: "127.0.0.1", DataBase: "db_test", Port: 3306, Prefix: "tb_", User: "root", Pass: "DRsXT5ZJ6Oi55LPQ", Charset: "utf8"},
		Read: struct {
			Host     string
			DataBase string
			Port     int
			Prefix   string
			User     string
			Pass     string
			Charset  string
		}{Host: "127.0.0.1", DataBase: "db_stocks", Port: 3306, Prefix: "tb_", User: "root", Pass: "DRsXT5ZJ6Oi55LPQ", Charset: "utf8"}}

	var vDataList []DataList

	if gormDbMysql, err := gorm_v2.GetSqlDriver("mysql", 0, confPrams); err == nil {
		gormDbMysql.Raw("select id,username,status,last_login_ip from tb_users").Find(&vDataList)
		fmt.Printf("Read 数据库查询结果：%v\n", vDataList)
		res := gormDbMysql.Exec("update tb_users  set  real_name='Write数据库更新' where   id<=2 ")
		fmt.Printf("Write 数据库更新以后的影响行数：%d\n", res.RowsAffected)
	}
}

func TestSqlserver(t *testing.T) {
	var users []tb_users

	result := variable.GormDbSqlserver.Exec("update   tb_users  set  remark='update 操作 write数据库' where   id=?", 1)
	if result.Error != nil {
		t.Errorf("单元测试失败，错误明细:%s\n", result.Error.Error())
	}

	result = variable.GormDbSqlserver.Table("tb_users").Select("id", "user_name", "pass", "remark").Where("id > ?", 0).Find(&users)
	if result.Error != nil {
		t.Errorf("单元测试失败，错误明细：%s\n", result.Error.Error())
	}
	fmt.Printf("sqlserver数据查询结果：%v\n", users)
}

func TestPostgreSql(t *testing.T) {
	var users []tb_users

	result := variable.GormDbPostgreSql.Exec("update   web.tb_users  set  remark='update 操作 write数据库' where   id=?", 1)
	if result.Error != nil {
		t.Errorf("单元测试失败，错误明细:%s\n", result.Error.Error())
	}

	result = variable.GormDbPostgreSql.Table("web.tb_users").Select("").Select("id", "user_name", "age", "addr", "remark").Where("id > ?", 0).Find(&users)
	if result.Error != nil {
		t.Errorf("单元测试失败，错误明细：%s\n", result.Error.Error())
	}
	fmt.Printf("sqlserver数据查询结果：%v\n", users)
}
