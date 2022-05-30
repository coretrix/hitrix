package helper

import (
	"strings"

	"github.com/coretrix/hitrix/service"
)

const (
	Options = "fit=${fit},format=${format},metadata=none,onerror=redirect,quality=${quality},width=${width},dpr=${dpr}/"
)

func GetImageURLTemplate(image string) string {
	return service.DI().Config().MustString("oss.cdn_url") + Options + image
}

func GetImageURLTemplateFilled(image, fit, format, quality, width, dpr string) string {
	image = service.DI().Config().MustString("oss.cdn_url") + Options + image
	image = strings.Replace(image, "${fit}", fit, -1)
	image = strings.Replace(image, "${format}", format, -1)
	image = strings.Replace(image, "${quality}", quality, -1)
	image = strings.Replace(image, "${width}", width, -1)
	image = strings.Replace(image, "${dpr}", dpr, -1)
	return image
}
