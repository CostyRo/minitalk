package main

type InputHandler struct{}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

func (h *InputHandler) Complete(input string, promptFn func(string) (string, error)) (string, error) {
	parensOpen, parensClose := 0, 0
	brackOpen, brackClose := 0, 0

	count := func(s string) {
		for _, c := range s {
			switch c {
			case '(':
				parensOpen++
			case ')':
				parensClose++
			case '[':
				brackOpen++
			case ']':
				brackClose++
			}
		}
	}

	count(input)

	for parensOpen > parensClose || brackOpen > brackClose {
		cont, err := promptFn("... ")
		if err != nil {
			return "", err
		}
		input += "\n" + cont
		count(cont)
	}

	return input, nil
}
