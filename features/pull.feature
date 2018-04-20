Feature: Pull command
  As a newcomer to lpn
  I want to be able to pull the commerce, nightly builds or the releases images for Liferay Portal

  Scenario Outline: Pull command when image exists
    When I run `lpn pull <type> -t <tag>`
    Then the output should contain:
    """
    latest: Pulling from <repository>
    """
    And I run `lpn checkImage <type> -t <tag>`
    And the output should contain:
    """
    The image [<repository>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | tag | repository |
    | commerce | latest | liferay/liferay-commerce:latest |
    | nightly | latest | mdelapenya/liferay-portal-nightlies:latest |
    | release | latest | mdelapenya/liferay-portal:latest |

  Scenario Outline: Pull command when image does not exist
    When I run `lpn pull <type> -t foo`
    Then the output should contain:
    """
    Error response from daemon: manifest for <image> not found
    """
    And the exit status should be 1

  Examples:
    | type    | image |
    | commerce | liferay/liferay-commerce:foo |
    | nightly | mdelapenya/liferay-portal-nightlies:foo |
    | release | mdelapenya/liferay-portal:foo |