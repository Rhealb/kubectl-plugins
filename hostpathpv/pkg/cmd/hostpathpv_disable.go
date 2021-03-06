/*Copyright 2017 The Kubernetes Authors.

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

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/Rhealb/csi-plugin/hostpathpv/pkg/hostpath/xfsquotamanager"
	xfs "github.com/Rhealb/csi-plugin/hostpathpv/pkg/hostpath/xfsquotamanager/common"

	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/kubectl/util/templates"
)

var (
	disable_valid_resources = `Valid resource types include:

    * node
    `

	disable_long = templates.LongDesc(`
		Disable or unDisable node quota disk.

		` + disable_valid_resources)

	disable_example = templates.Examples(`
		# Disable node nodename quota disk.
		kubectl hostpathpv setdisable node nodename --disk diskpath --disable=true
		
		# Undisable node nodename quota disk.
		kubectl hostpathpv setdisable node nodename --disk diskpath --disable=false
		`)
)

func NewCmdHostPathPVDisable(client *kubernetes.Clientset, out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "setdisable (TYPE/NAME ...) [flags]",
		Short:   T("Disable or unDisable node quota disk"),
		Long:    disable_long,
		Example: disable_example,
		Run: func(cmd *cobra.Command, args []string) {
			err := RunDisable(client, out, errOut, cmd, args)
			if err != nil {
				fmt.Fprintf(out, "run RunDisable err:%v\n", err)
			}
		},
	}
	cmd.Flags().String("disk", "", "Disable or unDisable quota disk of node")
	cmd.Flags().Bool("disable", true, "Set disk diskable or undiskable")
	return cmd
}

func RunDisable(clientset *kubernetes.Clientset, out, errOut io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Fprint(errOut, "You must specify the type of resource to disable. ", disable_valid_resources)
		usageString := "Required resource not specified."
		return UsageErrorf(cmd, usageString)
	}

	resource := args[0]
	diskpath := GetFlagString(cmd, "disk")
	disable := GetFlagBool(cmd, "disable")

	if diskpath == "" {
		return fmt.Errorf("please input diskpath")
	}
	switch {
	case resource == "nodes" || resource == "node":
		nodeName := ""
		if len(args) >= 2 {
			nodeName = args[1]
		}
		err := disableNode(clientset, nodeName, diskpath, disable)
		if err != nil {
			fmt.Fprintf(errOut, "disable err: %v\n", err)
		} else {
			fmt.Fprintf(out, "set success\n")
		}
	default:
		fmt.Fprint(errOut, "You must specify the type of resource to disable. ", disable_valid_resources)
		usageString := "Required resource not suport."
		return UsageErrorf(cmd, usageString)
	}

	return nil
}

func disableNode(clientset *kubernetes.Clientset, nodeName, diskpath string, disable bool) error {
	node, errGetNode := clientset.Core().Nodes().Get(nodeName, metav1.GetOptions{})
	if errGetNode != nil {
		return fmt.Errorf("get node err:%v", errGetNode)
	}
	list := getNodeQuotaDisks(node)
	for _, disk := range list {
		if path.Clean(disk.MountPath) == path.Clean(diskpath) {
			if disk.Disabled == disable {
				return nil
			} else {
				changed := addOrRemoveNodeDisableDisk(node, diskpath, disable)
				if changed == false {
					return nil
				}
				_, errUpdate := clientset.Core().Nodes().Update(node)
				return errUpdate
			}
		}
	}
	return fmt.Errorf("%s is not in quota disks", diskpath)
}

func addOrRemoveNodeDisableDisk(node *v1.Node, diskpath string, add bool) (changed bool) {
	if node == nil {
		return false
	}
	if node.Annotations == nil {
		node.Annotations = make(map[string]string)
	}
	paths := getNodeQuotadiskDisableList(node)

	for i, addedPath := range paths {
		if path.Clean(addedPath) == path.Clean(diskpath) {
			if add == true {
				return false
			} else {
				paths = append(paths[0:i], paths[i+1:]...)
				changed = true
				break
			}
		}
	}
	if add {
		paths = append(paths, diskpath)
		changed = true
	}
	node.Annotations[xfs.NodeDiskQuotaDisableListAnn] = strings.Join(paths, ",")
	return changed
}
func getNodeQuotadiskDisableList(node *v1.Node) []string {
	ret := make([]string, 0, 10)
	if node == nil || node.Annotations == nil || node.Annotations[xfs.NodeDiskQuotaDisableListAnn] == "" {
		return ret
	}
	list := node.Annotations[xfs.NodeDiskQuotaDisableListAnn]
	strs := strings.Split(list, ",")
	for _, str := range strs {
		str = strings.Trim(str, " ")
		if str != "" {
			ret = append(ret, str)
		}
	}
	return ret
}
func getNodeQuotaDisks(node *v1.Node) xfsquotamanager.NodeDiskQuotaInfoList {
	list := make(xfsquotamanager.NodeDiskQuotaInfoList, 0)
	if node == nil || node.Annotations == nil || node.Annotations[xfs.NodeDiskQuotaInfoAnn] == "" {
		return list
	}
	diskinfo := node.Annotations[xfs.NodeDiskQuotaInfoAnn]

	errUmarshal := json.Unmarshal([]byte(diskinfo), &list)
	if errUmarshal != nil {
		return list
	}
	return list
}
