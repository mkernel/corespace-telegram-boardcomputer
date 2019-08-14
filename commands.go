package main

var List = []Command{
}

func Run(args []string) error {
	if(len(args) == 0) {
		args=[]string{"help"}
	}
	for _,cmd := range List {
		if(cmd.Command() == args[0]) {
			var cmdArgs=args[1:];
			return cmd.Execute(cmdArgs);
		}
	}
	return nil;
}