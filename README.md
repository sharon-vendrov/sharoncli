# sharoncli tools

The tool assume that: 
1. codefreshcli tool is already installed and configured.
2. after the cluster creation, the user will create a pipeline which will be linked to the created cluster.
3. Docker deamon is installed for the "on-prem" option

examples:
sharoncli create runtime --cluster-name kubernetes-admin@kind --cloud-provider on-prem
sharoncli test runtime --name "default/project"

[![asciicast](https://asciinema.org/a/Dic6DbdELMRPFOuj7xlUvuSSx.svg)](https://asciinema.org/a/Dic6DbdELMRPFOuj7xlUvuSSx)
