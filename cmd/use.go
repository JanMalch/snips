package cmd

import (
	"os"
	"snips/internal/cnfg"
	"snips/internal/snippets"
	"snips/internal/use"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useCmdCopy bool
var useCmdOverwrite bool
var useCmdFile string

var useCmd = &cobra.Command{
	Use:     "use [flags] [id]",
	Aliases: []string{"u"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Use a snippet",
	Long: `Use a snippet.

Examples:
  snips use fun/data.kt -f my-data.kt    # appends to my-data.kt
  snips use fun/data.kt -o -f my-data.kt # overwrites my-data.kt
  snips use fun/data.kt -f.              # . will show destination suggestions based on the group and name
  snips use utils/strings.ts             # prints to stdout
  snips use api/example.json -c          # copies to system clipboard`,
	Run: func(cmd *cobra.Command, args []string) {
		var arg string
		if len(args) == 0 {
			selected, err := useSnippetTree()
			cobra.CheckErr(err)
			arg = selected.String()
		} else {
			arg = args[0]
		}
		cobra.CheckErr(use.Use(arg, useCmdFile, viper.GetBool(cnfg.KEY_APPLY_ALLOW_DIRTY), useCmdOverwrite, useCmdCopy, os.Stdout))
	},
}

func lastPathSegment(s string) string {
	segments := strings.Split(s, "/")
	return segments[len(segments)-1]
}

var groupColor = tcell.ColorGray
var fileTextStyle = tcell.StyleDefault.Bold(true)

func useSnippetTree() (snippets.Id, error) {
	ids, err := snippets.ListAll()
	if err != nil {
		return snippets.Id{}, err
	}

	var selected snippets.Id
	app := tview.NewApplication()
	treeRoot := tview.NewTreeNode("Snippets").SetSelectable(false)
	buildTree(treeRoot, ids)

	inputField := tview.NewInputField().
		SetLabel("Search: ")

	tree := tview.NewTreeView().
		SetRoot(treeRoot).
		SetCurrentNode(treeRoot)

	inputField.SetChangedFunc(func(text string) {
		if text == "" {
			buildTree(treeRoot, ids)
			return
		}
		matches := fuzzy.FindFromNoSort(text, ids)
		filteredIds := make(snippets.Ids, 0)
		for _, match := range matches {
			filteredIds = append(filteredIds, ids[match.Index])
		}
		buildTree(treeRoot, filteredIds)
	})
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
			app.SetFocus(tree)
			return nil
		}
		return event
	})
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab, tcell.KeyBacktab, tcell.KeyEscape:
			app.SetFocus(tree)
		case tcell.KeyEnter:
			id := findFirstFile(tree.GetRoot())
			if id != nil {
				selected = *id
				app.Stop()
			}
		}
	})

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// TODO: KeyUp focus inputField, when at top of list
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			app.SetFocus(inputField)
			return nil
		} else if event.Key() == tcell.KeyLeft {
			node := tree.GetCurrentNode()
			if node != nil && node.GetReference() == nil {
				node.Collapse()
			}
			return nil
		} else if event.Key() == tcell.KeyRight {
			node := tree.GetCurrentNode()
			if node != nil && node.GetReference() == nil {
				node.Expand()
			}
			return nil
		}
		return event
	})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref != nil {
			selected = ref.(snippets.Id)
			app.Stop()
		}
	})

	container := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0).
		SetBorders(false).
		AddItem(inputField, 0, 0, 1, 1, 1, 0, true).
		AddItem(tree, 1, 0, 1, 1, 0, 0, false)

	err = app.SetRoot(container, true).EnablePaste(true).SetFocus(container).Run()
	if err != nil {
		return snippets.Id{}, err
	}
	return selected, nil
}

func buildTree(root *tview.TreeNode, ids snippets.Ids) {
	root.ClearChildren()
	nodes := make(map[string]*tview.TreeNode, 0)
	nodes[""] = root
	for _, id := range ids {
		parentCrumb := ""
		for _, crumb := range id.Breadcrumbs() {
			_, ok := nodes[crumb]
			if !ok {
				node := tview.NewTreeNode(lastPathSegment(crumb)).SetColor(groupColor)
				nodes[parentCrumb].AddChild(node)
				nodes[crumb] = node
			}
			parentCrumb = crumb
		}
		nodes[id.Group].AddChild(tview.NewTreeNode(id.Name).SetReference(id).SetTextStyle(fileTextStyle))
	}
}

func findFirstFile(node *tview.TreeNode) *snippets.Id {
	if node == nil {
		return nil
	}
	ref := node.GetReference()
	if ref != nil {
		id := ref.(snippets.Id)
		return &id
	}
	children := node.GetChildren()
	for _, child := range children {
		cid := findFirstFile(child)
		if cid != nil {
			return cid
		}
	}
	return nil
}

func init() {
	useCmd.Flags().Bool("allow-dirty", false, "allow to run in dirty git repository")
	viper.BindPFlag(cnfg.KEY_APPLY_ALLOW_DIRTY, useCmd.Flags().Lookup("allow-dirty"))
	useCmd.Flags().BoolVarP(&useCmdCopy, "copy", "c", false, "copy output to clipboard")
	useCmd.Flags().BoolVarP(&useCmdOverwrite, "overwrite", "o", false, "whether to overwrite or append to the file (default false, appends)")
	useCmd.Flags().StringVarP(&useCmdFile, "file", "f", "", "file destination (\".\" to get suggestions)")
	useCmd.MarkFlagsMutuallyExclusive("copy", "file")
	// TODO: require file if overwrite
}
