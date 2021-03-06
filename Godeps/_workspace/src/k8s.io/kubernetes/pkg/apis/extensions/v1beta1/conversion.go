/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package v1beta1

import (
	"fmt"
	"reflect"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/conversion"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/intstr"
)

func addConversionFuncs(scheme *runtime.Scheme) {
	// Add non-generated conversion functions
	err := scheme.AddConversionFuncs(
		Convert_api_PodSpec_To_v1_PodSpec,
		Convert_v1_PodSpec_To_api_PodSpec,
		Convert_extensions_DeploymentSpec_To_v1beta1_DeploymentSpec,
		Convert_v1beta1_DeploymentSpec_To_extensions_DeploymentSpec,
		Convert_extensions_DeploymentStrategy_To_v1beta1_DeploymentStrategy,
		Convert_v1beta1_DeploymentStrategy_To_extensions_DeploymentStrategy,
		Convert_extensions_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment,
		Convert_v1beta1_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment,
		Convert_extensions_DaemonSetSpec_To_v1beta1_DaemonSetSpec,
		Convert_v1beta1_DaemonSetSpec_To_extensions_DaemonSetSpec,
		Convert_extensions_DaemonSetUpdateStrategy_To_v1beta1_DaemonSetUpdateStrategy,
		Convert_v1beta1_DaemonSetUpdateStrategy_To_extensions_DaemonSetUpdateStrategy,
		Convert_extensions_RollingUpdateDaemonSet_To_v1beta1_RollingUpdateDaemonSet,
		Convert_v1beta1_RollingUpdateDaemonSet_To_extensions_RollingUpdateDaemonSet,
		Convert_extensions_ReplicaSetSpec_To_v1beta1_ReplicaSetSpec,
		Convert_v1beta1_ReplicaSetSpec_To_extensions_ReplicaSetSpec,

		Convert_api_VolumeSource_To_v1beta1_VolumeSource,
		Convert_v1beta1_VolumeSource_To_api_VolumeSource,
	)
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}

	// Add field label conversions for kinds having selectable nothing but ObjectMeta fields.
	for _, kind := range []string{"DaemonSet", "Deployment", "Ingress"} {
		err = api.Scheme.AddFieldLabelConversionFunc("extensions/v1beta1", kind,
			func(label, value string) (string, string, error) {
				switch label {
				case "metadata.name", "metadata.namespace":
					return label, value, nil
				default:
					return "", "", fmt.Errorf("field label %q not supported for %q", label, kind)
				}
			})
		if err != nil {
			// If one of the conversion functions is malformed, detect it immediately.
			panic(err)
		}
	}

	err = api.Scheme.AddFieldLabelConversionFunc("extensions/v1beta1", "Job",
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name", "metadata.namespace", "status.successful":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		})
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}

// The following two PodSpec conversions functions where copied from pkg/api/conversion.go
// for the generated functions to work properly.
// This should be fixed: https://github.com/kubernetes/kubernetes/issues/12977
func Convert_api_PodSpec_To_v1_PodSpec(in *api.PodSpec, out *v1.PodSpec, s conversion.Scope) error {
	return v1.Convert_api_PodSpec_To_v1_PodSpec(in, out, s)
}

func Convert_v1_PodSpec_To_api_PodSpec(in *v1.PodSpec, out *api.PodSpec, s conversion.Scope) error {
	return v1.Convert_v1_PodSpec_To_api_PodSpec(in, out, s)
}

func Convert_extensions_DeploymentSpec_To_v1beta1_DeploymentSpec(in *extensions.DeploymentSpec, out *DeploymentSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.DeploymentSpec))(in)
	}
	out.Replicas = new(int32)
	*out.Replicas = int32(in.Replicas)
	if in.Selector != nil {
		out.Selector = new(LabelSelector)
		if err := Convert_unversioned_LabelSelector_To_v1beta1_LabelSelector(in.Selector, out.Selector, s); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := v1.Convert_api_PodTemplateSpec_To_v1_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_extensions_DeploymentStrategy_To_v1beta1_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	if in.RevisionHistoryLimit != nil {
		out.RevisionHistoryLimit = new(int32)
		*out.RevisionHistoryLimit = int32(*in.RevisionHistoryLimit)
	}
	out.MinReadySeconds = int32(in.MinReadySeconds)
	out.Paused = in.Paused
	if in.RollbackTo != nil {
		out.RollbackTo = new(RollbackConfig)
		out.RollbackTo.Revision = int64(in.RollbackTo.Revision)
	} else {
		out.RollbackTo = nil
	}
	return nil
}

