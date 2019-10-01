variable "ssh_public_key_path" {
  type = string
  default = "~/.ssh/id_rsa.pub"
}

variable "ssh_bastion_host" {
  type = string
  default = ""
}

variable "ssh_bastion_user" {
  type = string
  default = ""
}

variable "image" {
  type = string
  default = "ubuntu-18-04-cloud-amd64"
}

variable "load_generator_flavor" {
  type = string
  default = "m1.medium"
}

variable "nginx_on_docker_flavor" {
  type = string
  default = "m1.large"
}

provider "openstack" {
}

resource "openstack_compute_keypair_v2" "ssh" {
  name = "ssh"
  public_key = "${file(pathexpand(var.ssh_public_key_path))}"
}

resource "openstack_compute_servergroup_v2" "cfs-test" {
  name     = "cfs-test"
  policies = ["anti-affinity"]
}

resource "openstack_compute_instance_v2" "nginx-on-docker" {
  name = "nginx-on-docker"
  image_name = var.image
  flavor_name = var.nginx_on_docker_flavor
  key_pair = openstack_compute_keypair_v2.ssh.name
  network {
    name = "shared"
  }
  scheduler_hints {
    group = openstack_compute_servergroup_v2.cfs-test.id
  }

  connection {
    host = self.network[0].fixed_ip_v4
    user = "ubuntu"
    bastion_host = var.ssh_bastion_host
    bastion_user = var.ssh_bastion_user
  }
  provisioner "remote-exec" {
    inline = [
      "echo Ready!"
    ]
  }
}

resource "ansible_host" "nginx-on-docker" {
  inventory_hostname = openstack_compute_instance_v2.nginx-on-docker.name
  vars = {
    ansible_user = "ubuntu"
    ansible_host = openstack_compute_instance_v2.nginx-on-docker.network[0].fixed_ip_v4
    ansible_become = "yes"
  }
}

resource "openstack_compute_instance_v2" "load-generator" {
  name = "load-generator"
  image_name = var.image
  flavor_name = var.load_generator_flavor
  key_pair = openstack_compute_keypair_v2.ssh.name
  network {
    name = "shared"
  }
  scheduler_hints {
    group = openstack_compute_servergroup_v2.cfs-test.id
  }

  connection {
    host = self.network[0].fixed_ip_v4
    user = "ubuntu"
    bastion_host = var.ssh_bastion_host
    bastion_user = var.ssh_bastion_user
  }
  provisioner "remote-exec" {
    inline = [
      "echo Ready!"
    ]
  }
}

resource "ansible_host" "load-generator" {
  inventory_hostname = openstack_compute_instance_v2.load-generator.name
  vars = {
    ansible_user = "ubuntu"
    ansible_host = openstack_compute_instance_v2.load-generator.network[0].fixed_ip_v4
    ansible_become = "yes"
  }
}
