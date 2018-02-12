package translate

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Translate struct {
	Code       int    `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Source     string `json:"source,omitempty"`
	Target     string `json:"target,omitempty"`
	TargetText string `json:"targetText,omitempty"`
	CodeDesc   string `json:"codeDesc,omitempty"`
}

type Conf struct {
	SecretID  string
	SecretKey string
	Region    string
}

func New(secretID, secretKey, region string) *Conf {
	conf := Conf{}
	conf.SecretID = secretID
	conf.SecretKey = secretKey
	conf.Region = region

	return &conf
}

func (conf *Conf) Do(text string) (Translate, error) {
	translate := Translate{}
	params := map[string]interface{}{
		"Action":     "TextTranslate",
		"Nonce":      randomString(4),
		"Region":     "gz",
		"SecretId":   conf.SecretID,
		"Timestamp":  time.Now().Unix(),
		"sourceText": text,
		"source":     "zh",
		"target":     "en",
	}
	url := fmt.Sprintf("https://tmt.api.qcloud.com/v2/index.php?%s&Signature=%s", getParamStr(params, true), getSignature(params, conf.SecretKey))

	resp, err := http.Get(url)
	if err != nil {
		translate.Code = resp.StatusCode
		translate.Message = err.Error()
		return translate, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		translate.Code = resp.StatusCode
		translate.Message = err.Error()
		return translate, err
	}

	json.Unmarshal(body, &translate)

	if translate.Code != 0 {
		return translate, fmt.Errorf("Error:%s", translate.CodeDesc)
	}

	return translate, nil
}

//生成指定长度随机字符串
func randomString(n int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func interfaceToString(i interface{}) string {
	if str, ok := i.(string); ok {
		return str
	} else if intval, ok := i.(int); ok {
		return strconv.Itoa(intval)
	} else if intv, ok := i.(int64); ok {
		return strconv.FormatInt(intv, 10)
	} else if intv, ok := i.(uint64); ok {
		return strconv.FormatUint(intv, 10)
	}

	return ""
}

func getParamStr(params map[string]interface{}, needEscape bool) string {
	if params == nil || len(params) == 0 {
		return ""
	}

	tmp := []string{}
	for k, v := range params {
		var value string
		if needEscape == true {
			value = url.QueryEscape(interfaceToString(v))
		} else {
			value = interfaceToString(v)
		}
		tmp = append(tmp, fmt.Sprintf("%s=%s", k, value))
	}
	sort.Strings(tmp)

	return strings.Join(tmp, "&")
}

func getSignature(params map[string]interface{}, secretKey string) string {
	paramsUrl := getParamStr(params, false)

	srcStr := fmt.Sprintf("GETtmt.api.qcloud.com/v2/index.php?%s", paramsUrl)

	//HMAC_SHA1加密
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(srcStr))

	return url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

}
