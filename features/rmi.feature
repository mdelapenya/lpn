Feature: Rmi command
  As a newcomer to lpn
  I want to be able to remove the images used by the tool

  Scenario Outline: Rm command when image exists
    When I run `lpn pull <type> -t <tag>`
    And I run `lpn rmi <type> -t <tag>`
    Then the output should contain:
    """
    [<image>:<tag>] was deleted.
    """
    And the exit status should be 0

  Examples:
    | type    | image | tag |
    | ce      | liferay/portal | 7.0.6-ga7 |
    | dxp     | liferay/dxp | 7.0.10.8 |
    | nightly | mdelapenya/liferay-portal-nightlies | latest |
    | release | mdelapenya/liferay-portal | latest |

  Scenario Outline: Rm command when image does not exist
    When I run `lpn rmi <type> -t <tag>`
    Then the output should contain:
    """
    Impossible to remove the image [<image>:<tag>]
    """
    And the exit status should be 1

  Examples:
    | type    | image | tag |
    | ce      | liferay/portal | foo |
    | dxp     | liferay/dxp | foo |
    | nightly | mdelapenya/liferay-portal-nightlies | foo |
    | release | mdelapenya/liferay-portal | foo |