/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package messages

const (
	// PubSub
	ErrPubsubEmpty              = "pubsub name is empty"
	ErrPubsubNotFound           = "pubsub %s not found"
	ErrTopicEmpty               = "topic is empty in pubsub %s"
	ErrPubsubCloudEventsSer     = "error when marshalling cloud event envelope for topic %s pubsub %s: %s"
	ErrPubsubPublishMessage     = "error when publish to topic %s in pubsub %s: %s"
	ErrPubsubCloudEventCreation = "cannot create cloudevent: %s"
	// Http.
	ErrNotFound             = "method %q is not found"
	ErrMalformedRequest     = "failed deserializing HTTP body: %s"
	ErrMalformedRequestData = "can't serialize request data field: %s"
	// State
	ErrStateStoresNotConfigured = "state store is not configured"
	ErrStateStoreNotFound       = "state store %s is not found"
	ErrStateGet                 = "fail to get %s from state store %s: %s"
	ErrStateDelete              = "failed deleting state with key %s: %s"
	ErrStateSave                = "failed saving state in state store %s: %s"
	ErrStateQuery               = "failed query in state store %s: %s"
	// StateTransaction
	ErrStateStoreNotSupported     = "state store %s doesn't support transaction"
	ErrNotSupportedStateOperation = "operation type %s not supported"
	ErrStateTransaction           = "error while executing state transaction: %s"
	//	Lock
	ErrLockStoresNotConfigured = "lock store is not configured"
	ErrResourceIdEmpty         = "ResourceId is empty in lock store %s"
	ErrLockOwnerEmpty          = "LockOwner is empty in lock store %s"
	ErrExpireNotPositive       = "Expire is not positive in lock store %s"
	ErrLockStoreNotFound       = "lock store %s not found"
	//	Sequencer
	ErrSequencerStoresNotConfigured = "Sequencer store is not configured"
	ErrSequencerKeyEmpty            = "Key is empty in sequencer store %s"
	ErrSequencerStoreNotFound       = "Sequencer store %s not found"

	// Binding.
	ErrInvokeOutputBinding = "error when invoke output binding %s: %s"

	// Secret
	ErrSecretStoreNotConfigured = "error when get secret but not find configured"
	ErrSecretStoreNotFound      = "error when get secret but not find : %s"
	ErrSecretGet                = "error when get secret : secret name => %s,store name =>%s,error => %s"
	ErrBulkSecretGet            = "error when bulk get secret %s: %s"
	ErrPermissionDenied         = "access denied by policy to get %s from %s"
)
