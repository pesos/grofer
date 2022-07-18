package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/pesos/grofer/pkg/core"
	"github.com/pesos/grofer/pkg/metrics/factory"
	"github.com/pesos/grofer/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	defaultCid                  = ""
	defaultContainerRefreshRate = 1000
)

// k8sCmd represents the kubernetes command
var k8sCmd = &cobra.Command{
	Use:     "kubernetes",
	Short:   "kubernetes command is used to get information related to the local kubernetes cluster",
	Long:    `kubernetes command is used to get information related to the local kubernetes cluster. It provides both overall and per pod metrics.`,
	Aliases: []string{"k8s", "kubernetes"},
	RunE: func(cmd *cobra.Command, args []string) error {

		// validate args and extract flags.
		containerCmd, err := constructContainerCommand(cmd, args)
		if err != nil {
			return err
		}

		// create a metric scraper factory that will help construct
		// a container metric specific MetricScraper.
		metricScraperFactory := factory.
			NewMetricScraperFactory().
			ForCommand(core.KubernetesCommand).
			WithScrapeInterval(containerCmd.refreshRate)

		/*
		if containerCmd.isPerContainer() {
			metricScraperFactory = metricScraperFactory.ForSingularEntity(containerCmd.cid)
		}
		*/

		// construct a container specific MetricScraper.
		kubeMetricScraper, err := metricScraperFactory.Construct()
		if err != nil {
			return err
		}

		if containerCmd.all {
			err = kubeMetricScraper.Serve(factory.WithAllAs(containerCmd.all))
		} else {
			err = kubeMetricScraper.Serve()
		}

		if err != nil && err != core.ErrCanceledByUser {
			if err == core.ErrInvalidContainer {
				utils.ErrorMsg("cid")
			}
			log.Printf("Error: %v\n", err)
		}

		return nil
	},
}

type k8sCommand struct {
	refreshRate uint64
	cid         string
	all         bool
}

func constructK8sCommand(cmd *cobra.Command, args []string) (*containerCommand, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("the container command should have no arguments, see grofer container --help for further info")
	}
	cid, err := cmd.Flags().GetString("container-id")
	if err != nil {
		return nil, errors.New("error extracting flag --container-id")
	}

	allFlag, err := cmd.Flags().GetBool("all")
	if err != nil {
		return nil, errors.New("error extracting flag --all")
	}

	containerRefreshRate, err := cmd.Flags().GetUint64("refresh")
	if err != nil {
		return nil, errors.New("error extracting flag --refresh")
	}

	if containerRefreshRate < 1000 {
		return nil, errors.New("invalid refresh rate: minimum refresh rate is 1000(ms)")
	}

	containerCmd := &containerCommand{
		refreshRate: containerRefreshRate,
		cid:         cid,
		all:         allFlag,
	}

	return containerCmd, nil
}

func (cc *containerCommand) isPerContainer() bool {
	return cc.cid != defaultCid
}

func init() {
	rootCmd.AddCommand(containerCmd)

	containerCmd.Flags().StringP(
		"container-id",
		"c",
		"",
		"specify container ID",
	)

	containerCmd.Flags().Uint64P(
		"refresh",
		"r",
		defaultContainerRefreshRate,
		"Container information UI refreshes rate in milliseconds greater than 1000",
	)

	containerCmd.Flags().BoolP(
		"all",
		"a",
		false,
		"Specify to list all containers or only running containers.",
	)
}
