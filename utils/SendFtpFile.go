package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jlaffaye/ftp"
)

func SendFileFtp(fileName string, data string) error {
	ftpClient, err := ftp.Dial(os.Getenv("FTP_HOST") + ":" + os.Getenv("FTP_PORT"))
	if err != nil {
		return fmt.Errorf("error connecting to the ftp server at file : %s", fileName)
	}
	defer ftpClient.Quit()

	if err = ftpClient.Login(os.Getenv("FTP_USER"), os.Getenv("FTP_PASS")); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error getting authenticaed at the ftp server at file : %s", fileName)
	}

	dataAsBuffer := bytes.NewBufferString(data)

	if err = ftpClient.Stor(fmt.Sprintf("/MK27015/ReceiptsList/%s", fileName), dataAsBuffer); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error uploading the file server at file : %s", fileName)
	}
	return nil
}
