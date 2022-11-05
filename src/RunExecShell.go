package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var (
	bar                = progressbar.Default(100)
	sshPort            = 22
	sshUser            = "root"
	enableCompression  = true
	enableCreateDocker = true
	compression        = "3"
	containerId        = ""
)

func initData() map[string]string {

	log.SetFormatter(&log.JSONFormatter{})
	var countryCapitalMap map[string]string
	countryCapitalMap = make(map[string]string)
	fmt.Println("请任意原环境输入ip：示例：127.0.0.1，注意：如果是当前服务器请使用：127.0.0.1")
	var sourceSshHost string
	sourceSshHost = getCommandStr()
	if sourceSshHost == "" {
		sourceSshHost = "10.203.56.7"
	}
	fmt.Printf("输入的原主机ip：%s,进行ip正确性校验\n", sourceSshHost)
	address := net.ParseIP(sourceSshHost)
	if address == nil {
		log.Fatal("ip地址格式不正确，请重新运行程序，程序自动停止，bye")
		os.Exit(0)
	} else {
		countryCapitalMap["sourceSshHost"] = sourceSshHost
		log.Println("原ip地址格式正确,继续运行....")
	}
	fmt.Println("请任意迁移环境输入ip：示例：127.0.0.1")
	var targetSshHost string
	targetSshHost = getCommandStr()
	if targetSshHost == "" {
		targetSshHost = "10.203.56.8"
	}
	fmt.Printf("输入的迁移主机ip：%s,进行ip正确性校验\n", targetSshHost)
	address = net.ParseIP(targetSshHost)
	if address == nil {
		log.Println("ip地址格式不正确，请重新运行程序，程序自动停止，bye")
		os.Exit(0)
	}
	countryCapitalMap["targetSshHost"] = targetSshHost
	log.Println("原ip地址格式正确,继续运行....")
	fmt.Println("输入的主机密码：注意原主机和迁移主机密码需要一致否则无法进行迁移")
	var passwd string
	passwd = getCommandStr()
	if passwd == "" {
		passwd = "1"
	}
	countryCapitalMap["passwd"] = passwd
	return countryCapitalMap

}
func getCommandStr() string {
	var script string
	// reader2 := bufio.NewReader(os.Stdin)
	// script, _ = reader2.ReadString('\n')
	script = ""
	return strings.TrimSpace(script)
}
func getDockerMap() map[string]string {
	log.Println("示例如下：docker run -d --security-opt seccomp:unconfined（请输入docker运行时候的参数）  --name looper（请输入需要创建的docker名称）  busybox（请输入需要使用的镜像名称如：busybox） /bin/sh -c \"i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done\"（docker命令行）")
	var dockerCapitalMap map[string]string
	dockerCapitalMap = make(map[string]string)
	fmt.Println("请输入docker运行时候的参数，如：--security-opt seccomp:unconfined")
	var dockerRunScript string
	dockerRunScript = getCommandStr()
	if dockerRunScript == "" {
		dockerRunScript = "--security-opt seccomp:unconfined"
	}
	dockerCapitalMap["dockerRunScript"] = dockerRunScript
	fmt.Println("docker命令行，如：/bin/sh -c 'i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done'")
	var script string
	script = getCommandStr()
	if script == "" {
		script = "/bin/sh -c 'i=0; while true; do echo $i; i=$(expr $i + 1); sleep 1; done'"
	}
	dockerCapitalMap["script"] = script
	fmt.Println("请输入需要创建的docker名称，如：looper")
	var dockerName string
	dockerName = getCommandStr()
	if dockerName == "" {
		dockerName = "looper"
	}
	dockerCapitalMap["dockerName"] = dockerName
	fmt.Println("请输入需要使用的镜像名称如：busybox")
	var dockerImage string
	dockerImage = getCommandStr()
	if dockerImage == "" {
		dockerImage = "centos"
	}
	dockerCapitalMap["dockerImage"] = dockerImage

	return dockerCapitalMap
}
func getSession(client *ssh.Client) *ssh.Session {
	session, err := client.NewSession()

	if err != nil {
		bar.State()
		log.Fatal("创建ssh session 失败", err)
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
	var level *int
	level = flag.Int("level", 0, "debug level")
	if level == nil {
		log.SetLevel(log.DebugLevel)
	}
	level2 := *level
	switch level2 {
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	getStep("第一步：用户信息初始化")
	countryCapitalMap := initData()
	getBar(4)
	getStep("第二步：获取docker运行参数")
	dockerCapitalMap := getDockerMap()
	getBar(4)
	fmt.Println("本次任务是否开启压缩方式：Y或者N,默认是N")
	var enable string
	enable = getCommandStr()
	if strings.EqualFold(enable, "y") {
		enableCompression = true
	}
	if enableCompression {
		fmt.Println("请选择压缩类型：1：zip；2：tar；3.snappy,默认是：1")
		compression = getCommandStr()
	}

	getStep("第三步：原服务器开始进行初始化")
	client := toObtainSshClient(countryCapitalMap["sourceSshHost"], sshPort, countryCapitalMap["passwd"], sshUser)

	getBar(1)

	getStep("第四步：原服务器任务开始执行创建镜像任务")
	fmt.Println("本次任务是否在宿主机进行容器创建：Y或者N,默认是Y")
	var enable2 string
	enable2 = getCommandStr()
	if strings.EqualFold(enable2, "n") {
		enableCreateDocker = false
	}
	if enableCreateDocker {

		run(client, fmt.Sprintf("docker ps -a| grep %s  | gawk '{cmd=\"docker stop \"$1; system(cmd)}' && docker ps -a| grep %s  | gawk '{cmd=\"docker rm \"$1; system(cmd)}'", dockerCapitalMap["dockerName"],
			dockerCapitalMap["dockerName"]), "docker任务执行失败")
		getBar(5)
		checkImagesSc := fmt.Sprintf("if [ -z `docker images --format  '{{.Repository}}'  %s | grep %s  | awk 'END{print $1}'`]; then echo '不存在'; else echo '存在';  fi",
			dockerCapitalMap["dockerImage"], dockerCapitalMap["dockerImage"])
		s := combinedOutput(client, checkImagesSc, fmt.Sprintf("查找镜像【%s】失败", dockerCapitalMap["dockerImage"]))
		if s == "不存在" {
			searchImageSc := fmt.Sprintf("searchName=`docker search %s --format '{{.Name}}' | grep -E '^%s' | awk 'END{print $1}'` && if [ -z $searchName ]; then echo '不存在'; else echo '存在';  fi", dockerCapitalMap["dockerImage"], dockerCapitalMap["dockerImage"])
			s2 := combinedOutput(client, searchImageSc, "查询镜像失败")
			if s2 == "不存在" {
				log.Errorf("无法找【%s】镜像，请确认镜像是否存在", dockerCapitalMap["dockerImage"])
				os.Exit(0)

			}
			pullImagesSc := fmt.Sprintf("docker pull %s", dockerCapitalMap["dockerImage"])
			run(client, pullImagesSc, fmt.Sprintf("拉取镜像【%s】失败", dockerCapitalMap["dockerImage"]))

		}
		dockerRun := fmt.Sprintf("docker run -d --name %s %s %s  %s ", dockerCapitalMap["dockerName"],
			dockerCapitalMap["dockerRunScript"], dockerCapitalMap["dockerImage"], dockerCapitalMap["script"])
		time.Sleep(5 * time.Second)
		// 2.获取containerId
		containerId = combinedOutput(client, dockerRun, "docker任务执行,创建镜像失败")

		log.Printf("命令输出:%s", strings.TrimSpace(containerId))

		if containerId == "" {
			log.Errorf("获取containerId失败，请重新执行任务")
			os.Exit(0)
		}
		getBar(6)
		getStep("第五步：等待日志生成，方便形成日志差异化")
		if enableCreateDocker {
			log.Printf("等待日志，静候30s....")
			time.Sleep(20 * time.Second)
			log.Printf("继续执行....")
		}
		getBar(2)
	} else {
		log.Println("获取containerId")
		isContainerRun := fmt.Sprintf("docker ps -q -f name=%s", dockerCapitalMap["dockerName"])
		isContainerRunStr := strings.TrimSpace(combinedOutput(client, isContainerRun, "获取containerId失败，请确认原主机容器是否启动正常"))
		if strings.TrimSpace(isContainerRunStr) == "" {
			log.Errorf("原主机：%s容器未正常启动，请确认容器是否启动？", dockerCapitalMap["dockerName"])
			os.Exit(0)

		}
		getContainer := fmt.Sprintf("docker exec -i %s head -1 /proc/self/cgroup|cut -d/ -f3", dockerCapitalMap["dockerName"])
		containerId = strings.TrimSpace(combinedOutput(client, getContainer, "获取containerId失败，请确认原主机容器是否启动正常"))
		if containerId == "" {
			log.Errorf("获取containerId失败，请确认原主机容器是否启动正常")
			os.Exit(0)
		}
		getBar(6 + 2 + 5)
	}
	getStep("第六步：创建checkpoint")
	// 3.创建checkpoint
	checkPoint := fmt.Sprintf("docker checkpoint create %s c1", dockerCapitalMap["dockerName"])
	run(client, checkPoint, "docker任务执行,checkpoint执行失败")
	getBar(3)
	getStep("第七步：开启应用转镜像任务")
	// 4. 转换为镜像
	if containerId == "" {
		log.Errorf("获取containerId失败，请确认原主机容器是否启动正常")
		os.Exit(0)
	}
	commit := fmt.Sprintf("docker commit %s checkpoint", strings.TrimSpace(containerId))
	run(client, commit, "docker任务执行,转换为镜像失败")
	getBar(3)
	getStep("第八步：开启保存本地任务")
	// 5.导出镜像
	save := "docker save -o /opt/checkpoint checkpoint"
	log.Debugln(save)
	run(client, save, "docker任务执行,导出为镜像包失败")
	var checkpointName = "checkpoint"
	if enableCompression {
		lsZip := "ls -lh /opt/checkpoint | awk '{print $5}'"
		output := combinedOutput(client, lsZip, "获取文件大小任务执行失败")
		log.Infof("当前文件大小：%s", strings.TrimSpace(output))
		// 1：zip；2：tar；3.snappy
		switch compression {
		case "2":
			log.Println("开启压缩镜像，使用tar模式压缩")
			compressionSc := "cd /opt/ && tar -zcvf checkpoint.tar.gz checkpoint  > /dev/null&& ls -lh checkpoint.tar.gz |awk '{print $5}'"
			size := combinedOutput(client, compressionSc, "使用tar模式解压缩失败")
			log.Infof("压缩之后的镜像大小:%s", strings.TrimSpace(size))
			checkpointName = "checkpoint.tar.gz"
		case "3":
			log.Println("开启压缩镜像，使用snappy模式的hadoop-snappy压缩")
			compressionSc := "cd /opt/ && snzip -t hadoop-snappy checkpoint  > /dev/null&& ls -lh checkpoint.snappy|awk '{print $5}'"
			size := combinedOutput(client, compressionSc, "使用snappy模式的hadoop-snappy解压缩失败")
			log.Infof("压缩之后的镜像大小:%s", strings.TrimSpace(size))
			checkpointName = "checkpoint.snappy"
		default:
			log.Println("开启压缩镜像，使用zip模式压缩")
			compressionSc := "cd /opt/ && zip -r checkpoint.zip checkpoint  > /dev/null&&ls -lh checkpoint.zip | awk '{print $5}'"
			size := combinedOutput(client, compressionSc, "使用zip模式解压缩失败")
			log.Infof("压缩之后的镜像大小:%s", strings.TrimSpace(size))
			checkpointName = "checkpoint.zip"
		}

	}

	getBar(2)
	getStep("第八步：获取源环境的全量日志")
	// 6.获取当前操作镜像的最后一行日志
	scanLog := fmt.Sprintf("docker logs -f   %s ", dockerCapitalMap["dockerName"])
	log.Println(scanLog)
	tailLogBySource := combinedOutput(client, scanLog, "docker任务执行,查看日志失败")
	getBar(5)
	getStep("第九步：开启发送到迁移主机任务")
	// 7.发送镜像包
	startT := time.Now()
	scp := fmt.Sprintf("scp -q -r -c aes192-cbc -o 'MACs umac-64@openssh.com' /opt/%s root@%s:/opt/%s",
		checkpointName, countryCapitalMap["targetSshHost"], checkpointName)
	log.Debugln(scp)
	run(client, scp, "docker任务执行,发送镜像包失败")
	tc := time.Since(startT)
	log.Printf("发送迁移主机耗时：%s", tc)

	getBar(15)
	getStep("第十步：进入迁移主机服务器，开始执行任务")
	// 8.进入target服务器
	clientTarget := toObtainSshClient(countryCapitalMap["targetSshHost"], sshPort, countryCapitalMap["passwd"], sshUser)
	run(clientTarget, fmt.Sprintf("docker ps -a| grep %s  | gawk '{cmd=\"docker stop \"$1; system("+
		"cmd)}' && docker ps -a| grep %s  | gawk '{cmd=\"docker rm \"$1; system(cmd)}'", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerName"]), "迁移主机上docker任务执行")
	getBar(1)
	getStep("第十一步：进入迁移主机服务器，执行load命令")

	// 9.执行target load命令
	if enableCompression {
		// 1：zip；2：tar；3.snappy
		switch compression {
		case "2":
			log.Println("开启解压镜像，使用tar模式压缩")
			compressionSc := "cd /opt/ && tar -zxvf checkpoint.tar.gz checkpoint"
			run(clientTarget, compressionSc, "开启解压镜像，使用tar模式压缩任务失败")
		case "3":
			log.Println("开启解压镜像，使用snappy模式的hadoop-snappy压缩")
			compressionSc := "cd /opt/ && snzip -d checkpoint.snappy"
			run(clientTarget, compressionSc, "开启解压镜像，使用snappy模式的hadoop-snappy压缩失败")
		default:
			log.Println("开启解压镜像，使用zip模式压缩")
			compressionSc := "unzip -q -n -d /opt/ /opt/checkpoint.zip"
			run(clientTarget, compressionSc, "开启解压镜像，使用zip模式压缩失败")
		}

	}

	load := "docker load -i /opt/checkpoint"
	run(clientTarget, load, "迁移主机上docker任务,导入镜像包失败")
	getBar(4)
	getStep("第十二步：进入迁移主机服务器，开启创建任务")
	// 10.创建任务
	dockerRunTarget := fmt.Sprintf("docker run -d --name %s %s %s  %s && docker stop %s", dockerCapitalMap["dockerName"],
		dockerCapitalMap["dockerRunScript"], "checkpoint", dockerCapitalMap["script"], dockerCapitalMap["dockerName"])
	combo1 := strings.TrimSpace(combinedOutput(clientTarget, dockerRunTarget, "迁移主机上docker任务执行,创建镜像包失败"))
	time.Sleep(5 * time.Second)
	containerIdTarget := strings.Split(combo1, "\n")[0]
	getBar(10)
	getStep("第十三步：进入源主机，进行拷贝checkpoint到迁移节点")
	// 11.拷贝checkpoint到目的节点
	var checkpointsName string
	if enableCompression {
		sizeSc := fmt.Sprintf("du -h --max-depth=0  /var/lib/docker/containers/%s/checkpoints/c1 |  awk '{print $1}'", strings.TrimSpace(containerId))
		log.Infof("当前checkpoints的压缩点位大小为：%s", combinedOutput(client, sizeSc, "获取checkpoints大小失败"))
		// 1：zip；2：tar；3.snappy
		switch compression {
		case "2":
			log.Println("开启checkpoints镜像，使用tar模式压缩")
			compressionSc := fmt.Sprintf("cd /var/lib/docker/containers/%s/checkpoints/ && tar -zcvf c1.tar.gz c1 > /dev/null && ls -lh c1.tar.gz |awk '{print $5}'", strings.TrimSpace(containerId))
			size := combinedOutput(client, compressionSc, "使用tar模式压缩失败")
			log.Infof("压缩之后的checkpoints大小:%s", strings.TrimSpace(size))
			checkpointsName = "c1.tar.gz"
		case "3":
			log.Println("开启checkpoints镜像，使用snappy模式的hadoop-snappy压缩")
			compressionSc := fmt.Sprintf("cd /var/lib/docker/containers/%s/checkpoints/ && tar cf - c1   | snzip -t hadoop-snappy  > archive.tar.sz  && ls -lh archive.tar.sz|awk '{print $5}'", strings.TrimSpace(containerId))
			size := combinedOutput(client, compressionSc, "使用snappy模式的hadoop-snappy压缩失败")
			log.Infof("压缩之后的checkpoints大小:%s", strings.TrimSpace(size))
			checkpointsName = "archive.tar.sz"
		default:
			log.Println("开启checkpoints镜像，使用zip模式压缩")
			compressionSc := fmt.Sprintf("cd /var/lib/docker/containers/%s/checkpoints/ && zip -r c1.zip c1 > /dev/null &&ls -lh c1.zip | awk '{print $5}'", strings.TrimSpace(containerId))
			size := combinedOutput(client, compressionSc, "使用zip模式压缩失败")
			log.Infof("压缩之后的checkpoints大小:%s", strings.TrimSpace(size))
			checkpointsName = "c1.zip"
		}
		scpCheckPoint := fmt.Sprintf(
			"scp -r /var/lib/docker/containers/%s/checkpoints/%s/ root@%s:/var/lib/docker/containers/%s/checkpoints/",
			strings.TrimSpace(containerId), checkpointsName, countryCapitalMap["targetSshHost"],
			strings.TrimSpace(containerIdTarget))
		run(client, scpCheckPoint, "迁移主机上docker任务执行,拷贝checkpoint到目的节点失败")
	} else {
		scpCheckPoint := fmt.Sprintf(
			"scp -r /var/lib/docker/containers/%s/checkpoints/c1/ root@%s:/var/lib/docker/containers/%s/checkpoints/",
			strings.TrimSpace(containerId), countryCapitalMap["targetSshHost"], strings.TrimSpace(containerIdTarget))
		run(client, scpCheckPoint, "迁移主机上docker任务执行,拷贝checkpoint到目的节点失败")
	}
	time.Sleep(5 * time.Second)
	if enableCompression {
		// 1：zip；2：tar；3.snappy
		switch compression {
		case "2":
			log.Println("开启解压checkpoints，使用tar模式压缩")
			compressionSc := fmt.Sprintf("cd /var/lib/docker/containers/%s/checkpoints/ && tar -zxvf c1.tar.gz c1 &&rm -rf /var/lib/docker/containers/%s/checkpoints/c1.tar.gz", strings.TrimSpace(containerIdTarget), strings.TrimSpace(containerIdTarget))
			run(clientTarget, compressionSc, "开启解压checkpoints，使用tar模式压缩任务失败")
		case "3":
			log.Println("开启解压checkpoints，使用snappy模式的hadoop-snappy压缩")
			compressionSc := fmt.Sprintf("cd /var/lib/docker/containers/%s/checkpoints/ && snzip -dc archive.tar.sz | tar xf -&&rm -rf /var/lib/docker/containers/%s/checkpoints/archive.tar.sz", strings.TrimSpace(containerIdTarget), strings.TrimSpace(containerIdTarget))
			run(clientTarget, compressionSc, "开启解压checkpoints，使用snappy模式的hadoop-snappy压缩失败")
		default:
			log.Println("开启解压checkpoints，使用zip模式压缩")
			compressionSc := fmt.Sprintf("unzip -q -n -d /var/lib/docker/containers/%s/checkpoints/ /var/lib/docker/containers/%s/checkpoints/c1.zip&&rm -rf /var/lib/docker/containers/%s/checkpoints/c1.zip", strings.TrimSpace(containerIdTarget), strings.TrimSpace(containerIdTarget), strings.TrimSpace(containerIdTarget))
			run(clientTarget, compressionSc, "开启解压checkpoints，使用zip模式压缩失败")
		}
	}

	// 13. 恢复位点差
	checkPointTarget := fmt.Sprintf("docker start --checkpoint c1 %s", strings.TrimSpace(containerIdTarget))

	getBar(20)
	getStep("第十三步：进入迁移主机，恢复位点差")
	run(clientTarget, checkPointTarget, "迁移主机上docker任务执行,启动checkPoint失败")
	getBar(10)
	getStep("第十四步：进入迁移主机，进行日志全量获取")
	// 13.再次检查进程日志正常，接着上次创建checkpoint的时间点打印
	scanLogTarget := fmt.Sprintf(
		"docker logs --tail  all %s ", strings.TrimSpace(containerIdTarget))
	time.Sleep(5 * time.Second)
	combo := combinedOutput(clientTarget, scanLogTarget, "迁移主机上docker任务执行,日志获取失败")
	getBar(5)
	log.Println("打印原主机日志：")

	fmt.Println(tailLogBySource)

	log.Printf("打印迁移主机日志：")
	fmt.Println(string(combo))
	getStep("第十四步：进行缓存清理")
	s1 := "docker images|grep none|awk '{print $3 }'"
	s2 := "docker images|grep none|awk '{print $3 }'|xargs docker rmi >/dev/null"
	output := combinedOutput(client, s1, "源主机清理任务失败")
	if output != "" {
		run(client, s2, "源主机清理任务失败")

	}
	output1 := combinedOutput(clientTarget, s1, "源主机清理任务失败")
	if output1 != "" {
		run(clientTarget, s2, "源主机清理任务失败")
	}

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

	getStep("bye")
}
func combinedOutput(client *ssh.Client, script string, error string) string {
	session := getSession(client)
	log.Debugln(script)
	combo, err := session.CombinedOutput(script)
	if err != nil {
		log.Errorf("%s:", error, err)
		os.Exit(0)
	}
	out := string(combo)
	closeSession(session)
	return strings.TrimSpace(out)
}
func output(client *ssh.Client, script string, error string) string {
	session := getSession(client)
	log.Debugln(script)
	combo, err := session.CombinedOutput(script)
	if err != nil {
		log.Errorf("%s:", error, err)
		os.Exit(0)
	}
	out := string(combo)
	closeSession(session)
	return strings.TrimSpace(out)
}
func run(client *ssh.Client, script string, error string) {
	log.Debugln(script)
	session := getSession(client)
	err := session.Run(script)
	if err != nil {
		log.Errorf("%s", error, err)
		os.Exit(0)
	}
	closeSession(session)
}
func getBar(num int) {
	log.Println("当前进度.................")
	bar.Add(num)
	log.Println("任务继续调度.................")

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
		// bar.load()
		log.Errorln("创建ssh client 失败", err)
	}

	return sshClient
}

func getStep(message string) string {
	return fmt.Sprintf("=================================%s====================================================", message)

}
