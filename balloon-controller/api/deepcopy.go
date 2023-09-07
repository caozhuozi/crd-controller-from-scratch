package api

import "k8s.io/apimachinery/pkg/runtime"

func (in *Balloon) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}

	out := new(Balloon)
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)

	return out
}

func (in *BalloonList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}

	out := new(BalloonList)

	*out = *in
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		in, out := &in.Items, &out.Items
		for i := range *in {
			c := (*in)[i].DeepCopyObject().(*Balloon)
			*out = append(*out, *c)
		}
	}

	return out
}
