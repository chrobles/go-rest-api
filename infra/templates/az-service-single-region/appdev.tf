data "azurerm_container_registry" "acr" {
  name                = local.resolved_acr_name
  resource_group_name = local.resolved_acr_rg_name
  depends_on          = ["azurerm_resource_group.svcplan", "module.container_registry"]
}

data "azurerm_client_config" "current" {}

# Build and push the app service image slot to enable continuous deployment scenarios. We're using ACR build tasks to remotely carry the docker build / push.
resource "null_resource" "acr_image_deploy" {
  count      = length(var.deployment_targets)
  depends_on = ["module.container_registry"]

  triggers = {
    images_to_deploy = "${join(",", [for target in var.deployment_targets : "${target.image_name}:${target.image_release_tag_prefix}"])}"
  }

  provisioner "local-exec" {
    command = <<EOF
      az acr build                              \
        --subscription "$SUBSCRIPTION_ID"       \
        --resource-group "$RESOURCE_GROUP_NAME" \
        -t $IMAGE:$TAG                          \
        -r $REGISTRY                            \
        -f $DOCKERFILE $SOURCE
      EOF

    environment = {
      SUBSCRIPTION_ID = data.azurerm_client_config.current.subscription_id
      RESOURCE_GROUP_NAME = local.resolved_acr_rg_name
      IMAGE = var.deployment_targets[count.index].image_name
      TAG = var.deployment_targets[count.index].image_release_tag_prefix
      REGISTRY = module.container_registry.container_registry_name
      DOCKERFILE = var.acr_build_docker_file
      SOURCE = var.acr_build_git_source_url
    }

  }
}

module "service_plan" {
  source = "../../modules/providers/azure/service-plan"
  resource_group_name = azurerm_resource_group.svcplan.name
  service_plan_name = local.sp_name
}

module "app_service" {
  source = "../../modules/providers/azure/app-service"
  app_service_name_prefix = local.app_svc_name_prefix
  service_plan_name = module.service_plan.service_plan_name
  service_plan_resource_group_name = azurerm_resource_group.svcplan.name
  uses_acr = true
  azure_container_registry_name = module.container_registry.container_registry_name
  docker_registry_server_url = module.container_registry.container_registry_login_server
  docker_registry_server_username = data.azurerm_container_registry.acr.admin_username
  docker_registry_server_password = data.azurerm_container_registry.acr.admin_password
  app_service_config = {
    for target in var.deployment_targets :
    target.app_name => {
      image = "${target.image_name}:${target.image_release_tag_prefix}"
    }
  }
}

resource "azurerm_role_assignment" "acr_pull" {
  count = length(var.deployment_targets)
  scope = module.container_registry.container_registry_id
  role_definition_name = "AcrPull"
  principal_id = module.app_service.app_service_identity_object_ids[count.index]
}

