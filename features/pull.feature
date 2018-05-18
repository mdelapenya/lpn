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
    The image [<repository>:<tag>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | tag | repository |
    | commerce | latest | liferay/liferay-commerce |
    | nightly | latest | mdelapenya/liferay-portal-nightlies |
    | release | latest | mdelapenya/liferay-portal |

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
    | commerce | latest | liferay/liferay-commerce |
    | nightly | latest | mdelapenya/liferay-portal-nightlies |
    | release | latest | mdelapenya/liferay-portal |

  Scenario Outline: Pull command forcing the removal of a non present image
    When I run `docker rmi -f <repository>:<tag>`
    And I run `lpn pull <type> -t <tag> -f`
    Then the output should contain:
    """
    The image [<repository>:<tag>] was not found in the local cache. Skipping removal.
    """
    And the exit status should be 0
    And I run `lpn checkImage <type> -t <tag>`
    And the output should contain:
    """
    The image [<repository>:<tag>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | tag | repository |
    | release | latest | mdelapenya/liferay-portal |