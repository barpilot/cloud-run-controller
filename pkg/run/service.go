package run

import (
	runApi "google.golang.org/api/run/v1alpha1"
)

// type Service runApi.Service

// type Service struct {
// 	*runApi.Service
// }

type Service runApi.Service

func (in *Service) DeepCopy() *Service {
	if in == nil {
		return nil
	}
	out := new(Service)
	in.DeepCopyInto(out)
	return out
}

func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
	return
}
