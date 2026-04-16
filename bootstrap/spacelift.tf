resource "spacelift_aws_integration" "control_plane_integration" {
  name     = "${var.landing_zone_name}-aws-integration"
  role_arn = "arn:${data.aws_partition.current.id}:iam::${data.aws_caller_identity.current.account_id}:role/${var.landing_zone_name}-control-plan"
}
