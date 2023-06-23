locals {
  # subscription_id_alias is the id of the newly created subscription, if it exists.
  subscription_id_alias = try(azurerm_subscription.this[0].subscription_id, null)

  # subscription_id is the id of the newly created subscription, or the id supplied by var.subscription_id.
  subscription_id = coalesce(local.subscription_id_alias, var.subscription_id)

  # get a set of the RP names and features to register
  resource_provider_feature_set = toset(flatten([
    for rp, features in var.subscription_register_resource_providers_and_features : [
      for feature in features : {
        resource_provider_name = rp
        feature_name           = feature
      } if length(features) > 0
    ]
  ]))

  # Turn the above into a map for the resource for_each
  resource_provider_feature_map = {
    for i in local.resource_provider_feature_set : "${i.resource_provider_name}/${i.feature_name}" => {
      resource_provider_name = i.resource_provider_name
      feature_name           = i.feature_name
  } }
}
