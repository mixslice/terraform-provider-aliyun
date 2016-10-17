package aliyun

import (
	"fmt"
	"log"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunEcsInstance() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        resourceAliyunEcsInstanceCreate,
		Read:          resourceAliyunEcsInstanceRead,
		Update:        resourceAliyunEcsInstanceUpdate,
		Delete:        resourceAliyunEcsInstanceDelete,
		Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// The methods defined below will get called for each resource that needs to
// get created (createFunc), read (readFunc), updated (updateFunc) and deleted (deleteFunc).
// For example, if 10 resources need to be created then `createFunc`
// will get called 10 times every time with the information for the proper
// resource that is being mapped.
//
// If at some point any of these functions returns an error, Terraform will
// imply that something went wrong with the modification of the resource and it
// will prevent the execution of further calls that depend on that resource
// that failed to be created/updated/deleted.

func resourceAliyunEcsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).ecsclient

	args := ecs.CreateInstanceArgs{
		RegionId:     common.Region(d.Get("region").(string)),
		ImageId:      d.Get("image").(string),
		InstanceType: d.Get("instance_type").(string),
		InstanceName: d.Get("name").(string),
	}

	instanceID, err := client.CreateInstance(&args)
	if err != nil {
		return fmt.Errorf("Failed to create instance from Image %s: %v", args.ImageId, err)
	}

	d.SetId(instanceID)

	// Update if we need to
	return resourceAliyunEcsInstanceUpdate(d, meta)
}

func resourceAliyunEcsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAliyunEcsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAliyunEcsInstanceRead(d, meta)
}

func resourceAliyunEcsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).ecsclient

	if err := aliyunTerminateInstance(client, d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// InstanceStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an ecs instance.
func InstanceStateRefreshFunc(client *ecs.Client, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := client.DescribeInstanceAttribute(instanceID)
		if err != nil {
			if reqErr, _ := err.(*common.Error); reqErr.StatusCode == 404 {
				return nil, "", nil
			}

			log.Printf("Failed to describe Instance %s attribute: %v\n", instanceID, err)
			return nil, "", err
		}

		return &instance, string(instance.Status), nil
	}
}

func aliyunTerminateInstance(client *ecs.Client, id string) error {
	log.Printf("[INFO] Terminating instance: %s", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating", "Running", "Starting", "Stopping"},
		Target:     []string{"Stopped"},
		Refresh:    InstanceStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to terminate: %s", id, err)
	}

	if err := client.DeleteInstance(id); err != nil {
		return fmt.Errorf("Error terminating instance: %s", err)
	}

	log.Printf("[DEBUG] Waiting for instance (%s) to become terminated", id)

	return nil
}
