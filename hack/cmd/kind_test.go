package cmd

import (
	"testing"

	"gopkg.in/yaml.v3"
)

type kindCluster struct {
	KindNetworking `yaml:"networking"`
}

type KindNetworking struct {
	IPFamily string `yaml:"ipFamily"`
}

func Test_clusterConfig(t *testing.T) {
	testCases := []struct {
		name string
		net  string
	}{
		{"ipv6Config", "ipv6"},
		{"ipv4Config", "ipv4"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			net = tc.net
			clusterCfg := kindCluster{}
			err := yaml.Unmarshal(clusterConfig(), &clusterCfg)
			if err != nil {
				t.Error(err)
			}
			got := clusterCfg.IPFamily
			if got != tc.net {
				t.Errorf("got %s, wanted %s", got, tc.net)
			}
		})
	}
}

func Test_exportKubeCfg(t *testing.T) {
	if err := exportKubeCfg(); err != nil {
		t.Error(err)
	}
}

func Test_installIngressNginx(t *testing.T) {
	if err := installIngressNginx(); err != nil {
		t.Error(err)
	}
}
