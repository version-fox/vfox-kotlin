// Copyright 2026 Han Li and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package common holds the target-independent helpers of the Kotlin E2E
// suites. The linux and windows packages drive a real vfox in throwaway
// containers, one case per container; common decides which versions to test,
// how many cases run at once, how to name containers and how to retry
// network-touching steps.
//
// Run the offline unit tests with `go test -short ./...`; the container suites
// are skipped in short mode and run with `go test -run TestE2E ./linux` or
// `go test -run TestE2E ./windows`. See docs/e2e.md.
package common
