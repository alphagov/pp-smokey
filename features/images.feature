Feature: image fallbacks

  @normal
  Scenario: Image fallbacks look vaguely correct
    Given I am benchmarking
    When I GET https://spotlight.{PP_FULL_APP_DOMAIN}/performance/carers-allowance/volumetrics.png
    Then I should receive an HTTP 200
     And the elapsed time should be less than 4 seconds
     And the image should be between 950 and 970 pixels wide
     And the image should be between 400 and 410 pixels high