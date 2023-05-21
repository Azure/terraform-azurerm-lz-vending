locals {
  public_routing_policy = {
    destinations = [
      "Internet"
    ]
    name    = "PublicTraffic"
    nextHop = ""
  }
  private_routing_policy = {
    destinations = [
      "PrivateTraffic"
    ]
    name    = "PrivateTraffic"
    nextHop = ""
  }
}
