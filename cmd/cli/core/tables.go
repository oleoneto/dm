package core

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TableOptions struct {
	ColumnConfig *[]table.ColumnConfig
	Title        *string
	Header       table.Row
	Footer       *table.Row
	Style        *table.Style
	OutputMirror io.Writer
	SortField    string
	SortOrder    table.SortMode
}

type TableRowOptions struct{}

func Initialize(opts TableOptions) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(opts.OutputMirror)

	if opts.Style != nil {
		t.SetStyle(*opts.Style)
	} else {
		t.SetStyle(table.StyleLight)
	}

	if opts.ColumnConfig != nil {
		t.SetColumnConfigs(*opts.ColumnConfig)
	}

	if opts.Title != nil {
		t.SetTitle(*opts.Title)
	}

	t.AppendHeader(opts.Header)

	if opts.Footer != nil {
		t.AppendFooter(*opts.Footer)
	}

	if opts.SortField != "" {
		t.SortBy([]table.SortBy{{Name: opts.SortField, Mode: func() table.SortMode {
			if opts.SortOrder < 0 {
				return opts.SortOrder
			}

			return table.Dsc
		}()}})
	}

	return t
}
