package run

import (
	runApi "google.golang.org/api/run/v1alpha1"
)

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

type IamPolicy runApi.Policy

func (in *IamPolicy) DeepCopy() *IamPolicy {
	if in == nil {
		return nil
	}
	out := new(IamPolicy)
	in.DeepCopyInto(out)
	return out
}

func (in *IamPolicy) DeepCopyInto(out *IamPolicy) {
	*out = *in
	return
}
