package simplecli

import (
	"context"
	"testing"
)

func TestRegisterCLI(t *testing.T) {
	// flush before test
	flush()
	cases := []struct {
		command   *Command
		expectErr bool
	}{
		{
			command: &Command{
				Name:        "test1",
				Description: "for help test1",
				Command: func(ctx context.Context, args []string) error {
					return nil
				},
			},
			expectErr: false,
		},
		{
			command: &Command{
				Name:        "test2",
				Description: "for help test2",
				Command: func(ctx context.Context, args []string) error {
					return nil
				},
			},
			expectErr: false,
		},
		// command should already exists
		{
			command: &Command{
				Name:        "test1",
				Description: "for help test1",
				Command: func(ctx context.Context, args []string) error {
					return nil
				},
			},
			expectErr: true,
		},
		// command with sub-command
		{
			command: &Command{
				Name:        "test3",
				Description: "for help test3",
				Command: func(ctx context.Context, args []string) error {
					return nil
				},
				SubCommand: []SubCommand{
					{
						Name:        "testing",
						Description: "this is for testing",
					},
				},
			},
			expectErr: false,
		},
		// command with sub-command, but failed because same sub-command
		{
			command: &Command{
				Name:        "test4",
				Description: "for help test4",
				Command: func(ctx context.Context, args []string) error {
					return nil
				},
				SubCommand: []SubCommand{
					{
						Name:        "testing",
						Description: "this is for testing",
					},
					{
						Name:        "testing",
						Description: "this is for testing bruh",
					},
				},
			},
			expectErr: true,
		},
	}

	for _, val := range cases {
		err := RegisterCommands(val.command)
		if !val.expectErr && err != nil {
			t.Error("Failed to register command", err.Error())
		}
	}
}

// func TestRun(t *testing.T) {
// 	cases := []struct {
// 		args      []string
// 		expectErr bool
// 	}{
// 		{
// 			args:      []string{"test1"},
// 			expectErr: false,
// 		},
// 		{
// 			args:      []string{"test1", "help"},
// 			expectErr: false,
// 		},
// 		{
// 			args:      []string{"help"},
// 			expectErr: false,
// 		},
// 	}

// 	// cli commands
// 	clis := []*Command{
// 		&Command{
// 			Name:        "test1",
// 			Description: "for help test1",
// 			Command: func(ctx context.Context, args []string) error {
// 				return nil
// 			},
// 		},
// 		&Command{
// 			Name:        "test2",
// 			Description: "for help test2",
// 			Command: func(ctx context.Context, args []string) error {
// 				return nil
// 			},
// 		},
// 	}
// }
