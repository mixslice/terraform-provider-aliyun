provider "aliyun" {
  api_key     = "s3cur3t0k3n=="
  endpoint    = "https://api.example.org/v1"
  timeout     = 60
  max_retries = 5
}

resource "aliyun_instance" "my-speedy-server" {
  name = "speedracer"
  cpus = 4
  ram  = 16384
}