func Convert_v1beta1_DeploymentSpec_To_extensions_DeploymentSpec(in *DeploymentSpec, out *extensions.DeploymentSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*DeploymentSpec))(in)
	}
	if in.Replicas != nil {
		out.Replicas = int(*in.Replicas)
	}

	if in.Selector != nil {
		out.Selector = new(unversioned.LabelSelector)
		if err := Convert_v1beta1_LabelSelector_To_unversioned_LabelSelector(in.Selector, out.Selector, s); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := v1.Convert_v1_PodTemplateSpec_To_api_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_DeploymentStrategy_To_extensions_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	if in.RevisionHistoryLimit != nil {
		out.RevisionHistoryLimit = new(int)
		*out.RevisionHistoryLimit = int(*in.RevisionHistoryLimit)
	}
	out.MinReadySeconds = int(in.MinReadySeconds)
	out.Paused = in.Paused
	if in.RollbackTo != nil {
		out.RollbackTo = new(extensions.RollbackConfig)
		out.RollbackTo.Revision = in.RollbackTo.Revision
	} else {
		out.RollbackTo = nil
	}
	return nil
}

func Convert_extensions_DeploymentStrategy_To_v1beta1_DeploymentStrategy(in *extensions.DeploymentStrategy, out *DeploymentStrategy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.DeploymentStrategy))(in)
	}
	out.Type = DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(RollingUpdateDeployment)
		if err := Convert_extensions_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_v1beta1_DeploymentStrategy_To_extensions_DeploymentStrategy(in *DeploymentStrategy, out *extensions.DeploymentStrategy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*DeploymentStrategy))(in)
	}
	out.Type = extensions.DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(extensions.RollingUpdateDeployment)
		if err := Convert_v1beta1_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_extensions_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(in *extensions.RollingUpdateDeployment, out *RollingUpdateDeployment, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.RollingUpdateDeployment))(in)
	}
	if out.MaxUnavailable == nil {
		out.MaxUnavailable = &intstr.IntOrString{}
	}
	if err := s.Convert(&in.MaxUnavailable, out.MaxUnavailable, 0); err != nil {
		return err
	}
	if out.MaxSurge == nil {
		out.MaxSurge = &intstr.IntOrString{}
	}
	if err := s.Convert(&in.MaxSurge, out.MaxSurge, 0); err != nil {
		return err
	}
	return nil
}

func Convert_v1beta1_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment(in *RollingUpdateDeployment, out *extensions.RollingUpdateDeployment, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RollingUpdateDeployment))(in)
	}
	if err := s.Convert(in.MaxUnavailable, &out.MaxUnavailable, 0); err != nil {
		return err
	}
	if err := s.Convert(in.MaxSurge, &out.MaxSurge, 0); err != nil {
		return err
	}
	return nil
}

func Convert_extensions_DaemonSetSpec_To_v1beta1_DaemonSetSpec(in *extensions.DaemonSetSpec, out *DaemonSetSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.DaemonSetSpec))(in)
	}
	// unable to generate simple pointer conversion for unversioned.LabelSelector -> v1beta1.LabelSelector
	if in.Selector != nil {
		out.Selector = new(LabelSelector)
		if err := Convert_unversioned_LabelSelector_To_v1beta1_LabelSelector(in.Selector, out.Selector, s); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := v1.Convert_api_PodTemplateSpec_To_v1_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_extensions_DaemonSetUpdateStrategy_To_v1beta1_DaemonSetUpdateStrategy(&in.UpdateStrategy, &out.UpdateStrategy, s); err != nil {
		return err
	}
	out.UniqueLabelKey = new(string)
	*out.UniqueLabelKey = in.UniqueLabelKey
	return nil
}

