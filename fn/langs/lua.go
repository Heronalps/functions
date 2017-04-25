package langs

type LuaLangHelper struct {
	BaseHelper
}

func (lh *LuaLangHelper) Entrypoint() string {
	return "th func.lua --type cuda"
}

func (lh *LuaLangHelper) HasPreBuild() bool {
	return false
}

// PreBuild for Go builds the binary so the final image can be as small as possible
func (lh *LuaLangHelper) PreBuild() error {
	return nil
}

func (lh *LuaLangHelper) AfterBuild() error {
	return nil
}
