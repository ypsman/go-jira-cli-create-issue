package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	jira "github.com/andygrunwald/go-jira"
	"github.com/metakeule/fmtdate"
	yaml "gopkg.in/yaml.v2"
)

type tickets struct {
	Jiraserver string
	Apiuser    string
	Apipass    string
	Project    string
	Type       string
	Assignee   string
}

var ticket tickets

// Flag for command
var flagSumm = flag.String("sum", "none", "Summary of the Ticket")
var flagDesc = flag.String("desc", "none", "Ticket Desription")
var flagDate = flag.String("date", "", "DueDate DD-MM-YYYY")

func main() {
	flag.Parse()
	_, err := fmtdate.Parse("YYYY-MM-DD", *flagDate)
	if err != nil {
		fmt.Printf("Please check the Date argument.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *flagSumm == "none" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *flagDesc == "none" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	loadConfig()
	makeIssue()
}

func makeIssue() {
	jiraClient, _ := jira.NewClient(nil, ticket.Jiraserver)
	res, err := jiraClient.Authentication.AcquireSessionCookie(ticket.Apiuser, ticket.Apipass)
	if err != nil || res == false {
		fmt.Printf("Result: %v\n", res)
		panic(err)
	}
	i := jira.Issue{

		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: ticket.Assignee,
			},
			Reporter: &jira.User{
				Name: ticket.Apiuser,
			},
			Description: *flagDesc,
			Duedate:     *flagDate,
			Type: jira.IssueType{
				Name: ticket.Type,
			},
			Project: jira.Project{
				Key: ticket.Project,
			},
			Summary: *flagSumm,
		},
	}

	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		panic(err)
	}
	fmt.Println(issue.Key)
}

func loadConfig() {
	_, filename, _, _ := runtime.Caller(1)
	configfile := path.Join(path.Dir(filename), "jira-ci.yml")
	source, err := ioutil.ReadFile(configfile)
	check(err)
	err = yaml.Unmarshal(source, &ticket)
	check(err)
	if ticket.Jiraserver == "" {
		fmt.Printf("Please set the Jiraserver in config.\n")
		os.Exit(1)
	}
	if ticket.Apiuser == "" {
		fmt.Printf("Please set the Apiuser in config.\n")
		os.Exit(1)
	}
	if ticket.Apipass == "" {
		fmt.Printf("Please set the Apipass in config.\n")
		os.Exit(1)
	}
	if ticket.Project == "" {
		fmt.Printf("Please set the Project in config.\n")
		os.Exit(1)
	}
	if ticket.Type == "" {
		fmt.Printf("Please set the Type in config.\n")
		os.Exit(1)
	}
	if ticket.Assignee == "" {
		fmt.Printf("Please set the Assignee in config.\n")
		os.Exit(1)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
