# IBM Cloud VPC with s390x Instance Example

This example demonstrates how to use the Pulumi IBM Cloud provider to create a VPC and launch an s390x (IBM LinuxONE) virtual server instance.

## Prerequisites

1. **IBM Cloud Account**: You need an active IBM Cloud account
2. **IBM Cloud API Key**: Create an API key from the IBM Cloud console
3. **Pulumi CLI**: Install the Pulumi CLI from [pulumi.com](https://www.pulumi.com/docs/get-started/install/)
4. **Go**: This example uses Go, so you need Go 1.24 or later installed

## What This Example Creates

- **VPC**: A Virtual Private Cloud for network isolation
- **Subnet**: A subnet within the VPC for the instance
- **SSH Key**: An SSH key resource for secure instance access
- **s390x Instance**: A virtual server instance using s390x architecture (IBM LinuxONE)

## Configuration

Before running this example, you need to set the following configuration values:

### Required Configuration

```bash
# Your SSH public key for accessing the instance
pulumi config set sshPublicKey "ssh-rsa AAAAB3NzaC1yc2EA..."

# s390x compatible image ID (you can find this in IBM Cloud console or CLI)
# Example for Ubuntu on s390x: r014-xxxx-xxxx-xxxx-xxxx
pulumi config set imageId "r014-xxxx-xxxx-xxxx-xxxx"

# IBM Cloud API key (this will be stored as a secret)
pulumi config set ibmcloud:apiKey --secret
```

### Optional Configuration

```bash
# IBM Cloud region (default: us-east)
pulumi config set region us-east

# Availability zone (default: {region}-1)
pulumi config set zone us-east-1

# Resource group (optional)
pulumi config set resourceGroup "Default"
```

## Finding a s390x Image

To find available s390x images in IBM Cloud:

```bash
# Using IBM Cloud CLI
ibmcloud is images --visibility public | grep s390x

# Or via API
curl -X GET "https://us-east.iaas.cloud.ibm.com/v1/images?version=2024-11-19&generation=2" \
  -H "Authorization: Bearer $IAM_TOKEN" | jq '.images[] | select(.operating_system.architecture == "s390x")'
```

Common s390x images include:
- Ubuntu for s390x
- Red Hat Enterprise Linux for s390x
- SUSE Linux Enterprise Server for s390x

## s390x Instance Profiles

This example uses the `bz2e-1x4` profile, which provides:
- Architecture: s390x (IBM LinuxONE)
- 1 vCPU
- 4 GB RAM

Other available s390x profiles include:
- `bz2e-1x4`: 1 vCPU, 4 GB RAM
- `bz2e-2x8`: 2 vCPU, 8 GB RAM
- `bz2e-4x16`: 4 vCPU, 16 GB RAM
- `bz2e-8x32`: 8 vCPU, 32 GB RAM

## Running the Example

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Initialize Pulumi stack**:
   ```bash
   pulumi stack init dev
   ```

3. **Configure the stack** (see Configuration section above)

4. **Preview the changes**:
   ```bash
   pulumi preview
   ```

5. **Deploy the resources**:
   ```bash
   pulumi up
   ```

6. **View the outputs**:
   ```bash
   pulumi stack output
   ```

## Outputs

After deployment, the following outputs are available:

- `vpcId`: The ID of the created VPC
- `vpcName`: The name of the VPC
- `subnetId`: The ID of the created subnet
- `subnetName`: The name of the subnet
- `instanceId`: The ID of the s390x instance
- `instanceName`: The name of the instance
- `instanceStatus`: The current status of the instance
- `primaryNetworkInterface`: Network interface details including IP address

## Accessing the Instance

Once the instance is running, you can SSH into it using the private key corresponding to your configured public key:

```bash
# Get the instance's primary IP
INSTANCE_IP=$(pulumi stack output primaryNetworkInterface | jq -r '.primary_ip.address')

# SSH into the instance
ssh -i ~/.ssh/your_private_key root@$INSTANCE_IP
```

**Note**: You may need to configure security group rules to allow SSH access (port 22) if the default security group doesn't permit it.

## Cleanup

To destroy all resources created by this example:

```bash
pulumi destroy
```

## Cost Considerations

Running this example will incur costs in your IBM Cloud account. s390x instances on IBM LinuxONE may have different pricing than other instance types. Make sure to:

1. Review IBM Cloud pricing for VPC and virtual server instances
2. Destroy resources when you're done experimenting to avoid ongoing charges
3. Monitor your IBM Cloud billing dashboard

## Troubleshooting

### Image Not Found

If you get an error about the image not being found:
- Ensure the image ID is valid and available in your selected region
- Verify the image supports s390x architecture
- Check that you have permission to use the image

### SSH Key Issues

If you can't connect via SSH:
- Verify your SSH public key format is correct
- Check security group rules allow inbound SSH (port 22)
- Ensure you're using the correct private key
- Wait a few minutes for the instance to fully initialize

### Region/Zone Availability

Not all IBM Cloud regions support s390x instances. Currently supported regions include:
- us-east (Washington DC)
- us-south (Dallas)
- eu-gb (London)
- eu-de (Frankfurt)

## Learn More

- [IBM Cloud VPC Documentation](https://cloud.ibm.com/docs/vpc)
- [IBM LinuxONE Virtual Servers](https://cloud.ibm.com/docs/vpc?topic=vpc-about-virtual-server-instances)
- [Pulumi IBM Cloud Provider](https://github.com/mapt-oss/pulumi-ibmcloud)
- [Pulumi Documentation](https://www.pulumi.com/docs/)
