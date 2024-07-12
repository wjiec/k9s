package dao

import (
	"context"
	"github.com/derailed/k9s/internal/client"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/scale"
)

var _ Scalable = (*Scaler)(nil)

// Scaler represents a generic resource.
type Scaler struct {
	Generic
}

func (s *Scaler) Scale(ctx context.Context, path string, replicas int32) error {
	ns, n := client.Namespaced(path)

	cfg, err := s.Client().RestConfig()
	if err != nil {
		return err
	}

	discoveryClient, err := s.Client().CachedDiscovery()
	if err != nil {
		return err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	scaleKindResolver := scale.NewDiscoveryScaleKindResolver(discoveryClient)

	scaleClient, err := scale.NewForConfig(cfg, mapper, dynamic.LegacyAPIPathResolverFunc, scaleKindResolver)
	if err != nil {
		return err
	}

	currentScale, err := scaleClient.Scales(ns).Get(ctx, *s.gvr.GR(), n, metav1.GetOptions{})
	if err != nil {
		return err
	}

	currentScale.Spec.Replicas = replicas
	updatedScale, err := scaleClient.Scales(ns).Update(ctx, *s.gvr.GR(), currentScale, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.Debug().Msgf("%s scaled to %d", path, updatedScale.Spec.Replicas)
	return nil
}
