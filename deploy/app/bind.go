package main

// Asset mock for no bootstrap mode
func Asset(name string) ([]byte, error) {
	return []byte{}, nil
}

// AssetDir mock for no bootstrap mode
func AssetDir(name string) ([]string, error) {
	return []string{}, nil
}

// RestoreAssets mock for no bootstrap mode
func RestoreAssets(dir, name string) error {
	return nil
}
