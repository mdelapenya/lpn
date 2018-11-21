Feature: Rm command
  As a newcomer to lpn
  I want to be able to remove the container created by the tool

  Scenario Outline: Rm command when container exists
    When I run `lpn run <type> -t <tag>`
    And I run `lpn rm <type>`
    Then the output should contain:
    """
    lpn-<type>
    """
    And the exit status should be 0

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | dxp     | 7.0.10.8 |
    | nightly | latest |
    | release | latest |

  Scenario Outline: Rm command when container does not exist
    When I run `lpn rm <type>`
    Then the output should contain:
    """
    Impossible to remove the container [lpn-<type>]
    """
    And the exit status should be 1

  Examples:
    | type    | tag |
    | ce      | 7.0.6-ga7 |
    | dxp     | 7.0.10.8 |
    | nightly | latest |
    | release | latest |