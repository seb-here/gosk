package test

import (
	"fmt"
	"testing"

	"github.com/seb-here/gosk"
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
		t.Fatal(err)
	}
	skill, err := kernel.ImportSkill("FunSkill")
	if err != nil {
		t.Fatal(err)
	}
	sf := skill["Joke"]
	response, tokenCount, err := sf("Engineer", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Token count: %v\n", tokenCount)
	fmt.Printf("%v\n", response)
}
