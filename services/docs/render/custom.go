package render

import (
	"bytes"
	"fmt"
	"io"

	"github.com/piquel-fr/api/models"
)

func renderCustom(md []byte, config *models.DocsInstance) ([]byte, error) {
	var err error
	md, err = renderSingleline(md, config)
	if err != nil {
		return nil, err
	}
	md, err = renderMultiline(md, config)
	if err != nil {
		return nil, err
	}

	return md, nil
}

func renderSingleline(md []byte, config *models.DocsInstance) ([]byte, error) {
	match := singleline.FindSubmatch(md)
	if match == nil {
		return md, nil
	}

	total, tag, param := match[0], match[1], match[2]

	var newMarkdown bytes.Buffer
	switch string(tag) {
	case "include":
		include, err := loadInclude(string(param), config)
		if err != nil {
			return nil, err
		}
		newMarkdown.Write(include)
	default:
		io.WriteString(&newMarkdown, fmt.Sprintf("Tag %s does not exist\n", tag))
	}

	md = bytes.Replace(md, total, newMarkdown.Bytes(), 1)
	return renderSingleline(md, config)
}

func renderMultiline(md []byte, config *models.DocsInstance) ([]byte, error) {
	match := multiline.FindSubmatch(md)
	if match == nil {
		return md, nil
	}

	total, tag, body := match[0], match[1], match[2]

	var newMarkdown bytes.Buffer
	switch string(tag) {
	case "warning":
		io.WriteString(&newMarkdown, "Warning:\n")
		newMarkdown.Write(body)
	default:
		io.WriteString(&newMarkdown, fmt.Sprintf("Tag %s does not exist\n", tag))
		newMarkdown.Write(body)
	}

	md = bytes.Replace(md, total, newMarkdown.Bytes(), 1)
	return renderMultiline(md, config)
}
