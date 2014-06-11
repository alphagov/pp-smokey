Feature: spotlight

  @normal, @knownfailing
  Scenario: I can access the homepage
    When I GET https://spotlight.{PP_APP_DOMAIN}/performance
    Then I should receive an HTTP 200
      And I should see a strong ETag
