
package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var checkCrashloopsCmd = &cobra.Command{
	Use:   "checkCrashloops",
	Short: "Detect pods in CrashLoopBackOff state",
	Run: func(cmd *cobra.Command, args []string) {
		
		fmt.Println("Checking for CrashLoopBackOff pods...")
		kcmd := exec.Command("kubectl","get","pods","-A")
		out,err := kcmd.Output()
		if err!= nil {
			fmt.Println(err)
		}
		lines:=strings.Split(string(out),"\n")
		count_CrashPods := 0
		fmt.Println("NAMESPACE       NAME                                              READY      STATUS            RESTARTS       AGE")	
		for _ , line := range lines{
			if strings.Contains(line,"CrashLoopBackOff"){
			fmt.Println(line)
			count_CrashPods++
		}
		}
		fmt.Printf("We have detected %v Pods with CrashLoopBackOff status\n",count_CrashPods)
		gauge := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "k8s_crashloop_pods",
				Help: "Number of CrashLoopBackOff Pods",
			},
		)
		gauge.Set(float64(count_CrashPods))
		_= push.New("http://localhost:9091","k8s_cli_job").Collector(gauge).Grouping("env","dev").Push();
	},
}


func init() {
	rootCmd.AddCommand(checkCrashloopsCmd)
}
