/*
Copyright 2019 JD Cloud

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

// Code generated by injection-gen. DO NOT EDIT.

package client

import (
	"context"

	injection "github.com/kubernetes-local-volume/kubernetes-local-volume/pkg/common/injection"
	logging "github.com/kubernetes-local-volume/kubernetes-local-volume/pkg/common/logging"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

func init() {
	injection.Default.RegisterClient(withClient)
}

// Key is used as the key for associating information with a context.Context.
type Key struct{}

func withClient(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, Key{}, kubernetes.NewForConfigOrDie(cfg))
}

// Get extracts the kubernetes.Interface client from the context.
func Get(ctx context.Context) kubernetes.Interface {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch k8s.io/client-go/kubernetes.Interface from context.")
	}
	return untyped.(kubernetes.Interface)
}
