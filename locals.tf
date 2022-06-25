locals {
  # subscription_id is the id of the subscription into which resources will be created.
  # We pick the created sub id first, if it exists, otherwise we pick the subscription_id variable.
  # If this is blank then the subscription submodule is disabled an no subscription id has been supplied as an input variabke.
  subscription_id = coalesce(module.subscription[0].subscription_id, var.subscription_id, "")

  # subscription_resource_id is the Azure resource id of the subscription id that was supplied in the inpuit variables.
  # If var.subscription_id is empty, then we will return en empty string so that we can correctly coalesce the subscription_resource_id output.
  supplied_subscription_resource_id = var.subscription_id == "" ? "" : "/subscriptions/${var.subscription_id}"

  # subscription_resource_id is the Azure resource id of the subscription into which resources will be created.
  # We use the created sub resoruce id first, if it exists, otherwise we pick the subscription_id variable.
  # If this is blank then the subscription submodule is disabled an no subscription id has been supplied as an input variabke.
  subscription_resource_id = coalesce(module.subscription[0].subscription_resource_id, local.supplied_subscription_resource_id, "")
}
