Feature: admin app

  @normal
  @not_on_staging
  Scenario: Quickly being redirected to Sign-on-o-tron
    Given I am benchmarking
    When I GET https://admin-beta.{PP_APP_DOMAIN}/login
    Then I should receive an HTTP redirect beginning with https://signon.{GOVUK_APP_DOMAIN}/oauth/authorize
    And the elapsed time should be less than 1 second

  @normal
  @not_on_staging
  Scenario: Can log in to the admin app using Sign-on-o-tron
    When I try to login to Signon from https://admin-beta.{PP_APP_DOMAIN}/login
    Then I should be on a page with a URL that begins https://admin-beta.{PP_APP_DOMAIN}/

  @normal
  @not_on_staging
  Scenario: Can see a list of data sets
    When I try to login to Signon from https://admin-beta.{PP_APP_DOMAIN}/login
    Then I should see a list of data sets containing test
