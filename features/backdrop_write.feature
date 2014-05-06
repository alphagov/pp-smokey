Feature: backdrop_write

  @normal
  Scenario: I can write to Backdrop
    Given I have the JSON data []
      And I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test
    Then I should receive an HTTP 200

  @normal
  Scenario: I cannot write to Backdrop if I use a trailing slash
    Given I have the JSON data []
      And I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test/
    Then I should receive an HTTP 404

  @normal
  Scenario: I cannot write to Backdrop without an Authorization header
    Given I have the JSON data []
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test?limit=1
    Then I should receive an HTTP 401

  @normal
  Scenario: I cannot write to Backdrop using an inappropriate Authorization header
    Given I have the JSON data []
      And I have the HTTP header "Authorization: some-other-thing"
    When I POST https://www.{PP_APP_DOMAIN}/data/test/test?limit=1
    Then I should receive an HTTP 401

  @normal
  Scenario: I cannot PUT to backdrop
    Given I have the JSON data []
      And I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I PUT https://www.{PP_APP_DOMAIN}/data/test/test?limit=1
    Then I should receive an HTTP 405

  @normal
  Scenario: I cannot DELETE resources in backdrop
    Given I have the HTTP header "Authorization: Bearer qwertyuiop"
    When I DELETE https://www.{PP_APP_DOMAIN}/data/test/test?limit=1
    Then I should receive an HTTP 405
