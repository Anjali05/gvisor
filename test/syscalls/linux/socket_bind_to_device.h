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

#ifndef GVISOR_TEST_SYSCALLS_LINUX_SOCKET_BIND_TO_DEVICE_H_
#define GVISOR_TEST_SYSCALLS_LINUX_SOCKET_BIND_TO_DEVICE_H_

#include <string>

#include "test/syscalls/linux/ip_socket_test_util.h"
#include "test/syscalls/linux/socket_test_util.h"

namespace gvisor {
namespace testing {

typedef struct TestAction {
  bool reuse;
  int device;
  bool release;
  int release_row;
  int want;
} TestAction;

typedef struct TestCase {
  std::string name;
  std::vector<TestAction> actions;
} TestCase;

// Test fixture for SO_BINDTODEVICE tests.
class BindToDeviceTest : public ::testing::TestWithParam<struct TestCase> {
 protected:
  BindToDeviceTest() {
    // gUnit uses printf, so so will we.
    printf("Testing case: %s\n", GetParam().name.c_str());
  }

  PosixErrorOr<std::unique_ptr<FileDescriptor>> NewSocket() const {
    return IPv4UDPUnboundSocket(0).Create();
  }
};

}  // namespace testing
}  // namespace gvisor

#endif  // GVISOR_TEST_SYSCALLS_LINUX_SOCKET_BIND_TO_DEVICE_H_
