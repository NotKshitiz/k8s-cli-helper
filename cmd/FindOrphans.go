package cmd

import (
	"os/exec"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var FindOrphans = &cobra.Command{
	Use:   "find-orphans",
	Short: "Lists services that have no active endpoints (orphaned services)",
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(" Scanning for orphaned services (no endpoints)...")

		kcmd := exec.Command("kubectl", "get", "endpoints", "-A")
		out, err := kcmd.Output()
		if err != nil {
			color.Red(" Failed to run kubectl: %v", err)
			return
		}
		total := 0
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)

			if len(fields) < 4 || fields[2] != "<none>" || fields[0] == "NAMESPACE" {
				continue
			}

			namespace := fields[0]
			service := fields[1]
			total = total + 1
			color.Red("Orphaned service: %s (Namespace: %s)", service, namespace)

		}
		Guage := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:"k8s_total_orphans",
				Help:"Total Orphans services with no endpoints",
			},
		)
		Guage.Set(float64(total))
		_ = push.New("http://localhost:9091","k8s_job").
			Collector(Guage).
			Grouping("env","dev").Push();

		color.Green("Done scanning for orphaned services.")
	},
}

func init() {
	rootCmd.AddCommand(FindOrphans)
}
