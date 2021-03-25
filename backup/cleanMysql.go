package main

import (
	"fmt"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

func main() {
	DB,err:=sql.Open("mysql","root:Morigo2020@/timesUpStudent?charset=utf8")
	if err!=nil{
		fmt.Println(err)
	}
	_,err=DB.Exec("DELETE FROM timesUpStudents")
	if err!=nil{
		fmt.Println(err)	
	}
	DB.Close()
}


