Feature: checki command
  As a newcomer to lpn
  I want to be able to check if the images representing nightly builds or releases are present

  Scenario Outline: checki command for latest tag
    Given I run `lpn pull <type> -t <tag>`
    When I run `lpn checki <type> -t <tag>`
    Then the output should contain:
    """
    Image has been pulled from Docker Hub
    """
    And the output should contain:
    """
    image="docker.io/<image>:<tag>"
    """
    And the exit status should be 0

  Examples:
    | type    | image | tag |
    | nightly | mdelapenya/portal-snapshot | master |
    | release | mdelapenya/liferay-portal | latest |

  Scenario Outline: checki command
    Given I run `lpn pull <type> -t <tag>`
    When I run `lpn checki <type> -t <tag>`
    Then the output should contain:
    """
    Image has been pulled from Docker Hub
    """
    And the output should contain:
    """
    image="docker.io/<image>:<tag>"
    """
    And the exit status should be 0

  Examples:
    | type | image       | tag       |
    | ce   | liferay/portal  | 7.0.6-ga7 |
    | commerce | liferay/commerce  | 1.1.1 |
    | dxp  | liferay/dxp | 7.0.10.8  |

  Scenario Outline: checki command when an image is not found
    Given I run `lpn checki <type> -t foo`
    Then the output should contain:
    """
    Image has NOT been pulled from Docker Hub
    """
    And the output should contain:
    """
    image="docker.io/<image>"
    """
    And the exit status should be 0
  
  Examples:
    | type    | image |
    | ce | liferay/portal:foo |
    | commerce | liferay/commerce:foo |
    | dxp | liferay/dxp:foo |
    | nightly | mdelapenya/portal-snapshot:foo |
    | release | mdelapenya/liferay-portal:foo |