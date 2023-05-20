package gosk

type SkillFunction func(parameters ...string) string

type Skill map[string]SkillFunction
