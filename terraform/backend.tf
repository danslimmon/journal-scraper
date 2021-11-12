terraform {
  backend "remote" {
    organization = "danslimmon"

    workspaces {
      name = "journal-scraper"
    }
  }
}
