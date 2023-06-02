package gosk

type SkillFunction func(parameters ...string) (string, int, error)

type Skill map[string]SkillFunction
