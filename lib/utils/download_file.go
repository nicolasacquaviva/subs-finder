package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filePath string, fileUrl string, ext string) error {
	res, err := http.Get(fileUrl + ext)

	if res.StatusCode != 200 {
		return DownloadFile(filePath, fileUrl, ".zip")
	}

	if err != nil {
		return err
	}

	defer res.Body.Close()

	out, err := os.Create(filePath + ext)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)
	fmt.Println("File downloaded")

	return err
}
