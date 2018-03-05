Feature: Rm command
  As a newcomer to lpn
  I want to be able to remove the container created by the tool

  Scenario Outline: Rm command when container exists
    When I run `lpn run <type> -t <tag>`
    And I run `lpn rm`
    Then the output should contain:
    """
    <message>
    """
    And the exit status should be 0

  Examples:
    | type    | tag | message |
    | nightly | latest | liferay-portal-nook |
    | release | latest | liferay-portal-nook |

  Scenario: Rm command when container does not exist
    When I run `lpn rm`
    Then the output should contain:
    """
    Impossible to remove the container
    """
    And the exit status should be 1