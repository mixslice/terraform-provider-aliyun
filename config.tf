variable "aliyun_access_key" {}
variable "aliyun_secret_key" {}

provider "aliyun" {
  access_key = "${var.aliyun_access_key}"
  secret_key = "${var.aliyun_secret_key}"
  region = "cn-beijing"
}

data "aliyun_ecs_image" "ubuntu" {
  name_regex = "ubuntu"
  owner_alias = "system"
  most_recent = true
}

resource "aliyun_ecs_instance" "test01" {
  image = "${data.aliyun_ecs_image.ubuntu.id}"
  name = "test01"
  instance_type = "ecs.n1.tiny"
}
