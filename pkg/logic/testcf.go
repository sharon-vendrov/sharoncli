package logic

import (
	"fmt"
	"os"

	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/codefresh-io/go-sdk/pkg/utils"
)

// ExecutePipeline execute CF pipeline
func ExecutePipeline(pipelineName string) error {
	path := fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
	options, err := utils.ReadAuthContext(path, "")
	if err != nil {
		fmt.Println("Failed to read codefresh config file")
		return (err)
	}
	clientOptions := codefresh.ClientOptions{Host: options.URL,
		Auth: codefresh.AuthOptions{Token: options.Token}}
	cf := codefresh.New(&clientOptions)
	runOptions := codefresh.RunOptions{Branch: "string"}
	resp, err := cf.Pipelines().Run(pipelineName, &runOptions)
	if err != nil {
		fmt.Println("Failed to get run pipeline")
		return (err)
	}

	fmt.Println(resp)
	return nil
}

// ListPipelines lists all pipelines
func ListPipelines() error {
	path := fmt.Sprintf("%s/.cfconfig", os.Getenv("HOME"))
	options, err := utils.ReadAuthContext(path, "")
	if err != nil {
		fmt.Println("Failed to read codefresh config file")
		return (err)
	}
	clientOptions := codefresh.ClientOptions{Host: options.URL,
		Auth: codefresh.AuthOptions{Token: options.Token}}
	cf := codefresh.New(&clientOptions)
	pipelines, err := cf.Pipelines().List()
	if err != nil {
		fmt.Println("Failed to get Pipelines from Codefresh API")
		return (err)
	}
	for _, p := range pipelines {
		fmt.Printf("Pipeline: %+v\n\n", p)
	}

	return nil
}
