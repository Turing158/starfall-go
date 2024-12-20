package service

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"starfall-go/util"
	"strconv"
	"strings"
)

type OtherService struct {
}

var redisUtil = util.RedisUtil{}

func (OtherService) GetCodeImage(c *gin.Context) {
	oldId := c.GetHeader("Captcha-Id")
	if oldId != "" {
		key := "captcha:" + oldId
		if redisUtil.Has(key) {
			redisUtil.Del(key)
		}
	}
	id, base64s, _, err := util.CreateAndSaveCaptcha()
	if err != nil {
		return
	}
	data, err := base64.StdEncoding.DecodeString(strings.Split(base64s, ",")[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Header("Content-length", strconv.Itoa(len(data)))
	c.Header("Base64Img", base64s)
	c.Header("Captcha-Id", id)
	c.Data(200, "image/png", data)
}
