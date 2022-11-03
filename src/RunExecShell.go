package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"./mylog"
	"golang.org/x/crypto/ssh"
)

var (
	bar     = NewBarWithGraph(0, 100, "#")
	sshPort = 22
	sshUser = "root"
)

func initData() map[string]string {
	var countryCapitalMap map[string]string
	countryCapitalMap = make(map[string]string)
	fmt.Println("请任意原环境输入ip：示例：127.0.0.1，注意：如果是当前服务器请使用：127.0.0.1")
	var sourceSshHost string
	reader := bufio.NewReader(os.Stdin)
	sourceSshHost, _ = reader.ReadString('\n')
	sourceSshHost = strings.TrimSpace(sourceSshHost)
	if sourceSshHost == "" {
		sourceSshHost = "34.205.125.129"
	}
	bar.Add(1)
	fmt.Printf("输入的原主机ip：%s,进行ip正确性校验\n", sourceSshHost)
	address := net.ParseIP(sourceSshHost)
	if address == nil {
		mylog.Error.Fatal("ip地址格式不正确，请重新运行程序，程序自动停止，bye")
		bar.Add(1)
		os.Exit(0)
	} else {
		countryCapitalMap["sourceSshHost"] = sourceSshHost
		bar.Add(1)
		mylog.Info.Println("原ip地址格式正确,继续运行....")
	}
	fmt.Println("请任意迁移环境输入ip：示例：127.0.0.1")
	var targetSshHost string
	reader1 := bufio.NewReader(os.Stdin)
	targetSshHost, _ = reader1.ReadString('\n')
	targetSshHost = strings.TrimSpace(targetSshHost)
	if targetSshHost == "" {
		targetSshHost = "44.212.35.1"
	}
	fmt.Printf("输入的迁移主机ip：%s,进行ip正确性校验\n", targetSshHost)
	address = net.ParseIP(targetSshHost)
	if address == nil {
		mylog.Info.Println("ip地址格式不正确，请重新运行程序，程序自动停止，bye")
		bar.Add(1)
		os.Exit(0)
	}
	countryCapitalMap["targetSshHost"] = targetSshHost
	bar.Add(2)
	mylog.Info.Println("\n原ip地址格式正确,继续运行....")
	fmt.Println("输入的主机密码：注意原主机和迁移主机密码需要一致否则无法进行迁移")
	var passwd string
	reader2 := bufio.NewReader(os.Stdin)
	passwd, _ = reader2.ReadString('\n')
	passwd = strings.TrimSpace(passwd)
	if passwd == "" {
		passwd = "shizeying"
	}
	countryCapitalMap["passwd"] = passwd
	return countryCapitalMap

}
func getDockerMap() map[string]string {
	mylog.Info.Println("示例如下：docker run -d --security-opt seccomp:unconfined（请输入docker运行时候的参数）  --name looper（请输入需要创建的docker名称）  busybox（请输入需要使用的镜像名称如：busybox） /bin/sh -c \"i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done\"（docker命令行）")
	var dockerCapitalMap map[string]string
	dockerCapitalMap = make(map[string]string)
	fmt.Println("请输入docker运行时候的参数，如：--security-opt seccomp:unconfined")
	var dockerRunScript string
	reader := bufio.NewReader(os.Stdin)
	dockerRunScript, _ = reader.ReadString('\n')
	dockerRunScript = strings.TrimSpace(dockerRunScript)
	if dockerRunScript == "" {
		dockerRunScript = "--security-opt seccomp:unconfined"
	}
	dockerCapitalMap["dockerRunScript"] = dockerRunScript
	fmt.Println("docker命令行，如：/bin/sh -c 'i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done'")
	var script string
	reader2 := bufio.NewReader(os.Stdin)
	script, _ = reader2.ReadString('\n')
	script = strings.TrimSpace(script)
	if script == "" {
		script = "/bin/sh -c 'i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done'"
	}
	dockerCapitalMap["script"] = script
	fmt.Println("请输入需要创建的docker名称，如：looper")
	var dockerName string
	reader3 := bufio.NewReader(os.Stdin)
	dockerName, _ = reader3.ReadString('\n')
	dockerName = strings.TrimSpace(dockerName)
	if dockerName == "" {
		dockerName = "looper"
	}
	dockerCapitalMap["dockerName"] = dockerName
	fmt.Println("请输入需要使用的镜像名称如：busybox")
	var dockerImage string
	reader4 := bufio.NewReader(os.Stdin)
	dockerImage, _ = reader4.ReadString('\n')
	dockerImage = strings.TrimSpace(dockerImage)
	if dockerImage == "" {
		dockerImage = "busybox"
	}
	dockerCapitalMap["dockerImage"] = dockerImage

	return dockerCapitalMap
}
func getSession(client *ssh.Client) *ssh.Session {
	session, err := client.NewSession()

	if err != nil {
		bar.load()
		mylog.Info.Fatal("创建ssh session 失败", err)
	}
	return session
}

