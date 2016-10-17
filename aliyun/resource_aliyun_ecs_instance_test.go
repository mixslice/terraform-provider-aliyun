package aliyun

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAliyunInstance_basic(t *testing.T) {
	var v ecs.InstanceAttributesType

	testCheck := func(*terraform.State) error {
		// if *v.Placement.AvailabilityZone != "us-west-2a" {
		// 	return fmt.Errorf("bad availability zone: %#v", *v.Placement.AvailabilityZone)
		// }

		// if len(v.SecurityGroups) == 0 {
		// 	return fmt.Errorf("no security groups: %#v", v.SecurityGroups)
		// }
		// if *v.SecurityGroups[0].GroupName != "tf_test_foo" {
		// 	return fmt.Errorf("no security groups: %#v", v.SecurityGroups)
		// }

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		// We ignore security groups because even with EC2 classic
		// we'll import as VPC security groups, which is fine. We verify
		// VPC security group import in other tests
		IDRefreshName:   "aliyun_ecs_instance.foo",
		IDRefreshIgnore: []string{"security_groups", "vpc_security_group_ids"},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"aliyun_ecs_instance.foo", &v),
					testCheck,
					resource.TestCheckResourceAttr(
						"aliyun_ecs_instance.foo",
						"image",
						"ubuntu1404_64_40G_cloudinit_20160727.raw"),
					resource.TestCheckResourceAttr(
						"aliyun_ecs_instance.foo", "instance_type", "ecs.s1.small"),
				),
			},
		},
	})
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*AliyunClient).ecsclient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aliyun_ecs_instance" {
			continue
		}

		// Try to find the resource
		instance, err := client.DescribeInstanceAttribute(rs.Primary.ID)
		if instance != nil {
			return fmt.Errorf("Found unterminated instance: %s", instance)
		}

		// Verify the error is what we want
		if reqErr, ok := err.(*common.Error); ok && reqErr.StatusCode == 404 {
			continue
		}

		return err
	}

	return nil
}

func testAccCheckInstanceExists(n string, i *ecs.InstanceAttributesType) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckInstanceExistsWithProviders(n, i, &providers)
}

func testAccCheckInstanceExistsWithProviders(n string, i *ecs.InstanceAttributesType, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*AliyunClient).ecsclient
			instance, err := client.DescribeInstanceAttribute(rs.Primary.ID)

			if reqErr, ok := err.(*common.Error); ok && reqErr.StatusCode == 404 {
				continue
			}
			if err != nil {
				return err
			}

			if instance != nil {
				*i = *instance
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

const testAccInstanceConfig = `
resource "aliyun_ecs_instance" "foo" {
	image = "ubuntu1404_64_40G_cloudinit_20160727.raw"
	instance_type = "ecs.s1.small"
}
`
