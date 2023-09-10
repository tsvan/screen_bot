package internal

type Action struct {
	opts ActionOpts
}

type ActionOpts struct {
}

func NewAction() *Action {
	return &Action{}
}


