data "aws_partition" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_iam_role" "spacelift_integration" {
  name = "control_plane"
  path = "/${var.landing_zone_name}/spacelift/"

  assume_role_policy = core::jsonencode({
    Version = "2012-10-17",
    Statement = [
      core::jsonencode(data.spacelift_aws_integration_attachment_external_id.control_plane.assume_role_policy_statement)
    ]
  })
}

resource "aws_iam_role_policy_attachment" "control_plane" {
  role = aws_iam_role.spacelift_integration.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
