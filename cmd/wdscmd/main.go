package main

func main() {
	rootCmd := NewRootCmd("wdscli")

	rootCmd.AddSubcommands(
		NewStartCmd("start", rootCmd.Name()+" "),
	)

	if err := rootCmd.Parse(nil); err != nil {
		panic(err)
	}

	if err := rootCmd.Run(); err != nil {
		panic(err)
	}
}
