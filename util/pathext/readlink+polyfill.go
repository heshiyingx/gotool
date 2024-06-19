//go:build windows

package pathext

func ReadLink(name string) (string, error) {
	return name, nil
}
