package main

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"
)

const commandHelp = `* |/outline search "<your query>" | - Search your teams documents.`

// All struct to easily manipulate outline json
type PayloadDocumentsSearch struct {
	Offset          int    `json:"offset"`
	Limit           int    `json:"limit"`
	Query           string `json:"query"`
	IncludeArchived bool   `json:"includeArchived"`
	IncludeDrafts   bool   `json:"includeDrafts"`
	DateFilter      string `json:"dateFilter"`
}

type DocumentsSearchPagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type DocumentsSearchDataDocument struct {
	Title string `json:"title"`
	UrlId string `json:"urlId"`
	URL   string `json:"url"`
}

type DocumentsSearchData struct {
	Context  string                      `json:"context"`
	Ranking  float64                     `json:"ranking"`
	Document DocumentsSearchDataDocument `json:"document"`
}

type DocumentsSearchSuccess struct {
	Data       []DocumentsSearchData     `json:"data"`
	Pagination DocumentsSearchPagination `json:"pagination"`
}

type DocumentsSearchError struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

// formatDocumentsSearchResult construct the mattermost string response
func formatDocumentsSearchResult(OutlineURL string, result []DocumentsSearchData, query string) string {
	formatResult := []string{}
	formatResult = append(formatResult, fmt.Sprintf("All documents related to `%v`:\n", query))
	for _, v := range result {
		_result := fmt.Sprintf("- [%v](%v%v)\n", v.Document.Title, OutlineURL, v.Document.URL)
		formatResult = append(formatResult, _result)
	}
	return strings.Join(formatResult, " ")
}

// DocumentsSearch allows you to search your teams documents with keywords.
func DocumentsSearch(c *configuration, query string) string {
	client := resty.New()
	//	client.SetTimeout(10)
	client.SetHeader("Accept", "application/json")
	client.SetHeader("Content-Type", "application/json")

	var payload = PayloadDocumentsSearch{
		Offset:          0,
		Limit:           25,
		Query:           query,
		IncludeArchived: true,
		IncludeDrafts:   true,
		DateFilter:      "year"}

	var documentsSearchError = DocumentsSearchError{}
	var documentsSearchSuccess = DocumentsSearchSuccess{}
	resp, err := client.R().
		SetBody(payload).
		SetAuthToken(c.OutlineToken).
		SetResult(&documentsSearchSuccess).
		SetError(&documentsSearchError).
		Post(fmt.Sprintf("%s/api/documents.search", c.OutlineURL))

	if err != nil {
		return fmt.Sprintf("unable to search documents %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Sprintf("%v - %v", resp.StatusCode(), documentsSearchError.Error)
	}

	if len(documentsSearchSuccess.Data) == 0 {
		return fmt.Sprintf("**Sadly no result found :(**\n %v", c.PageNotFoundURL)
	}
	return formatDocumentsSearchResult(c.OutlineURL, documentsSearchSuccess.Data, query)
}

func (p *Plugin) getCommand() (*model.Command, error) {
	iconData, err := command.GetIconData(p.API, "assets/mattermost-outline.svg")

	if err != nil {
		return nil, errors.Wrap(err, "failed to get icon data")
	}

	return &model.Command{
		Trigger:              "outline",
		DisplayName:          "Mattermost Outline",
		Description:          "Mattermost outline plugin allow you to search your teams documents.",
		AutoComplete:         true,
		AutoCompleteDesc:     "Available commands: help, search",
		AutoCompleteHint:     "[command]",
		AutocompleteIconData: iconData,
	}, nil
}

func (p *Plugin) getCommandResponse(args *model.CommandArgs, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: "ephemeral",
		Text:         text,
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	config := p.getConfiguration()

	var (
		split      = strings.Fields(args.Command)
		command    = split[0]
		action     string
		parameters []string
	)
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}
	if command != "/outline" {
		return &model.CommandResponse{}, nil
	}

	if action == "help" || action == "" {
		text := "###### Mattermost Outline Plugin - Slash Command Help\n" + strings.Replace(commandHelp, "|", "`", -1)
		return p.getCommandResponse(args, text), nil
	}

	switch action {
	case "search":
		text := DocumentsSearch(config, strings.Join(parameters, " "))
		return p.getCommandResponse(args, text), nil
	default:
		return p.getCommandResponse(args, "Unknown action, please use `/outline help` to see all actions available."), nil
	}
}
