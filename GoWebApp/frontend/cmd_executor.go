package main

import (
	"fmt"
	"os/exec"
)

type Params struct {
	Level           int    `json:"level"`
	SourceSshHost   string `json:"sourceSshHost"`
	TargetSshHost   string `json:"targetSshHost"`
	Passwd          string `json:"passwd"`
	DockerRunScript string `json:"dockerRunScript"`
	Script          string `json:"script"`
	DockerName      string `json:"dockerName"`
	DockerImage     string `json:"dockerImage"`
	Enable          string `json:"enable"`
	Compression     int    `json:"compression"`
	Enable2         string `json:"enable2"`
}

func ExecuteCmd(params Params) (string, error) {
	cmd := buildCmd(params)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func buildCmd(params Params) string {
	// 根据参数拼接命令行
	cmd := fmt.Sprintf("run-criu --level %d --sourceSshHost %s --targetSshHost %s --passwd %s --dockerRunScript %s --script %s --dockerName %s --dockerImage %s --enable %s --compression %d --enable2 %s",
		params.Level,
		params.SourceSshHost,
		params.TargetSshHost,
		params.Passwd,
		params.DockerRunScript,
		params.Script,
		params.DockerName,
		params.DockerImage,
		params.Enable,
		params.Compression,
		params.Enable2,
	)
	return cmd
}
