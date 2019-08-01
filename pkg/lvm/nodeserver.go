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
	"log"
	"os"
	"path/filepath"

	//"github.com/docker/docker/pkg/mount"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/kubernetes/pkg/util/mount"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/wavezhang/k8s-csi-lvm/pkg/lvmd"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
	client kubernetes.Interface
	nodeID string
	vgName string
}

func (ns *nodeServer) GetNodeID() string {
	return ns.nodeID
}

//todo: revisit!!!
const maxVolumesPerNode = 10

func (ns *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	fmt.Println("!!!! buka in NodeGetInfo cp1")
	glog.V(3).Infof("method node_get_info is called")
	return &csi.NodeGetInfoResponse{
		NodeId:            ns.nodeID,
		MaxVolumesPerNode: maxVolumesPerNode,

		// make sure that the driver works on this particular region only
		// AccessibleTopology: &csi.Topology{
		// 	Segments: map[string]string{
		// 		"region": ns.region,
		// 	},
		// },
	}, nil
}

// // GetPluginInfo returns metadata of the plugin
// func (ns *nodeServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
// 	fmt.Println("!!!! buka in get plugin info")
// 	resp := &csi.GetPluginInfoResponse{
// 		Name: "csi-lvmplugin",
// 		//VendorVersion: "notImportant",
// 	}

// 	// d.log.WithFields(logrus.Fields{
// 	// 	"response": resp,
// 	// 	"method":   "get_plugin_info",
// 	// }).Info("get plugin info called")
// 	return resp, nil
// }

// GetPluginCapabilities returns available capabilities of the plugin
func (ns *nodeServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	resp := &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			//{
			// Type: &csi.PluginCapability_Service_{
			// 	Service: &csi.PluginCapability_Service{
			// 		Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
			// 	},
			// },
			//},
		},
	}

	fmt.Println("!!! buka get plugin capabilities called")

	// d.log.WithFields(logrus.Fields{
	// 	"response": resp,
	// 	"method":   "get_plugin_capabilities",
	// }).Info("get plugin capabitilies called")
	return resp, nil
}

// Probe returns the health and readiness of the plugin
func (ns *nodeServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	//d.log.WithField("method", "probe").Info("probe called")
	fmt.Println("!!! buka method Probe called")

	//d.readyMu.Lock()
	//defer d.readyMu.Unlock()

	return &csi.ProbeResponse{
		//Ready: &wrappers.BoolValue{
		//	Value: d.ready,
		//},
	}, nil
}

// GetPluginCapabilities returns available capabilities of the plugin
// func (d *Driver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
// 	resp := &csi.GetPluginCapabilitiesResponse{
// 		Capabilities: []*csi.PluginCapability{
// 			{
// 				Type: &csi.PluginCapability_Service_{
// 					Service: &csi.PluginCapability_Service{
// 						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
// 					},
// 				},
// 			},
// 			{
// 				Type: &csi.PluginCapability_Service_{
// 					Service: &csi.PluginCapability_Service{
// 						Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	d.log.WithFields(logrus.Fields{
// 		"response": resp,
// 		"method":   "get_plugin_capabilities",
// 	}).Info("get plugin capabitilies called")
// 	return resp, nil
// }

// // Probe returns the health and readiness of the plugin
// func (d *Driver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
// 	d.log.WithField("method", "probe").Info("probe called")
// 	d.readyMu.Lock()
// 	defer d.readyMu.Unlock()

// 	return &csi.ProbeResponse{
// 		Ready: &wrappers.BoolValue{
// 			Value: d.ready,
// 		},
// 	}, nil
// }

// GetPluginInfo returns metadata of the plugin
// func (d *Driver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
// 	resp := &csi.GetPluginInfoResponse{
// 		Name:          DriverName,
// 		VendorVersion: version,
// 	}

// 	d.log.WithFields(logrus.Fields{
// 		"response": resp,
// 		"method":   "get_plugin_info",
// 	}).Info("get plugin info called")
// 	return resp, nil
//}

