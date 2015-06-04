Feature: stagecraft

  @normal
  Scenario: I can route a request to Stagecraft data-sets API without a trailing slash
    When I GET https://stagecraft.{PP_FULL_APP_DOMAIN}/data-sets
    Then I should receive an HTTP 403
    # Requires a secret bearer token - getting a 403 is enough to know we
    # routed right through to the app.

  @normal
  Scenario: I can access Stagecraft admin UI with a trailing slash
    When I GET https://stagecraft.{PP_FULL_APP_DOMAIN}/admin/
    Then I should receive an HTTP 302 redirect beginning with https://stagecraft.{PP_FULL_APP_DOMAIN}/admin/

  @normal
  Scenario: I can access Stagecraft admin UI without a trailing slash
    When I GET https://stagecraft.{PP_FULL_APP_DOMAIN}/admin
    Then I should receive an HTTP 301 redirect to https://stagecraft.{PP_FULL_APP_DOMAIN}/admin/
