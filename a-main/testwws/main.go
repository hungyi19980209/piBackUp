package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	"io"
"github.com/hajimehoshi/oto"
	"strings"
	"github.com/hajimehoshi/go-mp3"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/gorilla/websocket"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var (
	//Logger 用來顯示錯誤的
	Logger *log.Logger

	//File 打開儲存錯誤的檔案
	File *os.File
)
var addr = flag.String("addr", "baohugo.com", "http service address")

func main() {
	
	time.Sleep(30 * time.Second)
	run("open")
	//time.Sleep(1 * time.Minute)
	
	File, err := os.OpenFile("/home/pi/Desktop/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	Logger = log.New(File, "", log.LstdFlags)
	Logger.SetPrefix("Test- ") // 設定日誌字首
	Logger.SetFlags(log.LstdFlags | log.Lshortfile)

for{
	server()
	time.Sleep(1*time.Second)
	}
	
}

func server(){
	
		flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/wss/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		Logger.Println("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				Logger.Println("read:", err)
				return
			}
			msg := string(message)
			
			if msg == "again"{
				err = c.WriteMessage(websocket.TextMessage, []byte("110002"))
				if err != nil {
					Logger.Println("write:", err)
					return
				}
			}else if msg == "alive" {
				fmt.Println("getin")
			}else{
					if strings.Count(msg,"|")==2{
						getName := strings.Split(msg,"|")[0]
						checkVoice(getName)
						run(getName)
						
					}
			}
		}
	}()


	err = c.WriteMessage(websocket.TextMessage, []byte("110002"))
	if err != nil {
		Logger.Println("write:", err)
		return
	}
	for {
		err := c.WriteMessage(websocket.TextMessage, []byte("baohugoHeartTest"))
			if err != nil {
				Logger.Println("write:", err)
				return
			}
			time.Sleep(1*time.Second)
	}
	
}
//-----  檢查並確認是否有語音檔
func checkVoice(a string) {
	_, err := os.Open("/home/pi/Desktop/audio/" + a + ".mp3")
	if err != nil {
		Get(a)
	}

}

//Get 取得mp3檔
func Get(a string) {
	message := "<speak>" + a + "<break time=\"0.1s\"/>準備回家<break time=\"1s\"/>" +
		a + "<break time=\"0.1s\"/>準備回家<break time=\"1s\"/>" + "</speak>"
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		Logger.Println(err)
	}

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: message},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "cmn-tw",
			Name:         "cmn-TW-Standard-A",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:  0.8,
			Pitch:         0.1,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		Logger.Println(err)
	}

	// The resp's AudioContent is binary.
	filename := "/home/pi/Desktop/audio/" + a + ".mp3"
	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	if err != nil {
		Logger.Println(err)
	}
}
func run(name string) error {
	f, err := os.Open("/home/pi/Desktop/audio/" + name + ".mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 44000)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	//fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
