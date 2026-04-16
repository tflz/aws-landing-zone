variable "landing_zone_name" {
  type = string
  description = "The display name of the landing zone."
}

variable "spacelift_parent_space_id" {
  type = string
  description = "ID of the parent space where the landing zone will be deployed. Defaults to root"
  default = "root"
}

variable "github_repository_name" {
  type = string
  description = "The name of the repository without owner part."
}

variable "github_repository_branch" {
  type = string
  description = "Git branch to apply changes from. Defaults to main."
  default = "main"
}
