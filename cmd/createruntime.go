/*
Copyright © 2019 Sharon Vendrov <sharon.vendrov1@gmail.com>

*/
package cmd

import (
	"flag"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/create"
	"sigs.k8s.io/kind/pkg/util"

	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type flagpole struct {
	Name      string
	Config    string
	ImageName string
	Retain    bool
	Wait      time.Duration
}

var flags = &flagpole{}

// runtimeCmd represents the runtime command
var runtimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runE(flags, cmd, args)
	},
}

func init() {
	runtimeCmd.Flags().StringVar(&flags.Name, "name", cluster.DefaultName, "cluster context name")
	runtimeCmd.Flags().StringVar(&flags.Config, "config", "", "path to a kind config file")
	runtimeCmd.Flags().StringVar(&flags.ImageName, "image", "", "node docker image to use for booting the cluster")
	runtimeCmd.Flags().BoolVar(&flags.Retain, "retain", false, "retain nodes for debugging when cluster creation fails")
	runtimeCmd.Flags().DurationVar(&flags.Wait, "wait", time.Duration(120)*time.Second, "Wait for control plane node to be ready (default 120s)")
	createCmd.AddCommand(runtimeCmd)
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
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

	listPods(*ctx)

	return nil
}

func listPods(ctx cluster.Context) {
	var kubeconfig *string
	if home := ctx.KubeConfigPath(); home != "" {
		kubeconfig = flag.String("kubeconfig", home, "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "example-xxxxx"
		_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
		if apierr.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*apierr.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}
