Feature: Pull command
  As a newcomer to lpn
  I want to be able to pull the nightly builds or the releases images for Liferay Portal

  Scenario Outline: Pull command when image exists
    When I run `lpn pull <type> -t <tag>`
    Then the output should contain:
    """
    <message>
    """
    And I run `lpn checkImage <type> -t <tag>`
    And the output should contain:
    """
    <checkMessage>
    """
    And the exit status should be 0

  Examples:
    | type    | tag | message | checkMessage |
    | nightly | latest | latest: Pulling from mdelapenya/liferay-portal-nightlies | The image [mdelapenya/liferay-portal-nightlies:latest] has been pulled from Docker Hub |
    | release | latest | latest: Pulling from mdelapenya/liferay-portal | The image [mdelapenya/liferay-portal:latest] has been pulled from Docker Hub |

  Scenario Outline: Pull command when image does not exist
    When I run `lpn pull <type> -t foo`
    Then the output should contain:
    """
    <message>
    """
    And the exit status should be 1

  Examples:
    | type    | message |
    | nightly | Error response from daemon: manifest for mdelapenya/liferay-portal-nightlies:foo not found |
    | release | Error response from daemon: manifest for mdelapenya/liferay-portal:foo not found |