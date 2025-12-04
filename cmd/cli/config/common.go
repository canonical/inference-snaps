package config

import "github.com/spf13/cobra"

const groupID = "config"

func Group(title string) *cobra.Group {
	return &cobra.Group{
		ID:    groupID,
		Title: title,
	}
}
