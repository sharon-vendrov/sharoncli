module github.com/sharon-vendrov/sharoncli

go 1.13

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.7 // indirect
	github.com/codefresh-io/go-sdk v0.16.0
	github.com/codefresh-io/venona/venonactl v0.0.0-20190815092312-094052ae2519
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.4.0
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // indirect
	golang.org/x/sys v0.0.0-20190922100055-0a153f010e69 // indirect
	google.golang.org/appengine v1.5.0 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20190920012459-5008bf6f8cd6 // indirect
	sigs.k8s.io/kind v0.5.1
)

replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
