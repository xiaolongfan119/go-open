package controller

import (
	"fmt"
	"io"
	"os"
	"time"

	hp "github.com/xiaolongfan119/go-open/library/net/http/hypnus"
)

type FileController struct{}

func (ctr *FileController) Upload(ctx *hp.Context) {
	r := ctx.Request

	r.ParseMultipartForm(32 << 20)
	//获取文件句柄，然后对文件进行存储等处理
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println("form file err: ", err)
		return
	}
	defer file.Close()
	//fmt.Fprintf(w, "%v", handler.Header)

	//创建上传的目的文件
	f, err := os.OpenFile(fmt.Sprintf("./../assets/%v.png", time.Now().Format("20060102150405")), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err: ", err)
		return
	}
	defer f.Close()
	//拷贝文件
	io.Copy(f, file)

	ctx.JSON(nil, nil)
}
