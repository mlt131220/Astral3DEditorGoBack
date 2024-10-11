package utils

import (
	"errors"
	"regexp"
	"strings"
)

// IsUrl 判断是否是url
func IsUrl(str string) bool {
	if strings.HasPrefix(str, "#") || strings.HasSuffix(str, ".exe") || strings.HasSuffix(str, ":void(0);") {
		return false
	} else if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		return false
	} else if strings.EqualFold(str, "javascript:;") {
		return false
	} else {
		return true
	}
}

// RepImages html筛选图片
func RepImages(html string) ([]string, error) {
	var imgRE = regexp.MustCompile(`<img src=".*static/upload/(.*?)"`)
	if imgRE == nil {
		return nil, errors.New("MustCompile err")
	}
	// 提取关键字
	images := imgRE.FindAllStringSubmatch(html, -1)
	out := make([]string, len(images))
	for i := range out {
		out[i] = images[i][1]
	}
	return out, nil
}
