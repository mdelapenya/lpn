Feature: Stop command
  As a newcomer to lpn
  I want to be able to stop the container created by the tool

  Scenario Outline: Stop command when container exists
    When I run `lpn run <type> -t <tag>`
    And I run `lpn stop`
    Then the output should contain:
    """
    lpn-<type>
    """
    And the exit status should be 0
    And I run `lpn rm`

  Examples:
    | type    | tag |
    | commerce | latest |
    | nightly | latest |
    | release | latest |

  Scenario: Stop command when container does not exist
    When I run `lpn stop`
    Then the output should contain:
    """
    Impossible to stop the container
    """
    And the exit status should be 1