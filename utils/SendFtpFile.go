package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jlaffaye/ftp"
)

func SendFileFtp(fileName string, data string, brandName string) error {
	ftpClient, err := ftp.Dial(os.Getenv("FTP_HOST") + ":" + os.Getenv("FTP_PORT"))
	if err != nil {
		return fmt.Errorf("error connecting to the ftp server at file : %s", fileName)
	}
	defer ftpClient.Quit()

	var fptUser string
	var fptPass string

	if brandName == "simply" || brandName == "SIMPLY" {
		fptUser = os.Getenv("FTP_USER_SIMPLY")
		fptPass = os.Getenv("FTP_PASS_SIMPLY")
	} else if brandName == "magasin" || brandName == "MAGASIN" {
		fptUser = os.Getenv("FTP_USER_MAGASIN")
		fptPass = os.Getenv("FTP_PASS_MAGASIN")
	} else {
		return fmt.Errorf("unknown brand error getting the ftp user and pass at file : %s", fileName)
	}

	if err = ftpClient.Login(fptUser, fptPass); err != nil {
		fmt.Println("Error login into FTP Server: ", err)
		return fmt.Errorf("error getting authenticaed at the ftp server at file : %s", fileName)
	}

	dataAsBuffer := bytes.NewBufferString(data)

	if err = ftpClient.Stor(fmt.Sprintf("/%s/ReceiptsList/%s", fptUser, fileName), dataAsBuffer); err != nil {
		fmt.Println(err)
		return fmt.Errorf("error uploading the file server at file : %s", fileName)
	}
	return nil
}
