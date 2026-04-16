resource "spacelift_space" "landing_zone" {
  name            = var.landing_zone_name
  parent_space_id = var.spacelift_parent_space_id
}

resource "spacelift_stack" "control_plane" {
  name       = "${var.landing_zone_name}-control-plane"
  autodeploy = true

  repository = var.github_repository_name
  branch     = var.github_repository_branch
  project_root = "stacks/control-plane"
}

resource "spacelift_aws_integration" "control_plane" {
  name     = "${var.landing_zone_name}-aws-integration"
  role_arn = "arn:${data.aws_partition.current.id}:iam::${data.aws_caller_identity.current.account_id}:role/${var.landing_zone_name}/spacelift/control-plane"
  space_id = spacelift_space.landing_zone.id
}

resource "spacelift_aws_integration_attachment" "control_plane" {
  integration_id = spacelift_aws_integration.control_plane.id
  stack_id       = spacelift_stack.control_plane.id 
  read           = true
  write          = true

  depends_on = [ aws_iam_role.spacelift_integration ]
}

data "spacelift_aws_integration_attachment_external_id" "control_plane" {
  integration_id = spacelift_aws_integration.control_plane.id
  stack_id       = spacelift_stack.control_plane.id 
  read           = true
  write          = true
}
