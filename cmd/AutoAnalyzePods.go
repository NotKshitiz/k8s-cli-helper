package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
	"strings"
	"github.com/fatih/color"
)

var AutoAnalyzePods = &cobra.Command{
	Use: "auto-analyze",
	Short: "Auto analyzes the CrashOffLoop Pods",
	Run:func(cmd *cobra.Command,args []string){
		fmt.Println("Analyzing CrashOffLoop Pods")
		kcmd := exec.Command("kubectl","get","pods","-A")
		out,err := kcmd.Output()
		if err!= nil {
			fmt.Println(err)
		}
		lines:=strings.Split(string(out),"\n")
		for _,line := range lines{
			if strings.Contains(line,"CrashLoopBackOff"){
			fields := strings.Fields(line)
			namespace := fields[0]
			podName := fields[1]
			mcmd := exec.Command("kubectl","describe","pod",podName,"-n",namespace)
			out,err:= mcmd.Output()
			if err!=nil{
				fmt.Println(err)
			}
			color.Yellow("CrashLoop Pod: %s (Namespace: %s)\n", podName, namespace)
			Describe_Lines := strings.Split(string(out),"\n")
			inEvents := false
			for _,describe_line := range Describe_Lines{
							line := strings.TrimSpace(describe_line)

				if strings.HasPrefix(line, "Events:") {
					color.Blue("\n Events:")
					inEvents = true
					continue
				}

				if inEvents {
					if line == "" {
						inEvents = false
						continue
					}
					fmt.Println("  " + line)
				}
				if strings.Contains(describe_line,"Container ID:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"Image:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"Image ID:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"State:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"Reason:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"Exit Code:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"Restart Count:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"BackOff:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"Pulling:"){
					fmt.Println(describe_line)
				}
				if strings.Contains(describe_line,"ImagePullBackOff"){
					fmt.Println(describe_line)
				}
			}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(AutoAnalyzePods)
}