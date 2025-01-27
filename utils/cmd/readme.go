/*
Copyright (C) 2022 The Falco Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/falcosecurity/evolution/utils/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	readmeRepoFilePath    string
	readmeOutFilePath     string
	readmeTextStartTagFmt = "<!-- REPOSITORY-%s-TABLE -->\n"
	readmeTextEndTagFmt   = "<!-- /REPOSITORY-%s-TABLE -->\n"
)

func readmeTextEditor(s string, status utils.RepositoryStatus) (string, error) {
	startTag := fmt.Sprintf(readmeTextStartTagFmt, strings.ToUpper(status.String()))
	endTag := fmt.Sprintf(readmeTextEndTagFmt, strings.ToUpper(status.String()))
	if len(s) == 0 {
		s = startTag + endTag
	}
	repos, err := utils.ReadRepositoriesFromFile(readmeRepoFilePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	empty := true
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Name", "Description"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetRowSeparator("-")
	table.SetAutoWrapText(false)
	for _, r := range repos {
		if r.Status == status {
			row := []string{}
			row = append(row, fmt.Sprintf("[falcosecurity/%s](https://github.com/falcosecurity/%s)", r.Name, r.Name))
			row = append(row, r.Description)
			table.Append(row)
			empty = false
		}
	}
	if !empty {
		table.Render()
	}
	return utils.ReplaceTextTags(s, startTag, endTag, buf.String())
}

func readmeTextEditorOfficial(s string) (string, error) {
	return readmeTextEditor(s, utils.RepositoryStatusOfficial)
}

func readmeTextEditorIncubating(s string) (string, error) {
	return readmeTextEditor(s, utils.RepositoryStatusIncubating)
}

func readmeTextEditorSandbox(s string) (string, error) {
	return readmeTextEditor(s, utils.RepositoryStatusSandbox)
}

func readmeTextEditorSpecial(s string) (string, error) {
	return readmeTextEditor(s, utils.RepositoryStatusSpecial)
}

var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Generate README.md for falcosecurity/evolution",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(readmeRepoFilePath) == 0 {
			return fmt.Errorf("must specify a path to repositories.yaml")
		}
		if len(readmeOutFilePath) == 0 {
			return fmt.Errorf("must specify an output markdown file")
		}
		return utils.EditCreateTextFile(
			readmeOutFilePath,
			readmeTextEditorOfficial,
			readmeTextEditorIncubating,
			readmeTextEditorSandbox,
			readmeTextEditorSpecial,
		)
	},
}

func init() {
	readmeCmd.Flags().StringVarP(&readmeRepoFilePath, "repositories", "r", "", "Path to a repositories.yaml file")
	readmeCmd.Flags().StringVarP(&readmeOutFilePath, "output", "o", "", "Path to an output markdown file")
}
