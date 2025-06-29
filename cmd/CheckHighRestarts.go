package cmd 

import(
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
	"strings"
	"strconv"
	"github.com/fatih/color"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var CheckHighRestarts = &cobra.Command{
	Use: "check-high-restarts",
	Short: "Get the pods with high restarts",
	Run:func(cmd *cobra.Command,args []string){
		color.Yellow("Checking....")
		fmcd := exec.Command("kubectl","get","pods","-A")
		out,err := fmcd.Output()
		if err!=nil{
			fmt.Println(err)
		}
		lines:=strings.Split(string(out),"\n")
		fmt.Println("NAMESPACE       NAME                                              READY     STATUS      RESTARTS       AGE")
		total := 0
		for _,line := range lines{
			fields := strings.Fields(line)
			if len(fields) < 5 || fields[0] == "NAMESPACE" {
				continue
			}
			restarts,_ :=  strconv.Atoi(fields[4])
			if(restarts>=10){
				fmt.Println(line)
				total = total + restarts
			}
		}
		Counter := prometheus.NewCounter(
				prometheus.CounterOpts{
					Name:"k8s_total_restarts",
					Help:"Total number of restarts > 10",
				},
			)
			Counter.Add(float64(total))
			_= push.New("http://localhost:9091","k8s_cli_job").Collector(Counter).Grouping("env","dev").Push();
	},
}

func init(){
	rootCmd.AddCommand(CheckHighRestarts)
}