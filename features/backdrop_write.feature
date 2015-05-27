Feature: backdrop_write

  @normal
  Scenario: I can write to Backdrop
    Given I have the JSON data []
      And I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test
    Then I should receive an HTTP 200

  @normal
  Scenario: I cannot write to Backdrop without an Authorization header
    Given I have the JSON data []
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test
    Then I should receive an HTTP 401

  @normal
  Scenario: I cannot write to Backdrop using an inappropriate Authorization header
    Given I have the JSON data []
      And I have the HTTP header "Authorization: some-other-thing"
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test
    Then I should receive an HTTP 401

  @normal
  Scenario: I can PUT an empty JSON list to backdrop (empty functionality)
    Given I have the JSON data []
      And I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I PUT https://www.{PP_APP_DOMAIN}/data/test/test
    Then I should receive an HTTP 200

  @normal
  Scenario: I cannot DELETE resources in backdrop
    Given I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I DELETE https://www.{PP_APP_DOMAIN}/data/test/test
    Then I should receive an HTTP 405
