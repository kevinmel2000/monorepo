package simplecli

/*
Simple package for simple CLI.
This package provide all structure that is needed to build a simple CLI application.
Command with sub-command supported with this following rules:
1. Command that don't have sub command will run immediately:
	eg: ./megazord restart
	But ./megaozrd restart help will return help instead
	output:
		restart:
			- restart apps
2. Command with sub-commands will return usage help if sub-command not specified:
	eg: ./megazord cli
	output:
		cli:
			- description of cli
		usage:
			- service restart > restart apps
TODO: Simplify CLI logic
*/

import (
	"context"
	"fmt"
	"strings"
)

// Command struct to describe cli command spesification
type Command struct {
	Name        string
	Description string
	Command     func(context.Context, []string) error

	// sub command
	SubCommand       []SubCommand          // list of command that is allowed for 1st args
	subCommandMap    map[string]SubCommand // internal map of sub-command
	subCommandLength int
}

type SubCommand struct {
	Name        string
	Description string
}

var commands map[string]*Command

const help = "help"

// RegisterCommands for set map of cli commands
func RegisterCommands(cmds ...*Command) error {
	if commands == nil {
		// flush will create empty map of command
		flush()
	}
	for _, command := range cmds {
		commandName := strings.ToLower(command.Name)
		if _, ok := commands[commandName]; ok {
			return fmt.Errorf("Command %s is duplicated", command.Name)
		}

		// map sub-command to internal sub-command map
		command.subCommandLength = len(command.SubCommand)
		if command.subCommandLength > 0 {
			command.subCommandMap = make(map[string]SubCommand)
			for _, sub := range command.SubCommand {
				subName := strings.ToLower(sub.Name)
				// sub command cannot take help command
				if subName == help {
					return fmt.Errorf("Cannot register help sub-command in %s command", command.Name)
				}
				// check if sub-command duplicate
				if _, ok := command.subCommandMap[subName]; ok {
					return fmt.Errorf("Sub-command %s is duplicated in %s command", sub.Name, command.Name)
				}
				command.subCommandMap[subName] = sub
			}
		}
		commands[commandName] = command
	}
	return nil
}

func Run(args []string) {
	err := run(args)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func run(args []string) (err error) {
	argsLength := len(args)
	// this is most likely not happened, but check anyway
	if argsLength <= 0 {
		return err
	}

	// check if args1 is help, if yes then describe all the command and description
	if strings.ToLower(args[1]) == help {
		help := fmt.Sprintf("COMMANDS:\n") // COMMANDS:
		for _, command := range commands {
			help += fmt.Sprintf("%s - %s\n", command.Name, command.Description) // ex: restart - restart service
		}
		fmt.Print(help)
		return err
	}

	// check if the command actually exists
	command, ok := commands[args[1]]
	if !ok {
		return fmt.Errorf("Command %s is not exists, please use help to identify commands", args[1])
	}

	var (
		sub SubCommand

		showHelp     bool
		showUsage    bool
		isSubCommand bool
		offset       = 2 // offset for trimming arguments passed to command function
	)
	// show usage if have sub-command
	if command.subCommandLength > 0 {
		showUsage = true
	}

	// if command have sub-command but not send its sub-command. ex: ./megazord cli
	if argsLength <= 2 && showUsage {
		showHelp = true
	}

	if argsLength >= 3 {
		if args[2] == help { // check if argument is help
			showHelp = true
		} else if command.subCommandLength > 0 { // check if sub-command is actually exists

			sub, ok = command.subCommandMap[args[2]]
			if !ok {
				return fmt.Errorf("Sub-command %s is not exists", args[2])
			}
			offset = 3
			isSubCommand = true
		}

		// if using sub-command but need help
		if argsLength > 3 && args[3] == help {
			showHelp = true
		}
	}

	// execute if help is not executed
	if !showHelp {
		// trim args and send to command function
		args = args[offset:]
		err = command.Command(context.TODO(), args)
		if err != nil {
			// show help if error
			showHelp = true
			err = fmt.Errorf("CLI Error: %s", err.Error())
		}
	}

	if showHelp {
		showCommandHelp(command, sub, showUsage, isSubCommand)
	}
	return err
}

func showCommandHelp(c *Command, sub SubCommand, showUsage, isSubCommand bool) {
	var help string
	if !isSubCommand {
		help += fmt.Sprintf("%s:\n", c.Name)             // 					command:
		help += fmt.Sprintf("  - %s\n\n", c.Description) //   						- description of command
	} else {
		help += fmt.Sprintf("%s:\n", sub.Name)             // 					command:
		help += fmt.Sprintf("  - %s\n\n", sub.Description) //   					- description of command
		showUsage = false
	}

	if showUsage {
		help += fmt.Sprintf("%s:\n", "usage") // 								usage:
		for _, sub := range c.SubCommand {
			help += fmt.Sprintf("  - %s > %s\n", sub.Name, sub.Description) //		- subcommand > description
		}
	}
	fmt.Print(help + "\n")
}

// create emtry map of cli command
func flush() {
	commands = make(map[string]*Command)
}
