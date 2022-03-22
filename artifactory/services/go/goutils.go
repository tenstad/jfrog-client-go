package _go

import (
	"strings"

	"github.com/tenstad/jfrog-client-go/artifactory/services/utils"
)

func CreateUrlPath(moduleId, version, props, extension string, url *string) error {
	*url = strings.Join([]string{*url, moduleId, "@v", version + extension}, "/")
	properties, err := utils.ParseProperties(props)
	if err != nil {
		return err
	}

	*url = strings.Join([]string{*url, properties.ToEncodedString(true)}, ";")
	if strings.HasSuffix(*url, ";") {
		tempUrl := *url
		tempUrl = tempUrl[:len(tempUrl)-1]
		*url = tempUrl
	}
	return nil
}
