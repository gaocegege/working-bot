package weekly

import (
	"context"
	"fmt"

	"github.com/dyweb/gommon/util/logutil"
	"github.com/google/go-github/github"

	"github.com/gaocegege/working-bot/cli/working-bot/server/config"
	"github.com/gaocegege/working-bot/pkg/gh"
	"github.com/gaocegege/working-bot/pkg/util/weeklyutil"
)

const (
	labelWorking = "working"
)

var log = logutil.NewPackageLogger()

type Worker struct {
	config config.Config
}

func NewWorker(config config.Config) *Worker {
	return &Worker{
		config: config,
	}
}

func (w Worker) HandleWeekly(issue github.Issue) error {
	if err := w.removeWorkingLabel(issue); err != nil {
		return err
	}
	if err := w.OpenNewIssue(issue); err != nil {
		return err
	}
	if err := w.commitAndSubmitPR(issue); err != nil {
		return err
	}
	return nil
}

func (w Worker) removeWorkingLabel(issue github.Issue) error {
	gc := gh.GetGitHubClient()
	ctx := context.Background()

	_, err := gc.Client.Issues.RemoveLabelForIssue(ctx, gc.Owner(), gc.Repo(), *issue.Number, labelWorking)
	if err != nil {
		return err
	}
	return nil
}

func (w Worker) OpenNewIssue(oldIssue github.Issue) error {
	gc := gh.GetGitHubClient()
	ctx := context.Background()

	weeklyNum, err := weeklyutil.GetWeeklyNumber(oldIssue)
	if err != nil {
		return err
	}
	weeklyNum++

	title := fmt.Sprintf("Weekly-%d", weeklyNum)
	assignee := "gaocegege-bot"
	body := fmt.Sprintf("工作周报第 %d 期开始 :tada:", weeklyNum)
	newIssue := &github.IssueRequest{
		Title: &title,
		Labels: &[]string{
			labelWorking,
		},
		Assignee: &assignee,
		Body:     &body,
	}

	_, _, err = gc.Client.Issues.Create(ctx, gc.Owner(), gc.Repo(), newIssue)
	if err != nil {
		return err
	}
	return nil
}

func (w Worker) commitAndSubmitPR(issue github.Issue) error {
	newBranch := generateNewBranch()
	log.Infof("generate a new branch name %s", newBranch)

	// do prepare thing before cli and api doc generation.
	if err := w.prepareGitEnv(newBranch); err != nil {
		log.Errorf("failed to prepare git environment: %v", err)
		return err
	}

	if err := w.buildWeekly(issue); err != nil {
		return err
	}

	// commit and push branch
	if err := w.gitCommitAndPush(newBranch); err != nil {
		if err == ErrNothingChanged {
			// if nothing changed, no need to submit pull request.
			return nil
		}
		return err
	}

	num, err := weeklyutil.GetWeeklyNumber(issue)
	if err != nil {
		return err
	}
	// start to submit pull request
	if err := w.sumbitPR(newBranch, num); err != nil {
		return err
	}
	return nil
}
