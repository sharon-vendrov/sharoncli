package logic

import "testing"

func TestExecutePipeline(t *testing.T) {
	err := ExecutePipeline("default/MyPipeline")
	if err != nil {
		t.Fail()
	}

}

func TestListPipelines(t *testing.T) {
	err := ListPipelines()
	if err != nil {
		t.Fail()
	}

}
