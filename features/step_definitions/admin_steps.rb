When /^I try to login to Signon from (.*)$/ do |url|
  assert ENV["SIGNON_USERNAME"] && ENV["SIGNON_PASSWORD"], "Please ensure that the signon user credentials are available in the environment"

  # We need to actually visit the URL to make the CSRF protection happy
  url = replace_env_host(url)
  visit url

  fill_in "user_email", :with => ENV["SIGNON_USERNAME"]
  fill_in "user_password", :with =>ENV["SIGNON_PASSWORD"]
  click_button "Sign in"
end

When /^I visit the admin home page$/ do
  visit "#{application_base_url('admin')}"
end

Then /^I should be on the admin home page$/ do
  page.has_selector?("#user_username").should eq(true) # username input field
  page.has_selector?("#user_password").should eq(true) # password input field
end

Then /^I should be on the admin post-login page$/ do
  page.has_selector?(".alert-success").should eq(true),  page.text # Signed in successfully message
  page.has_text?("Signed in as ").should eq(true), page.text #page has "Signed in as X" message
  page.has_text?("Sign out").should eq(true), page.text # page has a logout link
end

Then /^I should see a list of data sets containing (.*)\/(.*)$/ do |data_group, data_type|
  page.find(:css, ".data-set-list .data-set").has_text?(data_group)
  page.find(:css, ".data-set-list .data-set").has_text?(data_type)
end
