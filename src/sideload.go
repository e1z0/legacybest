package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	sshclient "github.com/helloyi/go-sshclient"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type sftpClient struct {
	host, user, password string
	port                 int
	*sftp.Client
}

// Create a new SFTP connection by given parameters
func NewConn(host, user, password string, port int) (client *sftpClient, err error) {
	switch {
	case `` == strings.TrimSpace(host),
		`` == strings.TrimSpace(user),
		`` == strings.TrimSpace(password),
		0 >= port || port > 65535:
		return nil, errors.New("Invalid parameters")
	}

	client = &sftpClient{
		host:     host,
		user:     user,
		password: password,
		port:     port,
	}

	if err = client.connect(); nil != err {
		return nil, err
	}
	return client, nil
}

func (sc *sftpClient) connect() (err error) {
	config := &ssh.ClientConfig{
		User:            sc.user,
		Auth:            []ssh.AuthMethod{ssh.Password(sc.password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	addr := fmt.Sprintf("%s:%d", sc.host, sc.port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	// create sftp client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	sc.Client = client

	return nil
}

// Upload file to sftp server
func (sc *sftpClient) Put(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		sc.Mkdir(path)
	}

	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return
}

// Download file from sftp server
func (sc *sftpClient) Get(remoteFile, localFile string) (err error) {
	srcFile, err := sc.Open(remoteFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return
}

type Sshopts struct {
	User     string
	Password string
	Host     string
	Port     int
	Privkey  bool
	KeyFile  string
}

func TransferToRaspberry(opts Sshopts, filename string) (err error) {
	conn, err := NewConn(opts.Host, opts.User, opts.Password, opts.Port)
	if err == nil {
		fmt.Printf("Connected!\n")
		err := conn.Put(filename, filepath.Base(filename))
		if err != nil {
			fmt.Printf("Unable to put file to remote host! err: %s\n", err)
			return nil
		}
		fmt.Printf("File: %s was successfully copied!\n", filename)
	} else {
		fmt.Printf("Unable to connect! %s\n", err)
		return err
	}
	client, err := sshclient.DialWithPasswd(fmt.Sprintf("%s:%d", opts.Host, opts.Port), opts.User, opts.Password)
	if opts.Privkey {
		client, err = sshclient.DialWithKey(opts.Host+":"+string(opts.Port), opts.User, opts.KeyFile)
	}
	// Dial with private key and a passphrase to decrypt the key
	//client, err := DialWithKeyWithPassphrase("host:port", "username", "prikeyFile", "my-passphrase"))
	if err != nil {
		fmt.Printf("Unable to connect to remote ssh server: %s for command injection! err: %s\n", opts.Host, err)
		return err
	}
	defer client.Close()
	fmt.Printf("Waiting for user input on another side :)\n")
	out, err := client.Cmd(fmt.Sprintf("/usr/bin/kermit -s %s", filepath.Base(filename))).SmartOutput()
	if err != nil {
		fmt.Printf("Unable to execute command in remote ssh server: %s stdout: %s err: %s\n", opts.Host, out, err)
		return err
	}
	fmt.Println(string(out))
	return nil
}
