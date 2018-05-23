package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func main() {

}

func Args() {

	args := os.Args

	fmt.Println(args)

	programName := args[0]
	fmt.Printf("The binary name is: %s \n", programName)

	otherArgs := args[1:]
	fmt.Println(otherArgs)

	for idx, arg := range otherArgs {
		fmt.Printf("Arg %d = %s \n", idx, arg)
	}

}

func GetEnv() {
	connStr := os.Getenv("DB_CONN")
	log.Printf("Connection string: %s\n", connStr)
}

func SetEnv() {

	key := "MYROOT"
	os.Setenv(key, "/home/wxl/work/go")

	val := GetEnvDefault(key, "/usr/local/go")

	log.Println("The value is :" + val)

	os.Unsetenv(key)

	val = GetEnvDefault(key, "/usr/local/go")

	log.Println("The default value is :" + val)

}

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

func GetPid() {

	pid := os.Getpid()
	fmt.Printf("Process PID: %d \n", pid)

	prc := exec.Command("ps", "-p", strconv.Itoa(pid), "-v")
	out, err := prc.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

}

func Signal() {

	sChan := make(chan os.Signal, 1)

	signal.Notify(sChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	exitChan := make(chan int)
	go func() {
		signal := <-sChan
		switch signal {
		case syscall.SIGHUP:
			fmt.Println("closed")
			exitChan <- 0

		case syscall.SIGINT:
			fmt.Println(" CTRL+C")
			exitChan <- 1

		case syscall.SIGTERM:
			fmt.Println("kill SIGTERM  executed")
			exitChan <- 1

		case syscall.SIGKILL:
			fmt.Println("SIGKILL handler")
			exitChan <- 1

		case syscall.SIGQUIT:
			fmt.Println("kill SIGQUIT  executed ")
			exitChan <- 1
		}
	}()

	code := <-exitChan
	os.Exit(code)
}

func osExec() {
	prc := exec.Command("ls", "-a")
	out := bytes.NewBuffer([]byte{})
	prc.Stdout = out
	err := prc.Run()
	if err != nil {
		fmt.Println(err)
	}

	if prc.ProcessState.Success() {
		fmt.Println("Output:\n")
		fmt.Println(out.String())
	}
}

func ProcWait() {

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "timeout"
	} else {
		cmd = "sleep"
	}

	proc := exec.Command(cmd, "1")
	proc.Start()

	proc.Wait()

	fmt.Printf("PID: %d\n", proc.ProcessState.Pid())
	fmt.Printf("Process run time: %dms\n", proc.ProcessState.SystemTime()/time.Microsecond)
	fmt.Printf("Exited result : %t\n", proc.ProcessState.Success())
}

func procOutput() {

	var cmd string

	if runtime.GOOS == "windows" {
		cmd = "dir"
	} else {
		cmd = "ls"
	}

	proc := exec.Command(cmd)

	buff, err := proc.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buff))
}

func Pipe() {
	cmd := []string{"go", "run", "sample.go"}

	proc := exec.Command(cmd[0], cmd[1], cmd[2])

	stdin, _ := proc.StdinPipe()
	defer stdin.Close()

	stdout, _ := proc.StdoutPipe()
	defer stdout.Close()

	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			fmt.Println("Program says:" + s.Text())
		}
	}()

	proc.Start()

	fmt.Println("Writing input")
	io.WriteString(stdin, "Hello\n")
	io.WriteString(stdin, "Golang\n")
	io.WriteString(stdin, "is awesome\n")

	time.Sleep(time.Second * 2)

	proc.Process.Kill()
}

var writer *os.File

func osFile() {
	var err error
	path := "/tmp/log"
	pwd, _ := os.Getwd()
	fmt.Printf("pwd : %s")
	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("create dir:%s", path)
	fileName := fmt.Sprintf("test_%d.log", time.Now().Unix())
	filePath := path + fileName
	writer, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	io.WriteString(writer, "logging\n")
	writer.Close()

	os.Chdir("/tmp/log") // cd
	pwd, _ = os.Getwd()
	fmt.Printf("pwd : %s")

}

