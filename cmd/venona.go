package cmd

/*
Copyright 2019 The Codefresh Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/codefresh-io/venona/venonactl/pkg/plugins"
	"github.com/codefresh-io/venona/venonactl/pkg/store"
)

// installCmd represents the install command
func installVenona(installCmdOptions venonaInstallCmdOptions) {
	s := store.GetStore()
	lgr := createLogger("Install", verbose)
	buildBasicStore(lgr)
	extendStoreWithCodefershClient(lgr)
	extendStoreWithKubeClient(lgr)

	builder := plugins.NewBuilder(lgr)
	isDefault := isUsingDefaultStorageClass(installCmdOptions.storageClass)

	builderInstallOpt := &plugins.InstallOptions{
		CodefreshHost:         s.CodefreshAPI.Host,
		CodefreshToken:        s.CodefreshAPI.Token,
		MarkAsDefault:         installCmdOptions.setDefaultRuntime,
		StorageClass:          installCmdOptions.storageClass,
		IsDefaultStorageClass: isDefault,
		DryRun:                installCmdOptions.dryRun,
		KubernetesRunnerType:  installCmdOptions.kubernetesRunnerType,
	}

	if installCmdOptions.kubernetesRunnerType {
		builder.Add(plugins.EnginePluginType)
	}

	if isDefault {
		builderInstallOpt.StorageClass = plugins.DefaultStorageClassNamePrefix
	}

	if installCmdOptions.kube.context == "" {
		config := clientcmd.GetConfigFromFileOrDie(s.KubernetesAPI.ConfigPath)
		installCmdOptions.kube.context = config.CurrentContext
		lgr.Debug("Kube Context is not set, using current context", "Kube-Context-Name", installCmdOptions.kube.context)
	}
	if installCmdOptions.kube.namespace == "" {
		installCmdOptions.kube.namespace = "default"
	}

	s.KubernetesAPI.InCluster = installCmdOptions.kube.inCluster

	s.KubernetesAPI.ContextName = installCmdOptions.kube.context
	s.KubernetesAPI.Namespace = installCmdOptions.kube.namespace

	if installCmdOptions.dryRun {
		s.DryRun = installCmdOptions.dryRun
		lgr.Info("Running in dry-run mode")
	}
	if installCmdOptions.venona.version != "" {
		version := installCmdOptions.venona.version
		lgr.Info("Version set manually", "version", version)
		s.Image.Tag = version
		s.Version.Latest.Version = version
	}
	s.ClusterInCodefresh = installCmdOptions.clusterNameInCodefresh
	if installCmdOptions.installOnlyRuntimeEnvironment == true && installCmdOptions.skipRuntimeInstallation == true {
		dieOnError(fmt.Errorf("Cannot use both flags skip-runtime-installation and only-runtime-environment"))
	}
	if installCmdOptions.installOnlyRuntimeEnvironment == true {
		builder.Add(plugins.RuntimeEnvironmentPluginType)
	} else if installCmdOptions.skipRuntimeInstallation == true {
		if installCmdOptions.runtimeEnvironmentName == "" {
			dieOnError(fmt.Errorf("runtime-environment flag is required when using flag skip-runtime-installation"))
		}
		s.RuntimeEnvironment = installCmdOptions.runtimeEnvironmentName
		lgr.Info("Skipping installation of runtime environment, installing venona only")
		builder.Add(plugins.VenonaPluginType)
	} else {
		builder.
			Add(plugins.RuntimeEnvironmentPluginType).
			Add(plugins.VenonaPluginType)
	}
	if isDefault {
		builder.Add(plugins.VolumeProvisionerPluginType)
	} else {
		lgr.Info("Custom StorageClass is set, skipping installation of default volume provisioner")
	}

	builderInstallOpt.ClusterName = s.KubernetesAPI.ContextName
	builderInstallOpt.RegisterWithAgent = true
	if s.ClusterInCodefresh != "" {
		builderInstallOpt.ClusterName = s.ClusterInCodefresh
		builderInstallOpt.RegisterWithAgent = false
	}
	builderInstallOpt.KubeBuilder = getKubeClientBuilder(builderInstallOpt.ClusterName, s.KubernetesAPI.Namespace, s.KubernetesAPI.ConfigPath, s.KubernetesAPI.InCluster)
	builderInstallOpt.ClusterNamespace = s.KubernetesAPI.Namespace

	values := s.BuildValues()
	var err error
	for _, p := range builder.Get() {
		values, err = p.Install(builderInstallOpt, values)
		if err != nil {
			dieOnError(err)
		}
	}
	lgr.Info("Installation completed Successfully")
}
