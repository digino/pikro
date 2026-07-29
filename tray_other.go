//go:build !windows && !darwin

package main

// runWithTray is a no-op on non-Windows platforms; the server runs directly.
func runWithTray(_ int, startServer func() error) {
	if err := startServer(); err != nil {
		panic(err)
	}
}
