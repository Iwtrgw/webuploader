package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	//允许跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsAddress string

func init() {
	log.Println("系统初始化.......")
	viper.SetConfigName("core")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(fmt.Errorf("配置文件读取失败:%s", err))
		os.Exit(1)
	}
	wsAddress = viper.GetString("websocket.proto") + viper.GetString("websocket.host") + viper.GetString("websocket.port")
	log.Println("初始化完成，开始监听.....")
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	log.Println(r.RemoteAddr)
	t.Execute(w, wsAddress)
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)                     // 32M
		name := r.MultipartForm.Value["video-filename"][0] // 文件名
		video := r.MultipartForm.File["video-blob"][0]     // 图片文件
		//fmt.Println(video)
		file, err := video.Open()
		if err == nil {
			data, err := ioutil.ReadAll(file) // 读取二进制文件字节流
			if err == nil {
				// fmt.Fprintln(w, string(data))   // 将读取的字节信息输出
				// 将文件存储到项目根目录下的 images 子目录
				// 从上传文件中读取文件名并获取文件后缀
				//names := strings.Split(image.Filename, ".")
				fmt.Println(11, name)
				//suffix := names[len(names) - 1]

				// 将上传文件名字段值和源文件后缀拼接出新的文件名
				//filename := name + "." + suffix
				//filename := name
				// 创建这个文件
				newFile, _ := os.Create("videos/" + name)
				defer newFile.Close()
				// 将上传文件的二进制字节信息写入新建的文件
				_, err := newFile.Write(data)
				if err == nil {
					fmt.Fprintf(w, "图片上传成功，图片大小: %d 字节\n", 1000)
				}
			} else {
				fmt.Fprintln(w, err)
			}
		}
		if err != nil {
			fmt.Fprintln(w, err)
		}
	} else {
		fmt.Println("请求方法不对，请使用POST请求", r.Method)

	}
}

func wsCamera(w http.ResponseWriter, r *http.Request) {
	// 升级成websocket协议
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级Websocket失败")
		return
	}
	log.Println(r.Proto)
	fileName := ws.RemoteAddr().String() + "_" + time.Now().Format("2006-01-02 15:04:05") + ".mp4"
	filePath := "videos/"
	defer ws.Close()
	newFile, _ := os.Create(filePath + fileName)
	newFile.Close()
	f, err := os.OpenFile(filePath+fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer f.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("读取消息错误：", err)
			return
		}
		fmt.Println("消息正常读取")
		var start int64 = 0
		readFile(message, &start, f)
		// 消息写入
		//data := []byte{''}
		err = ws.WriteMessage(websocket.TextMessage, []byte("success"))
		if err != nil {
			fmt.Println("写入消息失败：", err)
			return
		}
	}
}

func wsDisplay(w http.ResponseWriter, r *http.Request) {
	// 升级成websocket协议
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级Websocket失败")
		return
	}
	log.Println(r.Proto)
	fileName := ws.RemoteAddr().String() + "_display_" + time.Now().Format("2006-01-02 15:04:05") + ".mp4"
	filePath := "videos/"
	f, _ := os.Create(filePath + fileName)
	//newFile.Close()
	// 文件持续写入
	//f,err :=os.OpenFile(filePath+fileName,os.O_WRONLY,0644)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer f.Close()
	defer func() {
		ws.Close()
		fmt.Println("WebSocket 关闭")
	}()
	//var start int64 =0
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("读取消息错误：", err)
			return
		}
		fmt.Println("消息正常读取")
		//readFile(message,&start,f)
		_, err = f.Write(message)
		if err != nil {
			fmt.Println("写入文件失败", err)
			return
		}
		// 消息写入
		//data := []byte{''}
		err = ws.WriteMessage(websocket.BinaryMessage, []byte("success"))
		if err != nil {
			fmt.Println("写入消息失败：", err)
			return
		}
	}
}

func readFile(data []byte, start *int64, f *os.File) {
	/*filePath := "./videos/test.webm"
	file,err :=os.OpenFile(filePath,os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for i := 0; i < 5; i++ {
		write.Write(data)
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()*/
	//fmt.Println(string(data))

	/*newFile, _ := os.Create("videos/" + "test.mp4")
	defer newFile.Close()
	// 将上传文件的二进制字节信息写入新建的文件
	_, err := newFile.Write(data)
	if err == nil {
		fmt.Printf("图片上传成功，图片大小: %d 字节\n\n", 1000)
	}else {
		fmt.Println(err)
	}*/
	// 文件持续写入
	/*f,err :=os.OpenFile(filePath,os.O_WRONLY,0644)
	if err!=nil {
		fmt.Println("文件打开失败")
		return
	}
	defer f.Close()*/
	// 查找文件末尾
	/*n,err := f.Seek(*start,2)
	if err != nil {
		fmt.Println("文件末尾催偏移失败")
		return
	}
	*start = *start + n
	fmt.Printf("文件开始偏移量:%d 文件结束偏移量:%d\n", *start,n)*/
	// 从末尾的偏移量开始写入内容
	//_, err = f.WriteAt(data, n)
	n, _ := f.Write(data)
	*start = *start + int64(n)

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadVideo)
	http.HandleFunc("/camera", wsCamera)
	http.HandleFunc("/display", wsDisplay)
	host := "0.0.0.0"
	port := viper.GetString("host.port")
	http.ListenAndServe(host+port, nil)
}
