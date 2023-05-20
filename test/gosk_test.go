package test

import (
	"testing"

	"github.com/mfmayer/gosk"
)

func TestKernel(t *testing.T) {
	kernel, err := gosk.NewKernel()
	if err != nil {
		t.Error(err)
	}
	t.Logf("hello kernel: %v", kernel)
}

func TestSkillImport(t *testing.T) {
	kernel, err := gosk.NewKernel()
	if err != nil {
		t.Error(err)
	}
	skill, err := kernel.ImportSkill("FunSkill")
	if err != nil {
		t.Fatal(err)
	}
	sf := skill["Joke"]
	sf("testinput", "teststyle")
}
