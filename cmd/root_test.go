package cmd

import "testing"

func TestExecute(t *testing.T) {
	Execute()
}

func TestRootCmd(t *testing.T) {
	// should panic because of nil values
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	rootCmd.Run(nil, nil)
}
