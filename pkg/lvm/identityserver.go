/*
Copyright 2017 The Kubernetes Authors.

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

package lvm

import (
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
)

type identityServer struct {
	*csicommon.DefaultIdentityServer
}

const PluginName = "csi-lvmplugin"

func (is *identityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	fmt.Println("!!!! buka in OUR get plugin info")
	resp := &csi.GetPluginInfoResponse{
		Name: PluginName,
		//VendorVersion: "notImportant",
	}

	// d.log.WithFields(logrus.Fields{
	// 	"response": resp,
	// 	"method":   "get_plugin_info",
	// }).Info("get plugin info called")
	return resp, nil
}
