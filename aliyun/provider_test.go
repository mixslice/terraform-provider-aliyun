package aliyun

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"aliyun": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALIYUN_ACCESS_KEY_ID"); v == "" {
		t.Fatal("ALIYUN_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ALIYUN_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("ALIYUN_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALIYUN_DEFAULT_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-qingdao as test region")
		os.Setenv("ALIYUN_DEFAULT_REGION", "cn-qingdao")
	}
}
