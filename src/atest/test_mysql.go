package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

)

func test() {
	//操作数据库代码
	//第一个参数是数据库驱动
	//第二个参数是链接数据库字符串
	//u:u@tcp(127.0.0.1:3306)
	conn, err := sql.Open("mysql", "root:Root@123@tcp(10.1.102.190:3306)/gxs_ty?charset=utf8")
	if err != nil {
		fmt.Println("链接错误")
		return
	}
	/*
	   //创建表
	   _, err = conn.Exec("create table user(name varchar(40) ,password varchar(40) )")
	   if err != nil {
	       beego.Error("创建表失败", err)
	   }
	*/
	/*
	   //插入数据
	   conn.Exec("insert into user values(?,?)","tom","000000")
	*/
	/*
	   //插入多条数据（预处理）
	   users:=[][]string{{"lilei","000000"},{"hanmeimei","000000"}}
	   stat,_:=conn.Prepare("insert into user values(?,?)")
	   for _,s:=range users{
	       stat.Exec(s[0],s[1])
	   }

	*/
	//查询多行数据
	var city City
	res, err := conn.Query("select * from city limit 100")
	for res.Next() {
		res.Scan(&city.id, &city.city_name,&city.parent_id)
		fmt.Println(city)
	}
	//关闭数据库连接
	conn.Close()

}

func main() {
	test()
}

type City struct {
	id string
	city_name string
	parent_id string
}