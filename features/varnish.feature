Feature: varnish

  # Static assets

  @normal
  Scenario: I can access static assets
    When I GET https://assets.{PP_APP_DOMAIN}/spotlight/stylesheets/spotlight.css
    Then I should receive an HTTP 200

  # Backdrop routing

  @normal
  Scenario: I can access Backdrop without a trailing slash
    When I GET https://www.{PP_APP_DOMAIN}/data/test/test?limit=1
    Then I should receive an HTTP 200

  @normal
  Scenario: I can access Backdrop with a trailing slash
    When I GET https://www.{PP_APP_DOMAIN}/data/test/test/
    Then I should receive an HTTP 301 redirect to /data/test/test


  # BUG: see https://www.pivotaltracker.com/story/show/67097376
  @normal, @knownfailing
  Scenario: I can access Backdrop with a trailing slash and query parameters
    When I GET https://www.{PP_APP_DOMAIN}/data/test/test/?limit=1
    Then I should receive an HTTP 301 redirect to www.{PP_APP_DOMAIN}/data/test/test?limit=1

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


  # Stagecraft routing

  @normal
  Scenario: I can access Stagecraft data-sets API without a trailing slash
    When I GET https://stagecraft.{PP_APP_DOMAIN}/data-sets
    Then I should receive an HTTP 200


  # BUG: see https://www.pivotaltracker.com/story/show/67096896
  @normal, @knownfailing
  Scenario: I can access Stagecraft data-sets API with a trailing slash
    When I GET https://stagecraft.{PP_APP_DOMAIN}/data-sets/
    Then I should receive an HTTP 301 redirect to /data-sets

  @normal
  Scenario: I can access Stagecraft admin UI with a trailing slash
    When I GET https://stagecraft.{PP_APP_DOMAIN}/admin/
    Then I should receive an HTTP 200


  @normal
  Scenario: I can access Stagecraft admin UI without a trailing slash
    When I GET https://stagecraft.{PP_APP_DOMAIN}/admin
    Then I should receive an HTTP 301 redirect to https://stagecraft.{PP_APP_DOMAIN}/admin/


  # Admin UI routing

  @normal
  Scenario: I can access the admin application without a trailing slash
    When I GET https://admin.{PP_APP_DOMAIN}/not-authorized
    Then I should receive an HTTP 200

  @normal
  Scenario: I can access the admin application with a trailing slash
    When I GET https://admin.{PP_APP_DOMAIN}/not-authorized/
    Then I should receive an HTTP 301 redirect to /not-authorized
