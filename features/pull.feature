Feature: Pull command
  As a newcomer to lpn
  I want to be able to pull the nightly builds or the releases images for Liferay Portal

  Scenario Outline: Pull command when image exists
    When I run `lpn pull <type> -t <tag>`
    Then the output should contain:
    """
   Pulling [<repository>:<tag>]
    """
    And I run `lpn checkImage <type> -t <tag>`
    And the output should contain:
    """
    The image [<repository>:<tag>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | tag | repository |
    | ce      | 7.0.6-ga7 | liferay/portal |
    | dxp     | 7.0.10.8 | liferay/dxp |
    | nightly | latest | mdelapenya/liferay-portal-nightlies |
    | release | latest | mdelapenya/liferay-portal |

  Scenario Outline: Pull command when image does not exist
    When I run `lpn pull <type> -t foo`
    Then the output should contain:
    """
    The image [<image>] could not be pulled
    """
    And the exit status should be 1

  Examples:
    | type    | image |
    | ce      | liferay/portal:foo |
    | dxp     | liferay/dxp:foo |
    | nightly | mdelapenya/liferay-portal-nightlies:foo |
    | release | mdelapenya/liferay-portal:foo |

  Scenario Outline: Pull command forcing the removal of already present image
    When I run `lpn pull <type> -t <tag>`
    And I run `lpn pull <type> -t <tag> -f`
    Then the output should contain:
    """
    [<repository>:<tag>] was deleted.
    """
    And I run `lpn checkImage <type> -t <tag>`
    And the output should contain:
    """
    The image [<repository>:<tag>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | tag | repository |
    | ce      | 7.0.6-ga7 | liferay/portal |
    | dxp     | 7.0.10.8 | liferay/dxp |
    | nightly | latest | mdelapenya/liferay-portal-nightlies |
    | release | latest | mdelapenya/liferay-portal |

  Scenario Outline: Pull command forcing the removal of a non present image
    When I run `docker rmi -f <repository>:<tag>`
    And I run `lpn pull <type> -t <tag> -f`
    Then the output should contain:
    """
    The image [<repository>:<tag>] was not found in the local cache. Skipping removal
    """
    And I run `lpn checkImage <type> -t <tag>`
    And the output should contain:
    """
    The image [<repository>:<tag>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | tag | repository |
    | release | latest | mdelapenya/liferay-portal |