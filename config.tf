variable "aliyun_access_key" {}
variable "aliyun_secret_key" {}

provider "aliyun" {
  access_key = "${var.aliyun_access_key}"
  secret_key = "${var.aliyun_secret_key}"
}

resource "aliyun_ecs_instance" "test01" {
  image = "ubuntu1404_64_40G_aliaegis_20160222.vhd"
  name = "test01"
  region = "cn-qingdao"
  instance_type = "ecs.t1.small"
}
