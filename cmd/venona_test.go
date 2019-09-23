package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInstallVenona(t *testing.T) {
	installCmdOptions.kube.context = "kubernetes-admin@kind"
	installCmdOptions.clusterNameInCodefresh = "kubernetes-admin@kind"
	installCmdOptions.kube.namespace = viper.GetString("kube-namespace")

	homePath := os.Getenv("HOME")
	kubeConfigPath = homePath + "/.kube/kind-config-kind"

	installVenona(*installCmdOptions)

	// TODO get runtime environment and validate creation
}
