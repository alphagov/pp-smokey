When /^I try to login to Signon from (.*)$/ do |url|
  assert ENV["SIGNON_USERNAME"] && ENV["SIGNON_PASSWORD"], "Please ensure that the signon user credentials are available in the environment"

  # We need to actually visit the URL to make the CSRF protection happy
  url = replace_env_host(url)
  visit url

  fill_in "user_email", :with => ENV["SIGNON_USERNAME"]
  fill_in "user_password", :with =>ENV["SIGNON_PASSWORD"]
  click_button "Sign in"
end

When /^I upload (.*) to the (.*) data type in the (.*) data group/ do |path, data_type, data_group|
  attach_file("#{data_group}-#{data_type}-file", path)
  page.find(:css, "##{data_group}-#{data_type}-form button[type=submit]").click
end

Then /^I should see a list of (.*) data sets containing (.*) data type$/ do |data_set_group_name, data_set_type_name|
  page.find(:css, "li[data-name='#{data_set_group_name}'] > ul.data-set-list > li .data-set-type").text.should == data_set_type_name
end

Then /^I should see a success message for the (.*) data type in the (.*) data group/ do |data_type, data_group|
  page.find(:css, "##{data_group}-#{data_type}-form + div.upload-messages").text.should have_content('uploaded successfully')
end

When /^I follow the "(.*)" link$/ do |link_text|
  click_link link_text 
end

When /^I follow the Administer dashboards link$/ do
  page.find(:css, "[href='/administer-dashboards']").click
end

Then /^I should be on the dashboard administration page$/ do
  h1 = page.find(:css, "h1")
  h1.text.should == "Administer dashboards"
end
