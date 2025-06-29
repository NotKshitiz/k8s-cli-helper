package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
	"strings"
	"github.com/fatih/color"
	"os"
	"log"
)

var GetLogsCrashLoopBackOff = &cobra.Command{
	Use: "logs-crashpods",
	Short:"Give logs of the CrashPods",
	Run:func(cmd *cobra.Command,args []string){
		fmt.Println("Logs for the CrashLoopBackOff Pods....")
		kcmd := exec.Command("kubectl","get","pods","-A")
		out,err := kcmd.Output()
		if err!= nil {
			fmt.Println(err)
		}
		lines:=strings.Split(string(out),"\n")
		for _ , line := range lines{
			if strings.Contains(line,"CrashLoopBackOff"){
			fields := strings.Fields(line)
			namespace := fields[0]
			podName := fields[1]
			color.Red("CrashLoop Pod Found: %s (Namespace: %s)\n", podName, namespace)
			fmt.Println("-----------------------------------------------------")
			mcmd := exec.Command("kubectl", "logs", "-n", namespace, podName)
			out,err:= mcmd.Output()
			if err!= nil{
				fmt.Println(err)
			}
			logs_pod := strings.Split(string(out),"\n")
			logfile,err:= os.OpenFile("logs/mycli.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
			if err!=nil{
				fmt.Println(err)
				return
			}
			defer logfile.Close()
			logger := log.New(logfile,"",log.LstdFlags|log.Lshortfile)
			for _,logs_line := range logs_pod{
				logger.Println("[CRASHPODS]",logs_line)
				fmt.Println(logs_line)
			}
		}
	}
},
}
func init() {
	rootCmd.AddCommand(GetLogsCrashLoopBackOff)
}
