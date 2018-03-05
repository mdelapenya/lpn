Feature: CheckContainer command
  As a newcomer to lpn
  I want to be able to check if the container "liferay-portal-nook" is running

  Scenario Outline: CheckContainer command when container is running
    Given I run `lpn run <type> -t latest`
    When I run `lpn checkContainer`
    Then the output should contain:
    """
    The container [liferay-portal-nook] is running
    """
    And the exit status should be 0
    And I run `lpn rm`

  Examples:
    | type    |
    | nightly |
    | release |
  
  Scenario: CheckContainer command when container is not running
    Given I run `lpn checkContainer`
    Then the output should contain:
    """
    The container [liferay-portal-nook] is NOT running
    """
    And the exit status should be 1