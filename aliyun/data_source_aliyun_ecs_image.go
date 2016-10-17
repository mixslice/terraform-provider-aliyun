package aliyun

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAliyunEcsImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliyunEcsImageRead,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"most_recent": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"owner_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

// dataSourceAliyunEcsImageDescriptionRead performs the AMI lookup.
func dataSourceAliyunEcsImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).ecsclient

	nameRegex, nameRegexOk := d.GetOk("name_regex")
	ownerAlias, ownerAliasOk := d.GetOk("owner_alias")

	if nameRegexOk == false && ownerAliasOk == false {
		return fmt.Errorf("One of name_regex, or owner_alias must be assigned")
	}

	params := &ecs.DescribeImagesArgs{
		RegionId: common.Region(d.Get("region").(string)),
	}
	if ownerAliasOk {
		params.ImageOwnerAlias = ecs.ImageOwnerAlias(ownerAlias.(string))
	}

	images, _, err := client.DescribeImages(params)
	log.Printf("miaow: %v", images)
	if err != nil {
		return err
	}

	var filteredImages []ecs.ImageType
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, image := range images {
			// Check for a very rare case where the response would include no
			// image name. No name means nothing to attempt a match against,
			// therefore we are skipping such image.
			if image.ImageName == "" {
				log.Printf("[WARN] Unable to find name to match against "+
					"for image ID %q, nothing to do.",
					image.ImageId)
				continue
			}
			if r.MatchString(image.ImageName) {
				filteredImages = append(filteredImages, image)
			}
		}
	} else {
		filteredImages = images[:]
	}

	var image ecs.ImageType
	if len(filteredImages) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(filteredImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] aliyun_ecs_image - multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(filteredImages)
		} else {
			return fmt.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		// Query returned single result.
		image = filteredImages[0]
	}

	log.Printf("[DEBUG] aliyun_ecs_image - Single Image found: %s", image.ImageId)
	return imageDescriptionAttributes(d, image)
}

type imageSort []ecs.ImageType

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime := time.Time(a[i].CreationTime)
	jtime := time.Time(a[j].CreationTime)
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent image out of a slice of images.
func mostRecentImage(images []ecs.ImageType) ecs.ImageType {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

// populate the numerous fields that the image description returns.
func imageDescriptionAttributes(d *schema.ResourceData, image ecs.ImageType) error {
	// Simple attributes first
	d.SetId(image.ImageId)
	d.Set("architecture", image.Architecture)
	d.Set("creation_time", image.CreationTime)
	if image.Description != "" {
		d.Set("description", image.Description)
	}
	d.Set("image_id", image.ImageId)
	if image.ImageOwnerAlias != "" {
		d.Set("image_owner_alias", image.ImageOwnerAlias)
	}
	d.Set("image_name", image.ImageName)
	if image.OSName != "" {
		d.Set("os_name", image.OSName)
	}
	d.Set("status", image.Status)
	return nil
}
