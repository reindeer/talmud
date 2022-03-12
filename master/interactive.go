//go:build !darwin && !ios && !dbus
// +build !darwin,!ios,!dbus

package master

func get() string {
	return scan()
}

func save() {
}
