package internal

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func Test_Handle(t *testing.T) {
	testHandleCases := []struct {
		name    string
		input   string
		want    Command
		wantErr bool
	}{
		{
			name:  "empty input should send to default search provider",
			input: "",
			want: Command{
				Action:   "lookup",
				Name:     "*",
				Location: fmt.Sprintf(DefaultSearchProvider, ""),
			},
		},
		{
			name:  "add shortcut: 'add gh https://github.com'",
			input: "add gh https://github.com",
			want: Command{
				Action:   "add",
				Name:     "gh",
				Location: "https://github.com",
			},
		},
		{
			name:    "add shortcut without location: 'add gh'",
			input:   "add gh",
			wantErr: true,
		},
		{
			name:  "add search shortcut: 'add so https://stackvoerflow.com?q=%s'",
			input: "add so https://stackoverflow.com?q=%s",
			want: Command{
				Action:   "add",
				Name:     "so",
				Location: "https://stackoverflow.com?q=%s",
			},
		},
		{
			name:  "remove shortcut: 'remove gh'",
			input: "remove gh",
			want: Command{
				Action: "remove",
				Name:   "gh",
			},
		},
		{
			name:    "remove command with incorrect number of args: 'remove'",
			input:   "remove",
			wantErr: true,
		},
		{
			name:  "forward search to search provider: 'golang unit testing frameworks'",
			input: "golang unit testing frameworks",
			want: Command{
				Action:   "lookup",
				Name:     "*",
				Location: fmt.Sprintf(DefaultSearchProvider, "golang unit testing frameworks"),
			},
		},
		{
			name:  "visit a shortcut: 'fb'",
			input: "fb",
			want: Command{
				Action:   "lookup",
				Name:     "fb",
				Location: "https://facebook.com",
			},
		},
		{
			name:  "visit a search shortcut: 'go net/http'",
			input: "go net/http",
			want: Command{
				Action:   "lookup",
				Name:     "go",
				Location: "https://godoc.org/net/http",
			},
		},
	}

	for _, testCase := range testHandleCases {
		tt := testCase
		t.Run(tt.name, func(t *testing.T) {
			cm := &CommandHandler{
				Mutex: &sync.Mutex{},
				Shortcuts: map[string]string{
					"fb": "https://facebook.com",
					"go": "https://godoc.org/%s",
				},
			}

			got, err := cm.Handle(tt.input)

			if tt.wantErr != (err != nil) {
				t.Errorf("wantErr: %v got '%v'", tt.wantErr, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInput('%s') = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
