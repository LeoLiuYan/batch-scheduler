/*
 * Copyright 2020 THL A29 Limited, a Tencent company.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"encoding/json"
	"fmt"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	v1 "k8s.io/api/core/v1"

	v12 "github.com/tenstack/batch-scheduler/pkg/apis/podgroup/v1"
)

// DefaultWaitTime is 60s if MaxScheduleTime is not specified.
const DefaultWaitTime = 60 * time.Second

// CreateMergePatch return patch generated from original and new interfaces
func CreateMergePatch(original, new interface{}) ([]byte, error) {
	pvByte, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	cloneByte, err := json.Marshal(new)
	if err != nil {
		return nil, err
	}
	patch, err := jsonpatch.CreateMergePatch(pvByte, cloneByte)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

// VerifyPodAnnSatisfied verifies if pod ann satisfies batch scheduling
func VerifyPodAnnSatisfied(pod *v1.Pod) (string, bool) {
	if pod.Annotations == nil {
		return "", false
	}
	if pod.Annotations[PodGroupAnn] == "" {
		return "", false
	}
	return pod.Annotations[PodGroupAnn], true
}

// VerifyPodLabelSatisfied verifies if pod ann satisfies batch scheduling
func VerifyPodLabelSatisfied(pod *v1.Pod) (string, bool) {
	if pod.Labels == nil {
		return "", false
	}
	if pod.Labels[PodGroupLabel] == "" {
		return "", false
	}
	return pod.Labels[PodGroupLabel], true
}

// GetPodGroupFullName verify if pod ann satisfies batch scheduling
func GetPodGroupFullName(pg *v12.PodGroup) string {
	if pg == nil {
		return ""
	}

	return fmt.Sprintf("%v/%v", pg.Namespace, pg.Name)
}

// GetWaitTimeDuration verify if pod ann satisfies batch scheduling
func GetWaitTimeDuration(pg *v12.PodGroup, defaultMaxScheTime *time.Duration) time.Duration {
	waitTime := DefaultWaitTime
	if defaultMaxScheTime != nil || *defaultMaxScheTime != 0 {
		waitTime = *defaultMaxScheTime
	}
	if pg != nil && pg.Spec.MaxScheduleTime != nil {
		return pg.Spec.MaxScheduleTime.Duration
	}
	return waitTime
}
