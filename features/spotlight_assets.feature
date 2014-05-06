Feature: spotlight_assets

  @normal
  Scenario: I can access static assets
    When I GET https://assets.{PP_APP_DOMAIN}/spotlight/stylesheets/spotlight.css
    Then I should receive an HTTP 200
