Feature: Pull command
  As a newcomer to lpn
  I want to be able to pull the nightly builds or the releases images for Liferay Portal

  Scenario Outline: Pull command when image exists
    Given I run `lpn pull <type> -t <tag>`
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
    | commerce | 1.1.1 | liferay/commerce |
    | dxp     | 7.0.10.8 | liferay/dxp |
    | nightly | master | liferay/portal-snapshot |
    | release | latest | liferay/liferay-portal |

  Scenario Outline: Pull command when image does not exist
    Given I run `lpn pull <type> -t foo`
    Then the output should contain:
    """
    The image [<image>] could not be pulled
    """
    And the exit status should be 1

  Examples:
    | type    | image |
    | ce      | liferay/portal:foo |
    | commerce | liferay/commerce:foo |
    | dxp     | liferay/dxp:foo |
    | nightly | liferay/portal-snapshot:foo |
    | release | liferay/liferay-portal:foo |

  Scenario Outline: Pull command forcing the removal of already present image
    Given I run `lpn pull <type> -t <tag>`
    When I run `lpn pull <type> -t <tag> -f`
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
    | commerce | 1.1.1 | liferay/commerce |
    | dxp     | 7.0.10.8 | liferay/dxp |
    | nightly | master | liferay/portal-snapshot |
    | release | latest | liferay/liferay-portal |

  Scenario Outline: Pull command forcing the removal of a non present image
    Given I run `docker rmi -f <repository>:<tag>`
    When I run `lpn pull <type> -t <tag> -f`
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
    | release | latest | liferay/liferay-portal |