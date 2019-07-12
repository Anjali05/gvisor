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

#include <arpa/inet.h>
#include <linux/if_tun.h>
#include <net/if.h>
#include <netinet/in.h>
#include <sys/ioctl.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>

#include <cstdio>
#include <cstring>
#include <map>
#include <memory>
#include <unordered_map>
#include <unordered_set>
#include <utility>
#include <vector>

#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "gtest/gtest.h"
#include "test/syscalls/linux/ip_socket_test_util.h"
#include "test/syscalls/linux/socket_test_util.h"
#include "test/syscalls/linux/uid_util.h"
#include "test/util/test_util.h"

namespace gvisor {
namespace testing {

using std::pair;
using std::string;
using std::vector;

std::unordered_set<string> get_interface_names() {
  struct if_nameindex* interfaces = if_nameindex();
  if (interfaces == nullptr) {
    return {};
  }
  std::unordered_set<string> names;
  for (auto interface = interfaces;
       interface->if_index != 0 || interface->if_name != nullptr; interface++) {
    names.insert(interface->if_name);
  }
  if_freenameindex(interfaces);
  return names;
}

class Tunnel {
 public:
  // Creates a tunnel and returns the file descriptor of the tunnel and the
  // name.
  Tunnel() {
    fd_ = open("/dev/net/tun", O_RDWR);
    if (fd_ < 0) {
      return;
    }

    struct ifreq ifr;
    memset(&ifr, 0, sizeof(ifr));
    ifr.ifr_flags = IFF_TUN;

    int err = ioctl(fd_, (TUNSETIFF), (void*)&ifr);
    if (err < 0) {
      close(fd_);
      fd_ = 0;
    }
    name_ = ifr.ifr_name;
  }

  const string& GetName() const { return name_; }

  ~Tunnel() {
    if (fd_ != 0) {
      close(fd_);
    }
  }

 private:
  int fd_;
  string name_;
};

std::unique_ptr<Tunnel> NewTunnel() { return absl::make_unique<Tunnel>(); }

TEST_P(BindToDeviceTest, BindToDevice) {
  // Only root can use SO_BINDTODEVICE.
  SKIP_IF(!ASSERT_NO_ERRNO_AND_VALUE(IsRoot()));
  auto test_case = GetParam();
  auto test_name = test_case.name;
  auto test_actions = test_case.actions;

  auto interface_names = get_interface_names();
  // devices maps from the device id in the test case to the name of the device.
  std::unordered_map<int, string> devices;
  int next_unused_eth = 1;
  std::vector<std::unique_ptr<Tunnel>> tunnels;
  for (const auto& action : test_actions) {
    if (action.device != 0 && devices.find(action.device) == devices.end()) {
      // Need to pick a new device.
      devices[action.device] = string("eth") + std::to_string(next_unused_eth);
      next_unused_eth++;
    }

    if (interface_names.find(devices[action.device]) == interface_names.end()) {
      // Need a tunnel.
      tunnels.push_back(NewTunnel());
      // Skip if failed to get a tunnel.
      SKIP_IF(tunnels.back()->GetName() == "");
      devices[action.device] = tunnels.back()->GetName();
    }
  }

  SCOPED_TRACE(test_name);

  int action_index = 0;
  // sockets_to_close is a map from action index to the socket that was created.
  std::unordered_map<int,
                     std::unique_ptr<gvisor::testing::FileDescriptor>>
      sockets_to_close;
  // All the actions will use the same port, whichever we are assigned.
  in_port_t port = htons(0);
  for (const auto& action : test_actions) {
    SCOPED_TRACE(std::string("Action index: ") + std::to_string(action_index));
    if (action.release) {
      // Close the socket that was made in a previous action.  The release_row
      // indicates which socket to close based on index into the list of
      // actions.
      sockets_to_close.erase(action.release_row);
      continue;
    }

    // Make the socket.
    sockets_to_close[action_index] = ASSERT_NO_ERRNO_AND_VALUE(NewSocket());
    auto socket_fd = sockets_to_close[action_index]->get();
    action_index++;

    // If reuse is indicated, do that.
    if (action.reuse) {
      int reuse = 1;
      EXPECT_THAT(setsockopt(socket_fd, SOL_SOCKET, SO_REUSEPORT, &reuse,
                             sizeof(reuse)),
                  SyscallSucceedsWithValue(0));
    }

    // If the device is non-zero, bind to that device.
    if (action.device != 0) {
      string device_name = devices[action.device];
      EXPECT_THAT(setsockopt(socket_fd, SOL_SOCKET, SO_BINDTODEVICE,
                             device_name.c_str(), device_name.size() + 1),
                  SyscallSucceedsWithValue(0));
      char getDevice[100];
      socklen_t get_device_size = 100;
      EXPECT_THAT(getsockopt(socket_fd, SOL_SOCKET, SO_BINDTODEVICE, getDevice,
                             &get_device_size),
                  SyscallSucceedsWithValue(0));
    }

    struct sockaddr_in addr = {};
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
    addr.sin_port = port;
    if (action.want == 0) {
      ASSERT_THAT(
          bind(socket_fd, reinterpret_cast<const struct sockaddr*>(&addr),
               sizeof(addr)),
          SyscallSucceeds());
    } else {
      ASSERT_THAT(
          bind(socket_fd, reinterpret_cast<const struct sockaddr*>(&addr),
               sizeof(addr)),
          SyscallFailsWithErrno(action.want));
    }

    if (port == 0) {
      // We don't yet know what port we'll be using so we need to fetch it and
      // remember it for future commands.
      socklen_t addr_size = sizeof(addr);
      ASSERT_THAT(
          getsockname(socket_fd, reinterpret_cast<struct sockaddr*>(&addr),
                      &addr_size),
          SyscallSucceeds());
      port = addr.sin_port;
    }
  }
}

}  // namespace testing
}  // namespace gvisor