func Convert_v1beta1_DaemonSetSpec_To_extensions_DaemonSetSpec(in *DaemonSetSpec, out *extensions.DaemonSetSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*DaemonSetSpec))(in)
	}
	// unable to generate simple pointer conversion for v1beta1.LabelSelector -> unversioned.LabelSelector
	if in.Selector != nil {
		out.Selector = new(unversioned.LabelSelector)
		if err := Convert_v1beta1_LabelSelector_To_unversioned_LabelSelector(in.Selector, out.Selector, s); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := v1.Convert_v1_PodTemplateSpec_To_api_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_DaemonSetUpdateStrategy_To_extensions_DaemonSetUpdateStrategy(&in.UpdateStrategy, &out.UpdateStrategy, s); err != nil {
		return err
	}
	if in.UniqueLabelKey != nil {
		out.UniqueLabelKey = *in.UniqueLabelKey
	}
	return nil
}

func Convert_extensions_DaemonSetUpdateStrategy_To_v1beta1_DaemonSetUpdateStrategy(in *extensions.DaemonSetUpdateStrategy, out *DaemonSetUpdateStrategy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.DaemonSetUpdateStrategy))(in)
	}
	out.Type = DaemonSetUpdateStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(RollingUpdateDaemonSet)
		if err := Convert_extensions_RollingUpdateDaemonSet_To_v1beta1_RollingUpdateDaemonSet(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_v1beta1_DaemonSetUpdateStrategy_To_extensions_DaemonSetUpdateStrategy(in *DaemonSetUpdateStrategy, out *extensions.DaemonSetUpdateStrategy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*DaemonSetUpdateStrategy))(in)
	}
	out.Type = extensions.DaemonSetUpdateStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(extensions.RollingUpdateDaemonSet)
		if err := Convert_v1beta1_RollingUpdateDaemonSet_To_extensions_RollingUpdateDaemonSet(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_extensions_RollingUpdateDaemonSet_To_v1beta1_RollingUpdateDaemonSet(in *extensions.RollingUpdateDaemonSet, out *RollingUpdateDaemonSet, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.RollingUpdateDaemonSet))(in)
	}
	if out.MaxUnavailable == nil {
		out.MaxUnavailable = &intstr.IntOrString{}
	}
	if err := s.Convert(&in.MaxUnavailable, out.MaxUnavailable, 0); err != nil {
		return err
	}
	out.MinReadySeconds = int32(in.MinReadySeconds)
	return nil
}

func Convert_v1beta1_RollingUpdateDaemonSet_To_extensions_RollingUpdateDaemonSet(in *RollingUpdateDaemonSet, out *extensions.RollingUpdateDaemonSet, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RollingUpdateDaemonSet))(in)
	}
	if err := s.Convert(in.MaxUnavailable, &out.MaxUnavailable, 0); err != nil {
		return err
	}
	out.MinReadySeconds = int(in.MinReadySeconds)
	return nil
}

func Convert_extensions_ReplicaSetSpec_To_v1beta1_ReplicaSetSpec(in *extensions.ReplicaSetSpec, out *ReplicaSetSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*extensions.ReplicaSetSpec))(in)
	}
	out.Replicas = new(int32)
	*out.Replicas = int32(in.Replicas)
	if in.Selector != nil {
		out.Selector = new(LabelSelector)
		if err := Convert_unversioned_LabelSelector_To_v1beta1_LabelSelector(in.Selector, out.Selector, s); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if in.Template != nil {
		out.Template = new(v1.PodTemplateSpec)
		if err := v1.Convert_api_PodTemplateSpec_To_v1_PodTemplateSpec(in.Template, out.Template, s); err != nil {
			return err
		}
	} else {
		out.Template = nil
	}
	return nil
}

func Convert_v1beta1_ReplicaSetSpec_To_extensions_ReplicaSetSpec(in *ReplicaSetSpec, out *extensions.ReplicaSetSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*ReplicaSetSpec))(in)
	}
	if in.Replicas != nil {
		out.Replicas = int(*in.Replicas)
	}
	if in.Selector != nil {
		out.Selector = new(unversioned.LabelSelector)
		if err := Convert_v1beta1_LabelSelector_To_unversioned_LabelSelector(in.Selector, out.Selector, s); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if in.Template != nil {
		out.Template = new(api.PodTemplateSpec)
		if err := v1.Convert_v1_PodTemplateSpec_To_api_PodTemplateSpec(in.Template, out.Template, s); err != nil {
			return err
		}
	} else {
		out.Template = nil
	}
	return nil
}

