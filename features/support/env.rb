require 'nokogiri'
require 'capybara/cucumber'
require 'capybara/mechanize'

def base_url
  ENV["PP_WEBSITE_ROOT"] || "https://www.preview.performance.service.gov.uk"
end

Capybara.default_driver = :mechanize
Capybara.app = "not nil"
Capybara.app_host = base_url
