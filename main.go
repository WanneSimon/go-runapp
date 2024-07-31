package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

/* 启动参数 */
type AppStart struct {
	AppName string   // 记录pid的文件
	Main    string   // 主命令
	Force   bool     // 是否忽略pid文件是否存在，强行启动
	Args    []string // 参数参数
}

/* 停止参数 */
type AppStop struct {
	AppName string // 记录pid的文件
}

/* 重启参数 */
type AppRestart struct {
	AppName string // 记录pid的文件
}

/*
 * 启动：
 * myapp starts -exec javaw -args "-jar xxx.jar"
 * 停止
 * myapp stop
 * 重启
 * myapp restart
 */
func main() {
	allArgs := os.Args
	if len(allArgs) < 3 {
		usage()
		return
	}

	var appName string = allArgs[1]
	var mode string = allArgs[2]

	if mode == "start" || mode == "starts" {
		var arg AppStart
		arg.AppName = appName

		// 是否是强制启动
		arg.Force = mode == "starts"

		num := len(allArgs)
		if num < 4 {
			usage()
			return
		}
		arg.Main = allArgs[3]

		if num > 4 {
			arg.Args = allArgs[4:]
		} else {
			arg.Args = make([]string, 0)
		}

		startApp(arg)
	} else if mode == "stop" {
		stopApp(AppStop{
			AppName: appName,
		})
	} else if mode == "restart" {
		restartApp(AppRestart{
			AppName: appName,
		})
	} else {
		usage()
	}
}

func usage() {
	s := OsSeparator()
	str := s + "Example:" + s + "  myapp start|starts javaw -jar myapp.jar" + s + "  myapp stop" + s + "  myapp restart"
	log.Println(str)
}

/*
 * 启动：
 * myapp start  javaw -args "-jar xxx.jar"
 */
func startApp(arg AppStart) {
	// 检查参数
	if arg.AppName == "" {
		log.Println("缺少APP名称")
		os.Exit(1)
	}

	var pidFile = getFidFile(arg.AppName)
	// 非强制启动时检查数据文件是否存在
	if !arg.Force {
		info, err := os.Stat(pidFile)
		if !os.IsNotExist(err) {
			log.Println("APP名称已存在：" + info.Name())
			os.Exit(1)
		}
	}

	if arg.Main == "" {
		log.Println("缺少执行命令")
		os.Exit(1)
	}

	// 执行
	execCMD(arg.Main, arg.Args, pidFile)
}

/* 执行命令
 * @param main 主命令
 * @param args 参数字符串
 * @param dataFile 存储数据的文件
 */
func execCMD(main string, args []string, dataFile string) {
	log.Println("Exec")
	log.Println(main, args)

	// 拆分参数字符串为数组
	// todo 参数有空格的情况
	var argsArr []string = make([]string, 0)
	if args != nil {
		argsArr = args[:]
	}

	cmd := exec.Command(main, argsArr...)
	// err2 := cmd.Run()
	err2 := cmd.Start()
	if err2 != nil {
		log.Println("命令执行失败")
		log.Println(err2)
		return
	}

	// 保存信息
	pid := cmd.Process.Pid
	s := OsSeparator()

	argsStr := strings.Join(args, s)
	var savedata = strconv.Itoa(pid) + s + main + s + argsStr
	err := os.WriteFile(dataFile, []byte(savedata), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func OsSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

/*
 * 停止
 * myapp stop
 */
func stopApp(arg AppStop) {
	// 检查参数
	if arg.AppName == "" {
		log.Println("缺少APP名称")
		os.Exit(1)
	}

	var pidFile = getFidFile(arg.AppName)
	pid, _ := loadPid(arg.AppName)

	su := stopPid(pid, false)
	//time.Sleep(time.Second * 1)
	if su {
		log.Println(arg.AppName + " 已停止")
		return
	}
	log.Println("停止失败，重试")

	// 检查后，再次尝试
	su = stopPid(pid, false)
	//time.Sleep(time.Second * 3)
	if su {
		log.Println(arg.AppName + " 已停止")
		return
	}
	log.Println("停止失败，强行停止")

	// 检查后，强制停止
	su = stopPid(pid, true)
	if !su {
		log.Println("无法停止")
		return
	}
	log.Println(arg.AppName + " 已停止")

	// 移除文件
	removeFile(pidFile)
	log.Println("文件已移除 " + pidFile)
}

/*
 * 重启
 * myapp restart
 */
func restartApp(arg AppRestart) {
	// 检查参数
	if arg.AppName == "" {
		log.Println("缺少APP名称")
		os.Exit(1)
	}

	_, appStart := loadPid(arg.AppName)
	pidFile := getFidFile(arg.AppName)

	// 停止
	stopApp(AppStop{
		AppName: arg.AppName,
	})
	// 启动
	execCMD(appStart.Main, appStart.Args, pidFile)
}

func stopPid(pid string, force bool) bool {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		arr := make([]string, 0)
		if force {
			arr = append(arr, "/f")
		}
		arr = append(arr, "/pid")
		arr = append(arr, pid)

		cmd = exec.Command("taskkill", arr...)
	} else {
		arr := make([]string, 0)
		if force {
			arr = append(arr, "-9")
		}
		arr = append(arr, pid)
		cmd = exec.Command("kill", arr...)
	}

	err := cmd.Run()
	// err := cmd.Start()
	// err := cmd.Wait()
	if err != nil {
		if err.Error() == "exit status 128" {
			return true
		}
		log.Println("停止失败: " + err.Error())
		return false
	}
	return true
}
func removeFile(file string) {
	err2 := os.Remove(file)
	if err2 != nil {
		log.Println("PID文件删除失败！")
		log.Println(err2)
	}

	_, err := os.Stat(file)
	if os.IsExist(err) {
		os.Remove(file)
	}
}

func loadPid(appName string) (string, AppStart) {
	var pidFile = getFidFile(appName)
	info, err := os.Stat(pidFile)
	if os.IsNotExist(err) {
		log.Println("APP信息不存在：" + info.Name())
		os.Exit(1)
	}

	bytes, err := os.ReadFile(pidFile)
	if err != nil {
		log.Println("程序信息读取失败！")
		os.Exit(1)
	}

	data := string(bytes)
	lines := strings.Split(data, OsSeparator())
	pid := lines[0]
	main := lines[1]
	args := make([]string, 0)
	if len(lines) > 2 {
		args = lines[2:]
	}

	return pid, AppStart{
		AppName: appName,
		Main:    main,
		Args:    args,
	}
}

func getFidFile(appName string) string {
	return appName + ".pid"
}
