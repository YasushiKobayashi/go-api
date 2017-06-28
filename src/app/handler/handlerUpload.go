package handler

import (
	"app/config"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

func Upload(file *multipart.FileHeader) (res string, err error) {
	src, err := file.Open()
	if err != nil {
		return res, err
	}
	defer src.Close()

	pos := strings.LastIndex(file.Filename, ".")
	var fileName = strconv.FormatInt(time.Now().Unix(), 10) + file.Filename[pos:]

	// Destination
	dst, err := os.Create(config.UPLOAD_DIR + fileName)
	if err != nil {
		return res, err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return res, err
	}
	res = "//" + config.URL + config.UPLOAD_PATH + fileName
	return res, err
}
