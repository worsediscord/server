package main

func main() {
	cmd := NewRootCmd(NewStartCmd())

	if err := cmd.Parse(nil); err != nil {
		panic(err)
	}

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
