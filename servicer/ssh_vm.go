package servicer

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"
)

type VmWare struct {
	IP 		string
	User 	string
	Pwd		string
	client  *govmomi.Client
	ctx 	context.Context
}

type VirtualMachines struct {
	Name   string
	System string
	Self   Self
	VM     types.ManagedObjectReference
}

type TemplateInfo struct {
	Name   string
	System string
	Self   Self
	VM     types.ManagedObjectReference
}

type DatastoreSummary struct {
	Datastore          Datastore `json:"Datastore"`
	Name               string    `json:"Name"`
	URL                string    `json:"Url"`
	Capacity           int64     `json:"Capacity"`
	FreeSpace          int64     `json:"FreeSpace"`
	Uncommitted        int64     `json:"Uncommitted"`
	Accessible         bool      `json:"Accessible"`
	MultipleHostAccess bool      `json:"MultipleHostAccess"`
	Type               string    `json:"Type"`
	MaintenanceMode    string    `json:"MaintenanceMode"`
	DatastoreSelf      types.ManagedObjectReference
}

type Datastore struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type HostSummary struct {
	Host        Host   `json:"Host"`
	Name        string `json:"Name"`
	UsedCPU     int64  `json:"UsedCPU"`
	TotalCPU    int64  `json:"TotalCPU"`
	FreeCPU     int64  `json:"FreeCPU"`
	UsedMemory  int64  `json:"UsedMemory"`
	TotalMemory int64  `json:"TotalMemory"`
	FreeMemory  int64  `json:"FreeMemory"`
	HostSelf    types.ManagedObjectReference
}

type Host struct {
	Type  string `json:"Type"`
	Value string `json:"Value"`
}

type HostVM struct {
	Host map[string][]VMS
}

type VMS struct {
	Name  string
	Value string
}

type DataCenter struct {
	Datacenter      Self
	Name            string
	VmFolder        Self
	HostFolder      Self
	DatastoreFolder Self
}

type ClusterInfo struct {
	Cluster      Self
	Name         string
	Parent       Self
	ResourcePool Self
	Hosts        []types.ManagedObjectReference
	Datastore    []types.ManagedObjectReference
}

type ResourcePoolInfo struct {
	ResourcePool     Self
	Name             string
	Parent           Self
	ResourcePoolList []types.ManagedObjectReference
	Resource         types.ManagedObjectReference
}

type FolderInfo struct {
	Folder      Self
	Name        string
	ChildEntity []types.ManagedObjectReference
	Parent      Self
	FolderSelf  types.ManagedObjectReference
}

type Self struct {
	Type  string
	Value string
}

type CreateMap struct {
	TempName    string
	Datacenter  string
	Cluster     string
	Host        string
	Resources   string
	Storage     string
	VmName      string
	SysHostName string
	Network     string
}

var vw *VmWare

func NewVmWare(VmUrl,username,password string) *VmWare {
	u := &url.URL{
		Scheme: "https",
		Host:   VmUrl,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(username, password)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Println(err)
	}
	vw = &VmWare{
		IP:     VmUrl,
		User:   username,
		Pwd:    password,
		client: client,
		ctx:    ctx,
	}

	return vw

}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func (vm *VmWare) getBase(tp string) (v *view.ContainerView, error error) {
	fmt.Println("client:",vw.ctx,vw.client)
	m := view.NewManager(vw.client.Client)

	v, err := m.CreateContainerView(vw.ctx, vw.client.Client.ServiceContent.RootFolder, []string{tp}, true)
	if err != nil {
		return nil, err
	}
	return v, nil
}


func (vm *VmWare) GetAllVmClient() (vmList []VirtualMachines, templateList []TemplateInfo, err error) {
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, nil, err
	}



	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, nil, err
	}
	for _, vm := range vms {
		//if vm.Summary.Config.Name == "测试机器" {
		//	v := object.NewVirtualMachine(vw.client.Client, vm.Self)
		//	vw.setIP(v)
		//}
		if vm.Summary.Config.Template {
			templateList = append(templateList, TemplateInfo{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		} else {
			vmList = append(vmList, VirtualMachines{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		}
	}
	fmt.Println(vmList)
	return vmList, templateList, nil
}


func (vm *VmWare) GetAllHost() (hostList []*HostSummary, err error) {
	v, err := vw.getBase("HostSystem")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var hss []mo.HostSystem
	err = v.Retrieve(vw.ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	for _, hs := range hss {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		hostList = append(hostList, &HostSummary{
			Host: Host{
				Type:  hs.Summary.Host.Type,
				Value: hs.Summary.Host.Value,
			},
			Name:        hs.Summary.Config.Name,
			UsedCPU:     int64(hs.Summary.QuickStats.OverallCpuUsage),
			TotalCPU:    totalCPU,
			FreeCPU:     freeCPU,
			UsedMemory:  int64((units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage)) * 1024 * 1024),
			TotalMemory: int64(units.ByteSize(hs.Summary.Hardware.MemorySize)),
			FreeMemory:  freeMemory,
			HostSelf:    hs.Self,
		})
	}
	return hostList, err
}

