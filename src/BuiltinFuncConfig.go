package main

type BuiltinFuncConfig struct {
	Name     string
	Args     []string
	Targets  []string
	LinkReqs []string
	Eval     EvalFn
}

func (cfg BuiltinFuncConfig) Allowed() bool {
	if cfg.Targets == nil {
		return true
	}

	for _, t := range cfg.Targets {
		if t == TARGET {
			return true
		}
	}

	return false
}
