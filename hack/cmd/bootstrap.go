package cmd

import (
	"embed"
	"fmt"
)

//go:embed templates/*
var fs embed.FS

func bootstrapCluster() error {
	if err := installIngressNginx(); err != nil {
		return err
	}
	return nil
}

func installIngressNginx() error {
	if err := exportKubeCfg(); err != nil {
		return err
	}
	ipFamily := "IPv6"
	if net == "ipv4" {
		ipFamily = "IPv4"
	}
	f, err := renderTemplate("templates/kind-ingress-nginx.tpl", ipFamily)
	if err != nil {
		return err
	}
	err = kubectlApply(f, "ingress-nginx")
	if err != nil {
		return err
	}
	validateArgs := []string{"wait", "--namespace", "ingress-nginx",
		"--for=condition=ready", "pod", "--selector=app.kubernetes.io/component=controller",
		"--timeout=90s",
	}
	out, stderr, err := kubectl(validateArgs...)
	fmt.Println(string(out))
	if err != nil {
		fmt.Println(string(stderr))
		return err
	}
	return nil
}
