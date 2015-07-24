Feature: spotlight

  @normal
  Scenario: I can access the homepage
    Given I am benchmarking
    When I GET https://spotlight.{GOVUK_APP_DOMAIN}/performance
    Then I should receive an HTTP 200
      And I should see a strong ETag
