package main

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

const (
	V6 = `https://root:awwtgl01@@10.31.39.214/sdk`
	V4 = `https://10.31.39.180/sdk`
)

func main() {
	u, _ := url.Parse(V6)
	c, err := govmomi.NewClient(context.Background(), u, true)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Is VCentre:", c.IsVC())

	vimClient := c.Client
	printVMs(vimClient)

}

func printVMs(c *vim25.Client) {
	m := view.NewManager(c)
	ctx := context.Background()

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		panic(err)
	}

	defer v.Destroy(ctx)
	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		panic(err)
	}

	// Print summary per vm (see also: govc/vm/info.go)

	//data, _ := json.MarshalIndent(vms, "", "  ")
	//fmt.Println(string(data))
	for _, vm := range vms {
		printVM(vm)
	}
}

func printVM(vm mo.VirtualMachine) {
	keys := map[string]string{
		"Type":  vm.Summary.Config.GuestFullName,
		"State": string(vm.Summary.Runtime.PowerState),
		// "MAC Addr:"
	}


	fmt.Println(vm.Summary.Config.Name)
	for k, v := range keys {
		fmt.Printf("\t%s\t: %s\n", k, v)
	}

	poweredOn := vm.Summary.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn
	if poweredOn && vm.Summary.Guest != nil {
		printExtendedData(vm)
	}


	//getMacAddr(vm)
	//fmt.Printf("%s\n", vm.Summary.Config.Name)
	//fmt.Printf("\tType:\t%s\n", vm.Summary.Config.GuestFullName)
	//
	//fmt.Printf("%s: %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
}

func printExtendedData(vm mo.VirtualMachine) {
	guest := vm.Summary.Guest
	keys := map[string]string{
		"Tools Installed": guest.ToolsRunningStatus,
		"Hostname": guest.HostName,
		"IP": guest.IpAddress,
		"Guest OS": guest.GuestFullName,
	}
	for k, v := range keys {
		fmt.Printf("\t%s\t: %s\n", k, v)
	}
}

func getMacAddr(vm mo.VirtualMachine) {
	//hw := vm.Config.Hardware.Device

	fmt.Println("\tNetwork Props:")
	for _, n := range vm.Network {
		fmt.Printf("\t\t%s: %s\n", n.Type, n.Value)
	}

	fmt.Println("\tDevices:")

	// vm.Network.
	// for _, d := range hw {
	// 	dev := d.GetVirtualDevice()
	// 	// dev.
	// }
}
