package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

var (
	net         string
	clusterName string
	force       bool
	provider    *cluster.Provider
)

var createCmd = &cobra.Command{
	Use:  "create",
	Long: `Spin up a kind cluster for testing argocd.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clusters, err := provider.List()
		if err != nil {
			return err
		}
		for _, c := range clusters {
			if c == clusterName {
				fmt.Printf("%s kind cluster already exists\n", clusterName)
				if !force {
					kubeCfg, err := provider.KubeConfig(clusterName, false)
					if err != nil {
						return err
					}
					fmt.Println(kubeCfg)
					return nil
				}
				provider.Delete(clusterName, "")
			}
		}
		if err := kindCreate(); err != nil {
			return err
		}
		kubeCfg, err := provider.KubeConfig(clusterName, false)
		if err != nil {
			return err
		}
		fmt.Println(kubeCfg)
		return nil
	},
}

func kindCreate() error {
	return provider.Create(clusterName,
		cluster.CreateWithDisplaySalutation(true),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithRawConfig(clusterConfig()),
	)
}

func clusterConfig() []byte {
	cfg := fmt.Sprintf(`kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  ipFamily: %s
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP`, net)

	return []byte(cfg)
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&net, "net", "ipv6", "kind cluster IP Family. Valid choices are ipv6 or ipv4.")
	createCmd.Flags().StringVar(&clusterName, "name", "test-argocd", "name")
	createCmd.Flags().BoolVarP(&force, "force", "f", false, "Force create cluster.  This will delete existing cluster and recreate.")
	if provider == nil {
		provider = cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	}
}
