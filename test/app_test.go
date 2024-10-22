package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"sigs.k8s.io/yaml"
)

func TestArgocd_Application(t *testing.T) {
	testCases := []struct {
		fileName string
		name     string
		want     bool
	}{
		{"app-guestbook.yaml", "validApp", true},
		{"invalid-app-guestbook.yaml", "missingNamespace", false},
		{"atlas-kibana.yml", "helm", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := os.ReadFile(tc.fileName)
			if err != nil {
				t.Error(err)
			}
			app := &v1alpha1.Application{}
			if err := yaml.Unmarshal(b, app); err != nil {
				t.Error(err)
			}
			result := t.Run("CheckDestination", checkDestination(app))
			if result != tc.want {
				t.Error()
			}
			t.Run("CheckSource", checkSource(app))
		})
	}
}

func checkDestination(app *v1alpha1.Application) func(t *testing.T) {
	return func(t *testing.T) {
		dest := app.Spec.Destination
		if dest.Name != "" && dest.Server != "" {
			t.Error("Both application.spec.destination.name and application.spec.destination.server are set, only set one")
		}
		if dest.Name == "" && dest.Server == "" {
			t.Error("must set application.spec.destination.name or application.spec.destination.server")
		}
		if dest.Namespace == "" {
			t.Error("must set application.spec.destination.namespace")
		}
	}
}

func checkSource(app *v1alpha1.Application) func(t *testing.T) {
	return func(t *testing.T) {
		source := app.Spec.Source
		if source.IsZero() {
			t.Error("must supply application.spec.source")
			return
		}
		if source.IsHelm() {
			fmt.Println(source.Helm.Values)
		}
	}
}
