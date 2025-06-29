package cmd
import(
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
	"strings"
	"github.com/fatih/color"
	"bufio"
	"os"
)


var FixCrashPods = &cobra.Command{
	Use: "fix-crashpods",
	Short: "Restarts the crash pods after deleting them",
	Run: func(cmd *cobra.Command,args []string){
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Do you want to delete and restart the crash pods(Standalone pods will be deleted forever) y/n")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input)) 
		if input != "y" {
			color.Yellow("Skipping pod deletion. No changes made.")
			return
		}
		fmt.Println("Restarting Pods....")
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
				fmt.Printf("Deleting pod %s (namespace: %s)...\n", podName, namespace)
				err := exec.Command("kubectl", "delete", "pod", podName, "-n", namespace).Run()
				if err != nil {
					color.Red("Failed to delete pod %s: %v", podName, err)
				} else {
					color.Green("Deleted pod %s", podName)
				}

			}
		}
		color.Green("\nPod cleanup complete. Kubernetes will attempt to restart the pods.")
	
},

}

func init() {
	rootCmd.AddCommand(FixCrashPods)
}