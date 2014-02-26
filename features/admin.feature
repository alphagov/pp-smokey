Feature: admin

  @normal
  Scenario: Quickly loading the admin home page
    Given the "admin" application has booted
      And I am benchmarking
    When I visit the admin home page
    Then the elapsed time should be less than 1 seconds

  @normal
  Scenario: Can log in
    Given the "admin" application has booted
    When I try to login as a valid admin user
    Then I should be on the admin post-login page
