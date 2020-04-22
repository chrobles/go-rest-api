resource "azurerm_resource_group" "svcplan" {
  name     = local.admin_rg_name
  location = local.region
}

module "container_registry" {
  source                           = "../../modules/providers/azure/container-registry"
  container_registry_name          = local.resolved_acr_name
  resource_group_name              = local.resolved_acr_rg_name
  container_registry_admin_enabled = true
  container_registry_tags          = var.azure_container_tags
}
