Feature: Rmi command
  As a newcomer to lpn
  I want to be able to remove the images used by the tool

  Scenario Outline: Rmi command when image exists
    Given I run `lpn pull <type> -t <tag>`
    When I run `lpn rmi <type> -t <tag>`
    Then the output should contain:
    """
    Image has been removed
    """
    And the output should contain:
    """
    image="docker.io/<image>:<tag>"
    """
    And the exit status should be 0

  Examples:
    | type    | image | tag |
    | ce      | liferay/portal | 7.0.6-ga7 |
    | commerce | liferay/commerce | 1.1.1 |
    | dxp     | liferay/dxp | 7.0.10.8 |
    | nightly | liferay/portal-snapshot | master |
    | release | mdelapenya/liferay-portal | latest |

  Scenario Outline: Rmi command when image does not exist
    Given I run `lpn rmi <type> -t <tag>`
    Then the output should contain:
    """
    Impossible to remove the image
    """
    And the output should contain:
    """
    image="docker.io/<image>:<tag>"
    """
    And the exit status should be 0

  Examples:
    | type    | image | tag |
    | ce      | liferay/portal | foo |
    | commerce | liferay/commerce | foo |
    | dxp     | liferay/dxp | foo |
    | nightly | liferay/portal-snapshot | foo |
    | release | mdelapenya/liferay-portal | foo |