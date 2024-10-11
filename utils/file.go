package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"es-3d-editor-go-back/controllers/system"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DelFile 删除文件
// @param path string 文件基础路径
func DelFile(path string) (err error) {
	if _, err := os.Stat(path); err == nil {
		//删除文件
		if delErr := os.Remove(path); delErr != nil {
			//如果删除失败则输出 file remove Error!
			return errors.New("文件 " + path + " 删除失败！")
		} else {
			//如果删除成功则输出 file remove OK!
			return nil
		}
	} else {
		return errors.New("文件 " + path + " 不存在！")
	}
}

// DownloadHttpImg 下载网络图片保存至又拍云
// @param savePath string 保存位置
// @param url string 网络图片地址
func DownloadHttpImg(savePath string, url string) (filePath string, err error) {
	fmt.Printf("DownloadHttpImg 下载网络图片.参数：保存位置：%s,图片地址：%s \n", savePath, url)
	pic := savePath + "/" + time.Now().Format("20060102")
	idx := strings.LastIndex(url, "/")
	if idx < 0 {
		pic += "/" + url
	} else {
		pic += url[idx:]
	}

	//文件名加上时间戳
	fmt.Printf("DownloadHttpImg 文件名加上时间戳：起始位置位置：%s,时间戳：%s，后缀：%s \n", pic[:strings.LastIndex(pic, ".")], strconv.FormatInt(time.Now().Unix(), 10), pic[strings.LastIndex(pic, "."):])
	pic = pic[:strings.LastIndex(pic, ".")] + "-" + strconv.FormatInt(time.Now().Unix(), 10) + pic[strings.LastIndex(pic, "."):]

	res, err := http.Get(url)
	if err != nil {
		return "", errors.New("下载网络图片失败！error:" + err.Error())
	}

	err = system.UpYunUpload(pic, res.Body)

	if err != nil {
		return "", errors.New("文件存至又拍云失败!error:" + err.Error())
	}

	return pic, nil
}

// WriteBase64Img 将Base64 图片保存为普通图片格式
// @param savePath string 保存位置
// @param base64Image string base64图片内容
func WriteBase64Img(savePath string, base64Image string) (filePath string, err error) {
	_, err = regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64Image)
	if err != nil {
		return "", errors.New("base64 信息编码失败！error:" + err.Error())
	}

	re, _ := regexp.Compile(`^data:\s*image\/(\w+);base64,`)
	allData := re.FindAllSubmatch([]byte(base64Image), 2)
	fileType := string(allData[0][1]) //png ，jpeg 后缀获取

	base64Str := re.ReplaceAllString(base64Image, "")

	pic := savePath + "/" + time.Now().Format("20060102")

	curFileStr := strconv.FormatInt(time.Now().UnixNano(), 10)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(99999)

	var file = pic + "/base64-" + curFileStr + strconv.Itoa(n) + "." + fileType
	by, _ := base64.StdEncoding.DecodeString(base64Str)

	// []byte 转换为 io.ReaderCloser
	reader := ioutil.NopCloser(bytes.NewReader(by))

	err = system.UpYunUpload(file, reader)
	if err != nil {
		return "", errors.New("base64 转换为文件保存至又拍云失败!error:" + err.Error())
	}

	return file, nil
}