func closeSession(session *ssh.Session) {

	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)
}
func main() {

	getStep("第一步：用户信息初始化")
	countryCapitalMap := initData()
	getStep("第二步：获取docker运行参数")
	dockerCapitalMap := getDockerMap()
	getStep("第三步：原服务器开始进行初始化")
	client := toObtainSshClient(countryCapitalMap["sourceSshHost"], sshPort, countryCapitalMap["passwd"], sshUser)

	bar.Add(4)
	session := getSession(client)
	err2 := session.Run(fmt.Sprintf("docker ps -a| grep %s  | gawk '{cmd=\"docker stop \"$1; system(cmd)}' && docker ps -a| grep %s  | gawk '{cmd=\"docker rm \"$1; system(cmd)}'", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerName"]))
	mylog.Debug.Println(fmt.Sprintf("docker ps -a| grep %s  | gawk '{cmd=\"docker stop \"$1; system("+
		"cmd)}' && docker ps | grep %s  -a| gawk '{cmd=\"docker rm \"$1; system(cmd)}'", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerName"]))
	if err2 != nil {
		mylog.Error.Fatal("docker任务执行", err2)
		bar.load()
		os.Exit(0)
	}
	closeSession(session)
	session = getSession(client)
	dockerRun := fmt.Sprintf("docker run -d --name %s %s %s  %s ", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerRunScript"], dockerCapitalMap["dockerImage"], dockerCapitalMap["script"])
	mylog.Debug.Println(fmt.Sprintf("执行命令：%s", dockerRun))
	_ = fmt.Sprintf("docker exec -i %s cat /proc/self/cgroup | head -1 |awk -F '/' '{print $NF}'", dockerCapitalMap["dockerName"])
	closeSession(session)
	time.Sleep(5 * time.Second)

	// 1.创建docker进程
	session = getSession(client)
	combo, err := session.CombinedOutput(dockerRun)
	if err != nil {
		mylog.Error.Fatal("docker任务执行,创建镜像失败", err)
		bar.load()
		os.Exit(0)
	}
	mylog.Debug.Println(fmt.Sprintf("命令输出:%s", string(combo)))
	// 2.获取containerId
	containerId := string(combo)

	closeSession(session)
	mylog.Info.Printf("等待日志，静候30s....")
	time.Sleep(20 * time.Second)
	mylog.Info.Printf("继续执行....")
	session = getSession(client)
	// 3.创建checkpoint
	checkPoint := fmt.Sprintf("docker checkpoint create %s c1", dockerCapitalMap["dockerName"])
	mylog.Debug.Println(checkPoint)
	_, err1 := session.CombinedOutput(checkPoint)

	if err1 != nil {
		mylog.Error.Fatal("docker任务执行,checkpoint执行失败", err1)
		bar.load()
		os.Exit(0)
	}
	closeSession(session)
	session = getSession(client)
	// 4. 转换为镜像
	commit := fmt.Sprintf("docker commit %s checkpoint", strings.TrimSpace(containerId))
	mylog.Debug.Println(commit)
	err = session.Run(commit)
	if err != nil {
		mylog.Error.Fatal("docker任务执行,转换为镜像失败", err)
		bar.load()
		os.Exit(0)
	}
	closeSession(session)
	session = getSession(client)
	// 5.导出镜像
	save := "docker save -o /opt/checkpoint checkpoint"
	mylog.Debug.Println(save)
	err = session.Run(save)
	if err != nil {
		mylog.Error.Fatal("docker任务执行,导出为镜像包失败", err)
		bar.load()
		os.Exit(0)
	}
	closeSession(session)
	session = getSession(client)
	// 6.获取当前操作镜像的最后一行日志
	scanLog := fmt.Sprintf("docker logs -f   %s ", dockerCapitalMap["dockerName"])
	mylog.Info.Println(scanLog)
	combo, err = session.CombinedOutput(scanLog)
	if err != nil {
		mylog.Error.Fatal("docker任务执行,查看日志失败", err)
		os.Exit(0)
	}
	tailLogBySource := string(combo)

	closeSession(session)
	session = getSession(client)
	// 7.发送镜像包
	scp := fmt.Sprintf("scp /opt/checkpoint root@%s:/opt/checkpoint", countryCapitalMap["targetSshHost"])
	mylog.Debug.Println(scp)
	combo, err = session.CombinedOutput(scp)
	if err != nil {
		mylog.Error.Fatal("docker任务执行,发送镜像包失败", err)
		os.Exit(0)
	}
	closeSession(session)

	// 8.进入target服务器
	clientTarget := toObtainSshClient(countryCapitalMap["targetSshHost"], sshPort, countryCapitalMap["passwd"], sshUser)
	session = getSession(clientTarget)
	err2 = session.Run(fmt.Sprintf("docker ps -a| grep %s  | gawk '{cmd=\"docker stop \"$1; system(cmd)}' && docker ps -a| grep %s  | gawk '{cmd=\"docker rm \"$1; system(cmd)}'", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerName"]))
	if err2 != nil {
		mylog.Error.Fatal("迁移主机上docker任务执行", err2)
		bar.load()
		os.Exit(0)
	}
	closeSession(session)
	session = getSession(clientTarget)
	// 9.执行target load命令
	load := "docker load -i /opt/checkpoint"
	mylog.Debug.Println(load)
	err = session.Run(load)
	if err != nil {
		mylog.Error.Fatal("迁移主机上docker任务,导入镜像包失败", err)
		os.Exit(0)
	}
	closeSession(session)
	session = getSession(clientTarget)
	// 10.创建任务
	dockerRunTarget := fmt.Sprintf("docker run -d --name %s %s %s  %s && docker stop %s", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerRunScript"], "checkpoint", dockerCapitalMap["script"], dockerCapitalMap["dockerName"])
	mylog.Debug.Println(dockerRunTarget)
	combo, err = session.CombinedOutput(dockerRunTarget)
	if err != nil {
		mylog.Error.Fatal("迁移主机上docker任务执行,创建镜像包失败", err)
		os.Exit(0)
	}
	time.Sleep(5 * time.Second)
	containerIdTarget := strings.Split(string(combo), "\n")[0]
	closeSession(session)
	session = getSession(client)
	// 11.拷贝checkpoint到目的节点
	scpCheckPoint := fmt.Sprintf(
		"scp -r /var/lib/docker/containers/%s/checkpoints/c1/ root@%s:/var/lib/docker/containers/%s/checkpoints/",
		strings.TrimSpace(containerId), countryCapitalMap["targetSshHost"], strings.TrimSpace(containerIdTarget))
	mylog.Debug.Println(scpCheckPoint)
	err = session.Run(scpCheckPoint)
	if err != nil {
		mylog.Error.Fatal("迁移主机上docker任务执行,拷贝checkpoint到目的节点失败", err)
		os.Exit(0)
	}
	time.Sleep(5 * time.Second)
	closeSession(session)
	session = getSession(client)
	// 12.启动容器
	scpCheckPointTarget := fmt.Sprintf(
		"scp -r /var/lib/docker/containers/%s/checkpoints/c1/ root@%s:/var/lib/docker/containers/%s/checkpoints/",
		strings.TrimSpace(containerId), countryCapitalMap["targetSshHost"], strings.TrimSpace(containerIdTarget))
	mylog.Debug.Println(scpCheckPointTarget)
	err = session.Run(scpCheckPointTarget)
	if err != nil {
		mylog.Error.Fatal("迁移主机上docker任务执行,拷贝checkpoint到目的节点失败", err)
		os.Exit(0)
	}
	// 13. 恢复位点差
	checkPointTarget := fmt.Sprintf("docker start --checkpoint c1 %s", strings.TrimSpace(containerIdTarget))
	mylog.Debug.Println(checkPointTarget)
	closeSession(session)
	session = getSession(clientTarget)
	combo, err = session.CombinedOutput(checkPointTarget)
	if err != nil {
		mylog.Error.Fatal("迁移主机上docker任务执行,启动checkPoint失败", err)
		os.Exit(0)
	}
	mylog.Info.Printf("命令输出:%s\n", string(combo))
	closeSession(session)
	session = getSession(clientTarget)
	// 13.再次检查进程日志正常，接着上次创建checkpoint的时间点打印
	scanLogTarget := fmt.Sprintf(
		"docker logs --tail  all %s ", strings.TrimSpace(containerIdTarget))
	time.Sleep(5 * time.Second)
	combo, err = session.CombinedOutput(scanLogTarget)
	if err != nil {
		mylog.Error.Fatal("迁移主机上docker任务执行,日志获取失败", err)
		os.Exit(0)
	}
	mylog.Info.Printf("打印原主机日志：\n%s\n", tailLogBySource)

	mylog.Info.Printf("打印迁移主机日志：\n%s\n", string(combo))

	defer func(client *ssh.Client) {
		var err = client.Close()
		if err != nil {

		}
	}(client)
	defer func(client *ssh.Client) {
		var err = client.Close()
		if err != nil {

		}
	}(clientTarget)
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
		bar.load()
		mylog.Error.Fatal("创建ssh client 失败", err)
	}

	return sshClient
}

func getStep(message string) string {
	return fmt.Sprintf("=================================%s====================================================", message)

}

type Bar struct {
	mu      sync.Mutex
	graph   string    // 显示符号
	rate    string    // 进度条
	percent int       // 百分比
	current int       // 当前进度位置
	total   int       // 总进度
	start   time.Time // 开始时间
}

func NewBar(current, total int) *Bar {
	bar := new(Bar)
	bar.current = current
	bar.total = total
	bar.start = time.Now()
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < bar.percent; i += 2 {
		bar.rate += bar.graph // 初始化进度条位置
	}
	return bar
}

func NewBarWithGraph(start, total int, graph string) *Bar {
	bar := NewBar(start, total)
	bar.graph = graph
	return bar
}

func (bar *Bar) getPercent() int {
	return int((float64(bar.current) / float64(bar.total)) * 100)
}

func (bar *Bar) getTime() (s string) {
	u := time.Now().Sub(bar.start).Seconds()
	h := int(u) / 3600
	m := int(u) % 3600 / 60
	if h > 0 {
		s += strconv.Itoa(h) + "h "
	}
	if h > 0 || m > 0 {
		s += strconv.Itoa(m) + "m "
	}
	s += strconv.Itoa(int(u)%60) + "s"
	return
}

func (bar *Bar) load() {
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\n\r[%-50s]% 3d%%    %2s   %d/%d\n", bar.rate, bar.percent, bar.getTime(), bar.current, bar.total)
}
func (bar *Bar) Reset(current int) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.current = current
	bar.load()

}

func (bar *Bar) Add(i int) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.current += i
	bar.load()
}
