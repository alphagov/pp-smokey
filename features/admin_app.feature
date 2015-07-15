Feature: admin app

  Scenario: Admin app is up
    When I GET https://admin.{PP_APP_DOMAIN}/
    Then I should receive an HTTP 302

  @normal
  @not_on_staging
  Scenario: Quickly being redirected to Sign-on-o-tron
    Given I am benchmarking
    When I GET https://admin.{PP_APP_DOMAIN}/
    Then I should receive an HTTP redirect beginning with https://signon.{GOVUK_APP_DOMAIN}/oauth/authorize
    And the elapsed time should be less than 1 second

  @normal
  @not_on_staging
  Scenario: Can log in to the admin app using Sign-on-o-tron
    When I try to login to Signon from https://admin.{PP_APP_DOMAIN}/
    Then I should be on a page with a URL that begins https://admin.{PP_APP_DOMAIN}/

  @normal
  @not_on_staging
  Scenario: Can see a list of data sets
    When I try to login to Signon from https://admin.{PP_APP_DOMAIN}/
    Then I should see a list of test data sets containing test data type

  @normal
  @not_on_staging
  Scenario: Can upload to a data set
    When I try to login to Signon from https://admin.{PP_APP_DOMAIN}/
    And I upload fixtures/test-data.csv to the test data type in the test data group
    Then I should see a success message for the test data type in the test data group

  @normal
  @not_on_staging
  Scenario: Can see the dashboard administration page
    When I try to login to Signon from https://admin.{PP_APP_DOMAIN}/
    And I follow the Your dashboards link
    Then I should be on the dashboard administration page
