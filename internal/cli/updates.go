package cli

import (
	"fmt"

	"context"

	"charm.land/lipgloss/v2"
	"github.com/google/go-github/v84/github"
)

func CheckForUpdate(v string) error {
	client := github.NewClient(nil)
	ctx := context.Background()
	release, _, err := client.Repositories.GetLatestRelease(ctx, "JanMalch", "snips")
	if err != nil {
		return err
	}
	latestTag := "v" + v
	if latestTag == release.GetTagName() {
		fmt.Printf("v%s is the lastest snips version.\n", v)
		return nil
	}

	url := release.GetHTMLURL()
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Align(lipgloss.Left)

	fmt.Println(style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Inline(true).Bold(true).Render("A new snips version is available!\n"),
			fmt.Sprintf("%s -> %s", latestTag, release.GetTagName()),
			lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("12")).Render(url),
		),
	))
	return nil
}
