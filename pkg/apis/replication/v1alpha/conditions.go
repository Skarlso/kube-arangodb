//
// DISCLAIMER
//
// Copyright 2018 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package v1alpha

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is a strongly typed condition name
type ConditionType string

const (
	// ConditionTypeConfigured indicates that the replication has been configured.
	ConditionTypeConfigured ConditionType = "Configured"
)

// Condition represents one current condition of a deployment or deployment member.
// A condition might not show up if it is not happening.
// For example, if a cluster is not upgrading, the Upgrading condition would not show up.
type Condition struct {
	// Type of  condition.
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// ConditionList is a list of conditions.
// Each type is allowed only once.
type ConditionList []Condition

// IsTrue return true when a condition with given type exists and its status is `True`.
func (list ConditionList) IsTrue(conditionType ConditionType) bool {
	c, found := list.Get(conditionType)
	return found && c.Status == v1.ConditionTrue
}

// Get a condition by type.
// Returns true if found, false if not found.
func (list ConditionList) Get(conditionType ConditionType) (Condition, bool) {
	for _, x := range list {
		if x.Type == conditionType {
			return x, true
		}
	}
	// Not found
	return Condition{}, false
}

// Update the condition, replacing an old condition with same type (if any)
// Returns true when changes were made, false otherwise.
func (list *ConditionList) Update(conditionType ConditionType, status bool, reason, message string) bool {
	src := *list
	statusX := v1.ConditionFalse
	if status {
		statusX = v1.ConditionTrue
	}
	for i, x := range src {
		if x.Type == conditionType {
			if x.Status != statusX {
				// Transition to another status
				src[i].Status = statusX
				now := metav1.Now()
				src[i].LastTransitionTime = now
				src[i].LastUpdateTime = now
				src[i].Reason = reason
				src[i].Message = message
			} else if x.Reason != reason || x.Message != message {
				src[i].LastUpdateTime = metav1.Now()
				src[i].Reason = reason
				src[i].Message = message
			} else {
				return false
			}
			return true
		}
	}
	// Not found
	now := metav1.Now()
	*list = append(src, Condition{
		Type:               conditionType,
		LastUpdateTime:     now,
		LastTransitionTime: now,
		Status:             statusX,
		Reason:             reason,
		Message:            message,
	})
	return true
}

// Remove the condition with given type.
// Returns true if removed, or false if not found.
func (list *ConditionList) Remove(conditionType ConditionType) bool {
	src := *list
	for i, x := range src {
		if x.Type == conditionType {
			*list = append(src[:i], src[i+1:]...)
			return true
		}
	}
	// Not found
	return false
}
