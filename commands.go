package main

var commands = []command{}

func run(args []string) error {
	if len(args) == 0 {
		args = []string{"help"}
	}
	for _, cmd := range commands {
		if cmd.Command() == args[0] {
			var cmdArgs = args[1:]
			return cmd.Execute(cmdArgs)
		}
	}
	return nil
}
