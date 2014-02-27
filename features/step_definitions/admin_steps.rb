When /^I try to login as a valid admin user$/ do
  assert ENV["SIGNON_USERNAME"] && ENV["SIGNON_PASSWORD"], "Please ensure that the signon user credentials are available in the environment"

  # Need to do it this way to comply with CSRF protection
  visit "#{admin_base_url}sign-in"
  fill_in "user_email", :with => ENV["SIGNON_USERNAME"]
  fill_in "user_password", :with =>ENV["SIGNON_PASSWORD"]
  click_button "Sign in"
end

When /^I visit the admin home page$/ do
  visit "#{application_base_url('admin')}"
end

Then /^I should be on the admin home page$/ do
  page.has_selector?("#user_username").should == true # username input field
  page.has_selector?("#user_password").should == true # password input field
end

Then /^I should be on the admin post-login page$/ do
  page.has_selector?(".alert-success").should == true, page.text # Signed in successfully message
  page.has_text?("Signed in as ").should == true, page.text #page has "Signed in as X" message
  page.has_text?("Sign out").should == true, page.text # page has a logout link
end
