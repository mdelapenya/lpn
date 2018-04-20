Feature: Stop command
  As a newcomer to lpn
  I want to be able to stop the container created by the tool

  Scenario Outline: Stop command when container exists
    When I run `lpn run <type> -t <tag>`
    And I run `lpn stop <type>`
    Then the output should contain:
    """
    lpn-<type>
    """
    And the exit status should be 0
    And I run `lpn rm <type>`

  Examples:
    | type    | tag |
    | commerce | latest |
    | nightly | latest |
    | release | latest |

  Scenario Outline: Stop command when container does not exist
    When I run `lpn stop <type>`
    Then the output should contain:
    """
    Impossible to stop the container [lpn-<type>]
    """
    And the exit status should be 1
  
  Examples:
    | type    | tag |
    | commerce | latest |
    | nightly | latest |
    | release | latest |