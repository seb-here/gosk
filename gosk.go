package gosk

import (
	"embed"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mfmayer/gopenai"
	"github.com/mfmayer/gosk/internal/skillconfig"
	"github.com/mfmayer/gosk/utils"
)

//go:embed skills/*
var embeddedSkillsDir embed.FS

// SemanticKernel
type SemanticKernel struct {
	chatClient *gopenai.ChatClient
}

type newKernelOption func(*newKernelOptions)

type newKernelOptions struct {
	openAIKey string
}

// WithOpenAIKey to use this OpenAI key when creating a new semantic kernel, otherwise it's tried to get the key from "OPENAI_API_KEY" environment variable or .env file in current working directory
func WithOpenAIKey(key string) newKernelOption {
	return func(opt *newKernelOptions) {
		opt.openAIKey = key
	}
}

// NewKernel creates new kernel and tries to retrieve the OpenAI key from "OPENAI_API_KEY" environment variable or .env file in current working directory
func NewKernel(opts ...newKernelOption) (kernel *SemanticKernel, err error) {
	options := &newKernelOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.openAIKey == "" {
		options.openAIKey, err = utils.GetOpenAIKey()
		if err != nil {
			return
		}
	}
	cClient := gopenai.NewChatClient(options.openAIKey)
	kernel = &SemanticKernel{
		chatClient: cClient,
	}
	return
}

func (k *SemanticKernel) ImportSkill(name string) (skill Skill, err error) {
	fs, err := fs.Sub(embeddedSkillsDir, "skills")
	if err != nil {
		return
	}
	return k.importSkill(fs, name)
}

func (k *SemanticKernel) importSkill(fsys fs.FS, skillName string) (skill Skill, err error) {
	skill = Skill{}
	err = fs.WalkDir(fsys, skillName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != skillName {
			// read skill prompt template
			skprompt, err := template.ParseFS(fsys, filepath.Join(path, "skprompt.txt"))
			_ = skprompt
			if err != nil {
				return err
			}
			// read skill config
			jsonFile, err := fsys.Open(filepath.Join(path, "config.json"))
			if err != nil {
				return err
			}
			defer jsonFile.Close()
			sConfigBytes, _ := ioutil.ReadAll(jsonFile)
			sConfig := skillconfig.DefaultSkillConfig()
			err = json.Unmarshal(sConfigBytes, &sConfig)
			if err != nil {
				return err
			}

			// use path base name as skill function name
			skillFunctionName := filepath.Base(path)
			skill[skillFunctionName] = k.createSkillFunction(skprompt, sConfig)
			//templates[dirName] = subDirTemplate
			_ = skillFunctionName
		}
		return nil
	})
	return
}

func (k *SemanticKernel) createSkillFunction(template *template.Template, config skillconfig.SkillConfig) (skillFunc SkillFunction) {
	paramMap := map[string]*string{}
	paramArray := []*string{}
	for _, p := range config.Input.Parameters {
		param := p.DefaultValue
		paramMap[p.Name] = &param
		paramArray = append(paramArray, &param)
	}

	skillFunc = func(parameters ...string) string {
		for i, param := range parameters {
			if i < len(paramArray) {
				*paramArray[i] = param
			}
		}
		template.Execute(os.Stdout, paramMap)
		return ""
	}
	return
}
