Feature: Stop command
  As a newcomer to lpn
  I want to be able to stop the container created by the tool

  Scenario Outline: Stop command when container exists
    Given I run `lpn run <type> -t <tag>`
    When I run `lpn stop <type>`
    Then the output should contain:
    """
    [lpn-<type>] stopped
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

  Scenario Outline: Stop command when container does not exist
    Given I run `lpn stop <type>`
    Then the output should contain:
    """
    Impossible to stop the container [lpn-<type>]
    """
    And the exit status should be 1
  
  Examples:
    | type    |
    | ce      |
    | commerce |
    | dxp     |
    | nightly |
    | release |