variable "aliyun_access_key" {}
variable "aliyun_secret_key" {}

provider "aliyun" {
  access_key = "${var.aliyun_access_key}"
  secret_key = "${var.aliyun_secret_key}"
}

resource "aliyun_ecs" "my-server" {
  image = "coreos-stable"
  name = "my-server"
  region = "nyc3"
  size = "512mb"
}
