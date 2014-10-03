Feature: spotlight

  @normal
  Scenario: I can access the homepage
    Given I am benchmarking
    When I GET https://spotlight.{PP_FULL_APP_DOMAIN}/performance
    Then I should receive an HTTP 200
      And I should see a strong ETag
      And the elapsed time should be less than 2 second
