package main

import (
	"ji_sign/util"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const (
	sckey  string = "xxxxx"
	serurl string = "https://sc.ftqq.com/" + sckey + ".send?text="
)

func init() {
	//获取执行文件路径
	util.GetExecutePath()
	//加载配置文件
	util.LoadConfig()
	util.OpenLogFile()

}

func main() {
	// linux下部署直接使用自带的crontab
	sign()
}

//登录并签到
func sign() {
	var sendMsg string = ""
	// create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("j02.space", "sc.ftqq.com"),
	)

	// authenticate
	err := c.Post("https://j02.space/signin", map[string]string{"email": util.AppConfig.GetString("email"), "passwd": util.AppConfig.GetString("passwd")})
	if err != nil {
		log.Fatal(err)
		util.Log(err.Error())
	}

	c.OnResponse(func(r *colly.Response) {
		v, _ := zhToUnicode(r.Body)
		if strings.Contains(v, "签到") {
			sendMsg = v
		}
		util.Log("response revice :" + v)

	})
	c.Visit("https://j02.space/xiaoma/get_user")
	//签到
	err = c.Post("https://j02.space/user/checkin", map[string]string{})
	if err != nil {
		log.Fatal(err)
		util.Log(err.Error())
	}
	if len(sendMsg) == 0 {
		sendMsg = "几鸡签到成功!"
	}
	c.Visit(serurl + sendMsg)
}

func zhToUnicode(raw []byte) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}
