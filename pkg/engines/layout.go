package engines

import "fmt"

type Layout struct {
	Symlink string `yaml:"symlink"`
}

func (layout Layout) Validate(target string) error {
	if target == "" {
		return fmt.Errorf("layout: required field is not set: target")
	}
	if layout.Symlink == "" {
		return fmt.Errorf("layout %q: required field is not set: symlink", target)
	}
	return nil
}

func ValidateLayout(layout map[string]Layout) error {
	for target, entry := range layout {
		if err := entry.Validate(target); err != nil {
			return err
		}
	}
	return nil
}
