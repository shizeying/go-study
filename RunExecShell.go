package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {

	fmt.Println("请任意输入ip：")
	var sshHost string
	fmt.Scan(&sshHost)
	fmt.Printf("输入的ip：%s,进行ip正确性校验\n", sshHost)
	address := net.ParseIP(sshHost)
	if address == nil {
		log.Fatal("ip地址格式不正确，请重新运行程序，程序自动停止，bye")
		os.Exit(0)
	} else {
		log.Println("ip地址格式正确,继续运行....")
	}
	log.Println("本程序默认使用的用户： root;默认使用的端口号: 22.....")
	sshUser := "root"
	sshPassword := "shizeying"
	sshPort := 22
	log.Println(getStep("第一步：登陆服务器"))
	client := toObtainSshClient(sshHost, sshPort, sshPassword, sshUser)
	// 创建ssh-session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
	// 执行远程命令
	combo, err := session.CombinedOutput("docker exec -i looper cat /proc/self/cgroup | head -1 |awk -F '/' '{print $NF}'")
	if err != nil {
		log.Fatal("远程执行cmd 失败", err)
	}
	log.Println("命令输出:", string(combo))

	defer func(client *ssh.Client) {
		var err = client.Close()
		if err != nil {

		}
	}(client)

}

func toObtainSshClient(host string, sshPort int, passwd string, user string) *ssh.Client {
	// 创建ssh登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, // ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 这个可以， 但是不够安全
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(passwd)}
	// dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}

	return sshClient
}

func getStep(message string) string {
	return fmt.Sprintf("=================================%s====================================================", message)

}
