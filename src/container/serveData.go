/*
Copyright © 2020 The PES Open Source Team pesos@pes.edu

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
package container

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pesos/grofer/src/utils"
)

// Serve serves overall container metrics
func Serve(ctx context.Context, cli *client.Client, all bool, dataChannel chan ContainerMetrics, refreshRate int64) error {
	return utils.TickUntilDone(ctx, refreshRate, func() error {
		metrics, err := GetOverallMetrics(ctx, cli, all)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- metrics:
		}

		return nil
	})
}

// ServeContainer serves data on a per container basis
func ServeContainer(ctx context.Context, cli *client.Client, cid string, dataChannel chan PerContainerMetrics, refreshRate int64) error {
	return utils.TickUntilDone(ctx, refreshRate, func() error {
		metrics, err := GetContainerMetrics(ctx, cli, cid)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case dataChannel <- metrics:
		}

		return nil
	})
}
