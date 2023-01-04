//go:build !darwin && !ios && !dbus

package provider

import "github.com/reindeer/talmud/internal/terminal"

func Get() string {
	return terminal.Scan()
}

func Save() {
}
