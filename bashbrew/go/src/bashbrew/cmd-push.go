package main

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
)

func cmdPush(c *cli.Context) error {
	repos, err := repos(c.Bool("all"), c.Args()...)
	if err != nil {
		return cli.NewMultiError(fmt.Errorf(`failed gathering repo list`), err)
	}

	uniq := c.Bool("uniq")
	namespace := c.String("namespace")

	if namespace == "" {
		return fmt.Errorf(`"--namespace" is a required flag for "push"`)
	}

	for _, repo := range repos {
		r, err := fetch(repo)
		if err != nil {
			return cli.NewMultiError(fmt.Errorf(`failed fetching repo %q`, repo), err)
		}

		for _, entry := range r.Entries() {
			if r.SkipConstraints(entry) {
				continue
			}

			for _, tag := range r.Tags(namespace, uniq, entry) {
				created := dockerCreated(tag)
				lastUpdated := fetchDockerHubTagMeta(tag).lastUpdatedTime()
				if created.After(lastUpdated) {
					fmt.Printf("Pushing %s\n", tag)
					err = dockerPush(tag)
					if err != nil {
						return cli.NewMultiError(fmt.Errorf(`failed pushing %q`, tag), err)
					}
				} else {
					fmt.Printf("Skipping %s (created %s, last updated %s)\n", tag, created.Local().Format(time.RFC3339), lastUpdated.Local().Format(time.RFC3339))
				}
			}
		}
	}

	return nil
}
