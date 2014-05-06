Feature: backdrop_read

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


