/*
Copyright © 2019 Sharon Vendrov <sharon.vendrov1@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/create"
	"sigs.k8s.io/kind/pkg/util"

	"github.com/spf13/viper"
)

type flagpole struct {
	Name          string
	Config        string
	ImageName     string
	Retain        bool
	Wait          time.Duration
	CloudProvider string
}

type venonaInstallCmdOptions struct {
	dryRun                 bool
	clusterNameInCodefresh string
	kube                   struct {
		namespace string
		inCluster bool
		context   string
	}
	storageClass string
	venona       struct {
		version string
	}
	setDefaultRuntime             bool
	installOnlyRuntimeEnvironment bool
	skipRuntimeInstallation       bool
	runtimeEnvironmentName        string
	kubernetesRunnerType          bool
}

var flags = &flagpole{}
var installCmdOptions = &venonaInstallCmdOptions{}

// runtimeCmd represents the runtime command
var runtimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "TODO",
	Long:  `TODO`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runE(flags, cmd, args)
	},
}

func init() {
	viper.BindEnv("kube-namespace", "KUBE_NAMESPACE")
	viper.BindEnv("kube-context", "KUBE_CONTEXT")

	runtimeCmd.PersistentFlags().StringVar(&kubeConfigPath, "kube-config-path", viper.GetString("kubeconfig"), "Path to kubeconfig file (default is $HOME/.kube/config) [$KUBECONFIG]")

	runtimeCmd.Flags().StringVar(&flags.Name, "name", cluster.DefaultName, "cluster context name")
	runtimeCmd.Flags().StringVar(&flags.Config, "config", "", "path to a kind config file")
	runtimeCmd.Flags().StringVar(&flags.ImageName, "image", "", "node docker image to use for booting the cluster")
	runtimeCmd.Flags().BoolVar(&flags.Retain, "retain", false, "retain nodes for debugging when cluster creation fails")
	runtimeCmd.Flags().DurationVar(&flags.Wait, "wait", time.Duration(120)*time.Second, "Wait for control plane node to be ready (default 120s)")
	runtimeCmd.Flags().StringVar(&flags.CloudProvider, "cloud-provider", "on-prem", "Define cloud provider")

	runtimeCmd.Flags().StringVar(&installCmdOptions.clusterNameInCodefresh, "cluster-name", "", "cluster name (if not passed runtime-environment will be created cluster-less)")
	runtimeCmd.Flags().StringVar(&installCmdOptions.venona.version, "venona-version", "", "Version of venona to install (default is the latest)")
	runtimeCmd.Flags().StringVar(&installCmdOptions.runtimeEnvironmentName, "runtime-environment", "", "if --skip-runtime-installation set, will try to configure venona on current runtime-environment")
	runtimeCmd.Flags().StringVar(&installCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which venona should be installed [$KUBE_NAMESPACE]")
	runtimeCmd.Flags().StringVar(&installCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which venona should be installed (default is current-context) [$KUBE_CONTEXT]")
	runtimeCmd.Flags().StringVar(&installCmdOptions.storageClass, "storage-class", "", "Set a name of your custom storage class, note: this will not install volume provisioning components")

	runtimeCmd.Flags().BoolVar(&installCmdOptions.skipRuntimeInstallation, "skip-runtime-installation", false, "Set flag if you already have a configured runtime-environment, add --runtime-environment flag with name")
	runtimeCmd.Flags().BoolVar(&installCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if venona is been installed from inside a cluster")
	runtimeCmd.Flags().BoolVar(&installCmdOptions.installOnlyRuntimeEnvironment, "only-runtime-environment", false, "Set to true to onlky configure namespace as runtime-environment for Codefresh")
	runtimeCmd.Flags().BoolVar(&installCmdOptions.dryRun, "dry-run", false, "Set to true to simulate installation")
	runtimeCmd.Flags().BoolVar(&installCmdOptions.setDefaultRuntime, "set-default", false, "Mark the install runtime-environment as default one after installation")
	runtimeCmd.Flags().BoolVar(&installCmdOptions.kubernetesRunnerType, "kubernetes-runner-type", false, "Set the runner type to kubernetes (alpha feature)")

	createCmd.AddCommand(runtimeCmd)
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	switch flags.CloudProvider {
	case "on-prem":

		// Check if the cluster name already exists
		known, err := cluster.IsKnown(flags.Name)
		if err != nil {
			return err
		}
		if known {
			return errors.Errorf("a cluster with the name %q already exists", flags.Name)
		}

		// create a cluster context and create the cluster
		ctx := cluster.NewContext(flags.Name)
		fmt.Printf("Creating cluster %q ...\n", flags.Name)
		if err = ctx.Create(
			create.WithConfigFile(flags.Config),
			create.WithNodeImage(flags.ImageName),
			create.Retain(flags.Retain),
			create.WaitForReady(flags.Wait),
		); err != nil {
			if utilErrors, ok := err.(util.Errors); ok {
				for _, problem := range utilErrors.Errors() {
					log.Error(problem)
				}
				return errors.New("aborting due to invalid configuration")
			}
			return errors.Wrap(err, "failed to create cluster")
		}
		installCmdOptions.kube.context = "kubernetes-admin@kind"
		installCmdOptions.clusterNameInCodefresh = "kubernetes-admin@kind"
		homePath := os.Getenv("HOME")
		kubeConfigPath = homePath + "/.kube/kind-config-kind"
	default:
		return errors.Errorf("The cloud-provider isn't supported")
	}

	installVenona(*installCmdOptions)

	return nil
}
