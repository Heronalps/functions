package main

import (
	"fmt"
	"os"
	"io/ioutil"

	"github.com/dutchcoders/goftp"
	"github.com/urfave/cli"
)

var VolumeContainerPath = "/home/taejoon/kjwook/iron/data/"

func fileUploadFn() cli.Command {

	return cli.Command {
		Name:		"fileupload",
		Usage:		"file upload to server",
		Description:	"file upload to server through ftp",
		ArgsUsage:	"<FTP SERVER ADDR><ID><PASSWORD><FilePath>",
		Action:		fileUpload,
	}
}

func fileUpload(c *cli.Context) error {
	addr := c.Args().Get(0)
	id := c.Args().Get(1)
	pw := c.Args().Get(2)
	filePath := c.Args().Get(3)

	fmt.Println(addr +", " + id + ", " + pw + ", " + filePath)

	var err error
	var ftp *goftp.FTP

	// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
	if ftp, err = goftp.Connect(addr); err != nil {
		panic(err)
	}

	defer ftp.Close()
	fmt.Println("Successfully connected to", ftp)

	// Username / password authentication
	if err = ftp.Login(id, pw); err != nil {
		panic(err)
	}

	if err = ftp.Cwd("/"); err != nil {
		panic(err)
	}

	file, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	switch mode := file.Mode(); {
	case mode.IsDir():
		fmt.Println("Directory")
		if err := ftp.Mkd(VolumeContainerPath + filePath); err != nil {
			panic(err)
		}
		files, _ := ioutil.ReadDir(filePath)
		fmt.Println(len(files))
		for _, f := range files {
			fmt.Println(f.Name())
			sendFileToServer(ftp, filePath + "/" + f.Name())

		}
	case mode.IsRegular():
		fmt.Println("File")
	}

	return nil
}


func sendFileToServer(ftp *goftp.FTP, fileName string) {
	var err error
	// Upload a file
	var file *os.File

	if file, err = os.Open(fileName); err != nil {
		panic(err)
	}
	if err := ftp.Stor(VolumeContainerPath + fileName, file); err != nil {
		panic(err)
	}
	fmt.Println(fileName)

}
