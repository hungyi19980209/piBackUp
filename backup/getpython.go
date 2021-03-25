package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
	"encoding/base64"
	"io/ioutil"
	"os/exec"
	"net/http"
	"time"
	"github.com/hegedustibor/htgo-tts"//TTS
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)
func takepic(name string) {
	fmt.Print("taking picture\n")
	cmd := exec.Command("/usr/bin/fswebcam", name)
	//創建獲取命令輸出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}
	//執行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return
	}
	//讀取所有輸出
	bytes, err := ioutil.ReadAll(stdout)

	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error())
		return
	}
	fmt.Printf("take picture successful\n%s", bytes)

}

func base_64(name string) {
	fmt.Print("turn to BASE64\n")
	ff, _ := os.Open(name)
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, _ := ff.Read(sourcebuffer)
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	newname := name + ".json"
	ioutil.WriteFile(newname, []byte(sourcestring), 0667)

}

func sendfile(name string) {
	fmt.Print("sending json\n")
	newname := name + ".json"
	file, err := os.Open(newname)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	sentUrl := "http://40.74.116.89:5000/upload/" + newname
	res, err := http.Post(sentUrl, "binary/octet-stream", file)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	message, _ := ioutil.ReadAll(res.Body)
	fmt.Printf(string(message))
}

func TTS() {
	//connect
	time.Sleep(30 * time.Second)
    speech := htgotts.Speech{Folder: "audio", Language: "en"}
	speech.Speak("System Ready")
	DB,err:=sql.Open("mysql","root:Morigo2020@/timesUpStudent?charset=utf8")
	if err != nil {
			fmt.Println(err)
	}
	for {
		 
		conn, err := net.Dial("tcp", "40.74.116.89:5059") //163.13.127.178
		if err!=nil{
				fmt.Println(err)
		}
		if err == nil {
			var text string
			text = "School251001" //SchoolID  !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Need to change every new school
			fmt.Fprint(conn, text+"\n")
			message, _ := bufio.NewReader(conn).ReadString('\n')

			if message != "nobody\n" {
				if message != "" {
						if err!=nil{
							fmt.Println(err)
						}
						fmt.Println("Front3:"+message[:3])
                        if(message[:3]=="DEL"){
							_,err:=DB.Exec("DELETE FROM timesUpStudents WHERE studentID = ?",message[3:13])
							if err!=nil{
								fmt.Println(err)	
							}
						}else if(message[:3]=="TDL"){
							_,err:=DB.Exec("DELETE FROM timesUpStudents WHERE studentName = ?",message[13:])
							if err!=nil{
								fmt.Println(err)	
							}
						}else{
							fmt.Println("Insert!!")
							_,err:=DB.Exec("INSERT INTO timesUpStudents(name,studentID) value(?,?)",message[13:],message[3:13])
							if err!=nil{
								fmt.Println(err)	
							}
						}
					
				}
			}

		} else {
			fmt.Println("Server error!")
			speech := htgotts.Speech{Folder: "audio", Language: "en"}
			speech.Speak("No Connection")
		}
		time.Sleep(3 * time.Second)
	}
	DB.Close()

}
func Call() {
	var date string
	DB,err:=sql.Open("mysql","root:Morigo2020@/timesUpStudent?charset=utf8")
	if err!=nil{
		fmt.Println(err)
	}
	VID ,err := DB.Query("select date from todaysDate")
	if err!=nil{
		fmt.Println(err)
	}
	for VID.Next(){
			err=VID.Scan(&date)
			if err!=nil{
				fmt.Println(err)	
			}
	}
	if date != time.Now().Format("01-02"){
		_,err:=DB.Exec("DELETE FROM timesUpStudents")
		if err!=nil{
			fmt.Println(err)	
		}
		_,err=DB.Exec("UPDATE todaysDate set date = ?",time.Now().Format("01-02"))
		if err!=nil{
			fmt.Println(err)	
		}
	}
	for {
		fmt.Println("Start Call!")
		var name string
		VID ,err := DB.Query("SELECT name FROM timesUpStudents")
		if err!=nil{
			fmt.Println(err)
		}
		for VID.Next(){
				err=VID.Scan(&name)
				if err!=nil{
				fmt.Println(err)	
				}
				fmt.Println("Call:"+name)
				for j := 0; j < 2; j++ {
						speech := htgotts.Speech{Folder: "audio", Language: "en"}
						speech.Speak("hi "+name + "Go Home")
				}
		}
		time.Sleep(3 * time.Second)
	}
	DB.Close()
}
func bigthread(name string) {
	takepic(name)
	base_64(name)
	sendfile(name)
}
func Server2() {
	fmt.Println("Launching Server......")
	ln, _ := net.Listen("tcp", ":8069")

	for {
		conn, err := ln.Accept()
		if err != nil {

			fmt.Println(err)
			continue
		} else {

		}
		go runHandle(conn)
	}
}
func runHandle(conn net.Conn) {
	handleServer(conn)
}
func handleServer(conn net.Conn) {
	var studentID string
	DB,err:=sql.Open("mysql","root:Morigo2020@/timesUpStudent?charset=utf8")
	if err!=nil{
		fmt.Println(err)
	}
	message, _ := bufio.NewReader(conn).ReadString('\n')
      if(message=="000200001"){
				message="0000010001"
	}
	if(message=="000200002"){
				message="0000010002"
	}
	if(message=="000200003"){
				message="0000010003"
	}
	VID ,err := DB.Query("SELECT studentID FROM timesUpStudents")
	if err!=nil{
		fmt.Println(err)
	}
	for VID.Next(){
		err=VID.Scan(&studentID)
		if err!=nil{
			fmt.Println(err)	
		}
		if(studentID == message){
			_,err:=DB.Exec("DELETE FROM timesUpStudents WHERE studentID = ?",message)
			if err!=nil{
				fmt.Println(err)	
			}
		}
	}
	if len(message) == 10 {
		go bigthread(message)
	}
	DB.Close()
}
func main() {
	go TTS()
	go Call()
	Server2()
}


