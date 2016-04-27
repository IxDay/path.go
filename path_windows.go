package path

func (self Path) NormCase() Path {
	return Path(strings.ToLower(string(self)))
}

func (self Path) ExpandUser() Path {
	return self
}

func (self Path) ExpandVars() Path {
	return self
}
