package client

import (
	"encoding/xml"
	"fmt"
)

type BaseResp struct {
	// "response" OR "error"
	XMLName   xml.Name
	ErrorCode int `xml:"code"`
}

type baseRespInterface interface {
	setTagName(tagName string)
}

func parseResp(input []byte, output interface{}) error {
	if output == nil {
		return fmt.Errorf("output is nil")
	}

	r, ok := output.(baseRespInterface)
	if !ok {
		return fmt.Errorf("is not part of baseRespInterface")
	}

	r.setTagName("response")
	err := xml.Unmarshal(input, output)
	if err == nil {
		return nil
	}

	r.setTagName("error")
	return xml.Unmarshal(input, output)
}
