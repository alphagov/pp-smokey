Given /^the "(.*)" application has booted$/ do |app_name|
  url = application_base_url(app_name)
  head_request(url)
end
