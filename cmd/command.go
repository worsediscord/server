package cmd

type Command interface {
	Name() string
	Description() string
	Parse(args []string) error
	Run() error
}

// StringSliceValue implements flag.Value.
type StringSliceValue []string

func (s *StringSliceValue) String() string {
	return ""
}

func (s *StringSliceValue) Set(value string) error {
	*s = append(*s, value)
	return nil
}
