// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "test/syscalls/linux/socket_bind_to_device.h"

#include <vector>

#include "test/syscalls/linux/ip_socket_test_util.h"
#include "test/syscalls/linux/socket_test_util.h"
#include "test/util/test_util.h"

namespace gvisor {
namespace testing {

std::vector<TestCase> GetTestCases() {
  return std::vector<TestCase>{
      {.name = "bind twice with device fails",
       .actions =
           {
               {.device = 3, .reuse = 0, .want = 0},
               {.device = 3, .reuse = 0, .want = EADDRINUSE},
           }},
      {.name = "bind to device",
       .actions =
           {
               {.device = 1, .want = 0},
               {.device = 2, .want = 0},
           }},
      {.name = "bind to device and then without device",
       .actions =
           {
               {.device = 123, .want = 0},
               {.device = 0, .want = EADDRINUSE},
           }},
      {.name = "bind without device",
       .actions =
           {
               {.want = 0},
               {.device = 123, .want = EADDRINUSE},
               {.device = 123, .reuse = true, .want = EADDRINUSE},
               {.want = EADDRINUSE},
               {.reuse = true, .want = EADDRINUSE},
           }},
      {.name = "bind with device",
       .actions =
           {
               {.device = 123, .want = 0},
               {.device = 123, .want = EADDRINUSE},
               {.device = 123, .reuse = true, .want = EADDRINUSE},
               {.device = 0, .want = EADDRINUSE},
               {.device = 0, .reuse = true, .want = EADDRINUSE},
               {.device = 456, .reuse = true, .want = 0},
               {.device = 789, .want = 0},
               {.want = EADDRINUSE},
               {.reuse = true, .want = EADDRINUSE},
           }},
      {.name = "bind with reuse",
       .actions =
           {
               {.reuse = true, .want = 0},
               {.device = 123, .want = EADDRINUSE},
               {.device = 123, .reuse = true, .want = 0},
               {.device = 0, .want = EADDRINUSE},
               {.device = 0, .reuse = true, .want = 0},
           }},
      {.name = "binding with reuse and device",
       .actions =
           {
               {.device = 123, .reuse = true, .want = 0},
               {.device = 123, .want = EADDRINUSE},
               {.device = 123, .reuse = true, .want = 0},
               {.device = 0, .want = EADDRINUSE},
               {.device = 456, .reuse = true, .want = 0},
               {.device = 0, .reuse = true, .want = 0},
               {.device = 789, .reuse = true, .want = 0},
               {.device = 999, .want = EADDRINUSE},
           }},
      {.name = "mixing reuse and not reuse by binding to device",
       .actions =
           {
               {.device = 123, .reuse = true, .want = 0},
               {.device = 456, .want = 0},
               {.device = 789, .reuse = true, .want = 0},
               {.device = 999, .want = 0},
           }},
      {.name = "can't bind to 0 after mixing reuse and not reuse",
       .actions =
           {
               {.device = 123, .reuse = true, .want = 0},
               {.device = 456, .want = 0},
               {.device = 0, .reuse = true, .want = EADDRINUSE},
           }},
      {.name = "bind and release",
       .actions =
           {
               {.device = 123, .reuse = true, .want = 0},
               {.device = 0, .reuse = true, .want = 0},
               {.device = 345, .reuse = false, .want = EADDRINUSE},
               {.device = 789, .reuse = true, .want = 0},

               // Release the bind to device 0 and try again.
               {.release = true, .release_row = 1},
               {.device = 345, .reuse = false, .want = 0},
           }},
      {.name = "bind twice with reuse once",
       .actions =
           {
               {.device = 123, .reuse = false, .want = 0},
               {.device = 0, .reuse = true, .want = EADDRINUSE},
           }},
  };
}

INSTANTIATE_TEST_SUITE_P(BindToDeviceSockets, BindToDeviceTest,
                         ::testing::ValuesIn(GetTestCases()));
}  // namespace testing
}  // namespace gvisor
