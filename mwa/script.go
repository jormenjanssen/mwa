package main

func NewScriptAction(script string) func() error {

	if script == "" {
		return NoScriptErrorFunc()
	}

	return func() error {
		return ExecuteScript(script)
	}
}

func NoScriptErrorFunc() func() error {
	return func() error {
		return nil
	}
}
