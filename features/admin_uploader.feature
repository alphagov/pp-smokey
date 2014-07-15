Feature: admin_uploader

  @normal
  Scenario: Quickly loading the admin home page
    Given I am benchmarking
    When I visit the admin home page
    Then the elapsed time should be less than 1 seconds

  @normal
  @not_on_staging
  Scenario: Can log in
    When I try to login as a valid admin user
    Then I should be on the admin post-login page

  # Admin UI routing

  @normal
  Scenario: I can access the admin application without a trailing slash
    When I GET https://admin.{PP_APP_DOMAIN}/not-authorized
    Then I should receive an HTTP 200

  @normal
  Scenario: I can access the admin application with a trailing slash
    When I GET https://admin.{PP_APP_DOMAIN}/not-authorized/
    Then I should receive an HTTP 301 redirect to /not-authorized
