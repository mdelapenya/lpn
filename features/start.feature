Feature: Start command
  As a newcomer to lpn
  I want to be able to start the container created by the tool

  Scenario Outline: Start command when container exists and it's started
    Given I run `lpn run <type> -t <tag>`
    And I run `lpn start <type>`
    Then the output should contain:
    """
    [lpn-<type>] started
    """
    And the exit status should be 0
    And I run `lpn rm <type>`

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | commerce | 1.1.1 |
    | dxp     | 7.0.10.8 |
    | nightly | master |
    | release | latest |

  Scenario Outline: Start command when container exists and it's stopped
    Given I run `lpn run <type> -t <tag>`
    When I run `lpn stop <type>`
    And I run `lpn start <type>`
    Then the output should contain:
    """
    [lpn-<type>] started
    """
    And the exit status should be 0
    And I run `lpn rm <type>`

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | commerce | 1.1.1 |
    | dxp     | 7.0.10.8 |
    | nightly | master |
    | release | latest |

  Scenario Outline: Start command when container does not exist
    Given I run `lpn start <type>`
    Then the output should contain:
    """
    Impossible to start the container [lpn-<type>]
    """
    And the output should contain:
    """
    No such container: lpn-<type>
    """
    And the exit status should be 1
  
  Examples:
    | type    |
    | ce      |
    | commerce |
    | dxp     |
    | nightly |
    | release |