///////////////////////
///////// CARRIES
///////////////////////

// This will Convert our internal represantation of VolumeSource to its v1 representation
// Used for keeping backwards compatibility for the Metadata field
func Convert_api_VolumeSource_To_v1beta1_VolumeSource(in *api.VolumeSource, out *v1.VolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.VolumeSource))(in)
	}

	if err := s.DefaultConvert(in, out, conversion.IgnoreMissingFields); err != nil {
		return err
	}

	if in.DownwardAPI != nil {
		out.DownwardAPI = new(v1.DownwardAPIVolumeSource)
		if err := Convert_api_DownwardAPIVolumeSource_To_v1_DownwardAPIVolumeSource(in.DownwardAPI, out.DownwardAPI, s); err != nil {
			return err
		}

		// also copy to Metadata
		out.Metadata = new(v1.MetadataVolumeSource)
		if err := Convert_api_DownwardAPIVolumeSource_To_v1_MetadataVolumeSource(in.DownwardAPI, out.Metadata, s); err != nil {
			return err
		}
	} else {
		out.DownwardAPI = nil
	}
	return nil
}

// downward -> metadata (api -> v1)
func Convert_api_DownwardAPIVolumeSource_To_v1_MetadataVolumeSource(in *api.DownwardAPIVolumeSource, out *v1.MetadataVolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.DownwardAPIVolumeSource))(in)
	}
	if in.Items != nil {
		out.Items = make([]v1.MetadataFile, len(in.Items))
		for i := range in.Items {
			if err := Convert_api_DownwardAPIVolumeFile_To_v1_MetadataFile(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	}
	return nil
}

func Convert_api_DownwardAPIVolumeFile_To_v1_MetadataFile(in *api.DownwardAPIVolumeFile, out *v1.MetadataFile, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.DownwardAPIVolumeFile))(in)
	}
	out.Name = in.Path
	if err := Convert_api_ObjectFieldSelector_To_v1_ObjectFieldSelector(&in.FieldRef, &out.FieldRef, s); err != nil {
		return err
	}
	return nil
}

// This will Convert the v1 representation of VolumeSource to our internal representation
// Used for keeping backwards compatibility for the Metadata field
func Convert_v1beta1_VolumeSource_To_api_VolumeSource(in *v1.VolumeSource, out *api.VolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.VolumeSource))(in)
	}

	if err := s.DefaultConvert(in, out, conversion.IgnoreMissingFields); err != nil {
		return err
	}

	// if specified Metadata will stomp DownwardAPI
	if in.Metadata != nil {
		out.DownwardAPI = new(api.DownwardAPIVolumeSource)
		if err := Convert_v1_MetadataVolumeSource_To_api_DownwardAPIVolumeSource(in.Metadata, out.DownwardAPI, s); err != nil {
			return err
		}
	} else {
		if in.DownwardAPI != nil {
			out.DownwardAPI = new(api.DownwardAPIVolumeSource)
			if err := Convert_v1_DownwardAPIVolumeSource_To_api_DownwardAPIVolumeSource(in.DownwardAPI, out.DownwardAPI, s); err != nil {
				return err
			}
		} else {
			out.DownwardAPI = nil
		}
	}
	return nil
}

// metadata -> downward (v1 -> api)
func Convert_v1_MetadataVolumeSource_To_api_DownwardAPIVolumeSource(in *v1.MetadataVolumeSource, out *api.DownwardAPIVolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.MetadataVolumeSource))(in)
	}
	if in.Items != nil {
		out.Items = make([]api.DownwardAPIVolumeFile, len(in.Items))
		for i := range in.Items {
			if err := Convert_v1_MetadataFile_To_api_DownwardAPIVolumeFile(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	}
	return nil
}

func Convert_v1_MetadataFile_To_api_DownwardAPIVolumeFile(in *v1.MetadataFile, out *api.DownwardAPIVolumeFile, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.MetadataFile))(in)
	}
	out.Path = in.Name
	if err := Convert_v1_ObjectFieldSelector_To_api_ObjectFieldSelector(&in.FieldRef, &out.FieldRef, s); err != nil {
		return err
	}

	return nil
}
