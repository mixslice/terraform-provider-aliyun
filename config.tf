variable "aliyun_access_key" {}
variable "aliyun_secret_key" {}

provider "aliyun" {
  access_key = "${var.aliyun_access_key}"
  secret_key = "${var.aliyun_secret_key}"
}

data "aliyun_ecs_image" "ubuntu" {
  name_regex = "ubuntu"
  owner_alias = "system"
  region = "cn-qingdao"
  most_recent = true
}

resource "aliyun_ecs_instance" "test01" {
  image = "${data.aliyun_ecs_image.ubuntu.id}"
  name = "test01"
  region = "cn-qingdao"
  instance_type = "ecs.n1.tiny"
}
