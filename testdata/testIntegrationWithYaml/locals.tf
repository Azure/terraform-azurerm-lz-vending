locals {
  # landing_zone_data_dir is the directory containing the YAML files for the landing zones.
  landing_zone_data_dir = "${path.root}/data"

  # landing_zone_files is the list of landing zone YAML files to be processed
  landing_zone_files = fileset(local.landing_zone_data_dir, "landing_zone_*.yaml")

  # landing_zone_data_map is the decoded YAML data stored in a map
  landing_zone_data_map = {
    for f in local.landing_zone_files :
    f => yamldecode(file("${local.landing_zone_data_dir}/${f}"))
  }
}
