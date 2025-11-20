package main

import (
	"github.com/mapt-oss/pulumi-ibmcloud/sdk/go/ibmcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		// Get configuration values with defaults
		region := cfg.Get("region")
		if region == "" {
			region = "us-east" // Default region
		}

		zone := cfg.Get("zone")
		if zone == "" {
			zone = region + "-1" // Default zone
		}

		resourceGroup := cfg.Get("resourceGroup")

		// SSH public key - user should provide this via config
		sshPublicKey := cfg.Require("sshPublicKey")

		// Create a VPC
		vpc, err := ibmcloud.NewIsVpc(ctx, "s390x-vpc", &ibmcloud.IsVpcArgs{
			Name:          pulumi.String("s390x-example-vpc"),
			ResourceGroup: pulumi.String(resourceGroup),
			Tags: pulumi.StringArray{
				pulumi.String("example"),
				pulumi.String("s390x"),
				pulumi.String("vpc"),
			},
		})
		if err != nil {
			return err
		}

		// Create a subnet in the VPC
		subnet, err := ibmcloud.NewIsSubnet(ctx, "s390x-subnet", &ibmcloud.IsSubnetArgs{
			Name:                 pulumi.String("s390x-example-subnet"),
			Vpc:                  vpc.ID(),
			Zone:                 pulumi.String(zone),
			TotalIpv4AddressCount: pulumi.Int(256),
			ResourceGroup:        pulumi.String(resourceGroup),
		})
		if err != nil {
			return err
		}

		// Create SSH key for instance access
		sshKey, err := ibmcloud.NewIsSshKey(ctx, "s390x-ssh-key", &ibmcloud.IsSshKeyArgs{
			Name:          pulumi.String("s390x-example-key"),
			PublicKey:     pulumi.String(sshPublicKey),
			ResourceGroup: pulumi.String(resourceGroup),
			Tags: pulumi.StringArray{
				pulumi.String("example"),
				pulumi.String("s390x"),
			},
		})
		if err != nil {
			return err
		}

		// Create a s390x virtual server instance
		// Using bz2e profile which is for s390x architecture (IBM LinuxONE)
		instance, err := ibmcloud.NewIsInstance(ctx, "s390x-instance", &ibmcloud.IsInstanceArgs{
			Name:          pulumi.String("s390x-example-instance"),
			Profile:       pulumi.String("bz2e-1x4"), // s390x profile: 1 vCPU, 4GB RAM
			Image:         pulumi.String(cfg.Get("imageId")), // User must provide a valid s390x image ID
			Zone:          pulumi.String(zone),
			ResourceGroup: pulumi.String(resourceGroup),
			Keys: pulumi.StringArray{
				sshKey.ID(),
			},
			PrimaryNetworkInterface: &ibmcloud.IsInstancePrimaryNetworkInterfaceArgs{
				Subnet: subnet.ID(),
				Name:   pulumi.String("eth0"),
			},
			VpcId: vpc.ID(),
			Tags: pulumi.StringArray{
				pulumi.String("example"),
				pulumi.String("s390x"),
				pulumi.String("linuxone"),
			},
		})
		if err != nil {
			return err
		}

		// Export important values
		ctx.Export("vpcId", vpc.ID())
		ctx.Export("vpcName", vpc.Name)
		ctx.Export("subnetId", subnet.ID())
		ctx.Export("subnetName", subnet.Name)
		ctx.Export("instanceId", instance.ID())
		ctx.Export("instanceName", instance.Name)
		ctx.Export("instanceStatus", instance.Status)
		ctx.Export("primaryNetworkInterface", instance.PrimaryNetworkInterface)

		return nil
	})
}
