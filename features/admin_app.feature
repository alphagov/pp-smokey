Feature: admin app

  @normal
  Scenario: Quickly being redirected to Sign-on-o-tron
    Given I am benchmarking
    When I GET https://admin-beta.{PP_APP_DOMAIN}/login
    Then I should receive an HTTP redirect beginning with https://signon.{GOVUK_APP_DOMAIN}/oauth/authorize
    And the elapsed time should be less than 1 second
