package internal

type Cdata struct {
	Cdata string `xml:",cdata"`
}

type Code struct {
	PythonCode Cdata `xml:"python_code"`

	// Empty for Python projects.
	ScratchDescription Cdata `xml:"scratch_description"`
}
