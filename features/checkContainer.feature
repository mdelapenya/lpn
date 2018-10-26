Feature: CheckContainer command
  As a newcomer to lpn
  I want to be able to check if the "lpn" container for each type is running

  Scenario Outline: CheckContainer command when container is running
    Given I run `lpn run <type> -t latest`
    When I run `lpn checkContainer <type>`
    Then the output should contain:
    """
    The container [lpn-<type>] DOES exist in the system
    """
    And the exit status should be 0
    And I run `lpn rm <type>`

  Examples:
    | type    |
    | nightly |
    | release |
  
  Scenario Outline: CheckContainer command when container is not running
    Given I run `lpn checkContainer <type>`
    Then the output should contain:
    """
    The container [lpn-<type>] does NOT exist in the system
    """
    And the exit status should be 1

  Examples:
    | type    |
    | nightly |
    | release |