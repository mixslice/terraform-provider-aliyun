package aliyun

import "github.com/hashicorp/terraform/helper/schema"

func resourceAliyunInstance() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        createFunc,
		Read:          readFunc,
		Update:        updateFunc,
		Delete:        deleteFunc,
		Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram": &schema.Schema{
				Type:     schema.TypeInt,
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

func createFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config)
	machine := Machine{
		Name: d.Get("name").(string),
		CPUs: d.Get("cpus").(int),
		RAM:  d.Get("ram").(int),
	}

	err := client.CreateMachine(&machine)
	if err != nil {
		return err
	}

	d.SetId(machine.Id())

	return nil
}

func readFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func updateFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func deleteFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}
