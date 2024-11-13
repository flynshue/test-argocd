package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"
)

func renderTemplate(filename string, data any) (manifest string, err error) {
	tpl, err := template.ParseFS(fs, filename)
	if err != nil {
		return "", err
	}
	manifest = fmt.Sprintf("/tmp/%s.yaml", path.Base(filename))
	f, err := os.Create(manifest)
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(f, data); err != nil {
		return "", err
	}
	return manifest, nil
}

func kubectl(args ...string) (stdout, stderr []byte, err error) {
	c := exec.Command("kubectl", args...)
	op, err := c.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	ep, err := c.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, nil, err
	}
	stdout, _ = io.ReadAll(op)
	stderr, _ = io.ReadAll(ep)
	return stdout, stderr, c.Wait()
}

func currentContext() (string, error) {
	out, _, err := kubectl("config", "current-context")
	if err != nil {
		return "", err
	}
	b := bytes.TrimRight(out, "\n")
	return string(b), nil
}

func exportKubeCfg() error {
	err := provider.ExportKubeConfig(clusterName, "", false)
	if err != nil {
		return err
	}
	context, err := currentContext()
	if err != nil {
		return err
	}
	fmt.Printf("using context %s\n", context)
	return nil
}

func kubectlApply(fileName, namespace string) error {
	args := []string{"apply", "-f", fileName}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	stdout, stderr, err := kubectl(args...)
	fmt.Println(string(stdout))
	if err != nil {
		fmt.Println(string(stderr))
		return err
	}
	return nil
}

func kubectlDeleteF(fileName, namespace string) error {
	args := []string{"delete", "-f", fileName}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	stdout, stderr, err := kubectl(args...)
	fmt.Println(string(stdout))
	if err != nil {
		fmt.Println(string(stderr))
		return err
	}
	return nil
}