// NodeGetInfo returns the supported capabilities of the node server. This
// should eventually return the droplet ID if possible. This is used so the CO
// knows where to place the workload. The result of this function will be used
// by the CO in ControllerPublishVolume.
// func (d *Driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
// 	d.log.WithField("method", "node_get_info").Info("node get info called")
// 	return &csi.NodeGetInfoResponse{
// 		NodeId:            d.nodeId,
// 		MaxVolumesPerNode: maxVolumesPerNode,

// 		// make sure that the driver works on this particular region only
// 		AccessibleTopology: &csi.Topology{
// 			Segments: map[string]string{
// 				"region": d.region,
// 			},
// 		},
// 	}, nil
// }

func (ns *nodeServer) createVolume(ctx context.Context, volumeId string) (*v1.PersistentVolume, error) {
	fmt.Println("!!!!!! buka in createVolume !!!!!!!!!!!!!!!!!!!")

	pv, err := getPV(ns.client, volumeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Failed to get pv by volumeId %s: %s", volumeId, err))
	}
	node, err := getNode(ns.client, ns.GetNodeID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to get node by nodeId %s: %s", ns.GetNodeID(), err))
	}

	nodeAffinityAnn, err := generateNodeAffinity(node)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to generate node affinity annotations for %v: %v", node.GetName(), err))
	}
	cap := pv.Spec.Capacity[v1.ResourceStorage]
	size := cap.Value()

	addr, err := getLVMDAddr(ns.client, ns.GetNodeID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to getLVMDAddr for %v: %v", node, err))
	}

	conn, err := lvmd.NewLVMConnection(addr, connectTimeout)
	defer conn.Close()
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to connect to %v: %v", addr, err))
	}

	resp, err := conn.CreateLV(ctx, &lvmd.LVMOptions{
		VolumeGroup: ns.vgName,
		Name:        volumeId,
		Size:        uint64(size),
	})
	glog.V(3).Infof("CreateLV: %v", resp)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Error in CreateLogicalVolume: err=%v",
			err)
	}

	pv.Spec.NodeAffinity = nodeAffinityAnn
	pv.Annotations[lvmNodeAnnKey] = node.GetName()
	return updatePV(ns.client, pv)
}
func (ns *nodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, fmt.Sprintf("NodeExpandVolume is not yet implemented"))
}

func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	// fsType := req.GetVolumeCapability().GetMount().GetFsType()
	// devicePath := req.GetPublishInfo()["DevicePath"]

	volumeId := req.GetVolumeId()
	devicePath := filepath.Join("/dev/", ns.vgName, volumeId)

	if _, err := os.Stat(devicePath); os.IsNotExist(err) {
		_, err := ns.createVolume(ctx, volumeId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(targetPath, 0750); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	log.Printf("Determining filesystem type at %v", devicePath)
	existingFstype, err := determineFilesystemType(devicePath)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot determine filesystem type: err=%v",
			err)
	}
	log.Printf("Existing filesystem type is '%v'", existingFstype)
	if existingFstype == "" {
		// There is no existing filesystem on the
		// device, format it with the requested
		// filesystem.
		log.Printf("The device %v has no existing filesystem, formatting with %v", devicePath, existingFstype)
		if err := formatDevice(devicePath, defaultFs); err != nil {
			return nil, status.Errorf(
				codes.Internal,
				"formatDevice failed: err=%v",
				err)
		}
		existingFstype = defaultFs
	}

	// Volume Mount
	if notMnt {
		// Get Options
		var options []string
		if req.GetReadonly() {
			options = append(options, "ro")
		} else {
			options = append(options, "rw")
		}
		mountFlags := req.GetVolumeCapability().GetMount().GetMountFlags()
		options = append(options, mountFlags...)

		// Mount
		mounter := mount.New("")
		err = mounter.Mount(devicePath, targetPath, defaultFs, options)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, status.Error(codes.NotFound, "Targetpath not found")
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if notMnt {
		return nil, status.Error(codes.NotFound, "Volume not mounted")
	}

	//TODO: double check!!!
	err = mount.New("").Unmount(targetPath)
	//vf err = util.UnmountPath(req.GetTargetPath(), mount.New(""))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return &csi.NodeStageVolumeResponse{}, nil
}
