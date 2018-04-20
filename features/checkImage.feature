Feature: CheckImage command
  As a newcomer to lpn
  I want to be able to check if the images representing commerce, nightly builds or releases are present

  Scenario Outline: CheckImage command
    Given I run `lpn pull <type> -t latest`
    When I run `lpn checkImage <type> -t latest`
    Then the output should contain:
    """
    The image [<image>] has been pulled from Docker Hub
    """
    And the exit status should be 0

  Examples:
    | type    | image |
    | commerce | liferay/liferay-commerce:latest |
    | nightly | mdelapenya/liferay-portal-nightlies:latest |
    | release | mdelapenya/liferay-portal:latest |

  Scenario Outline: CheckImage command when an image is not found
    When I run `lpn checkImage <type> -t foo`
    Then the output should contain:
    """
    The image [<image>] has NOT been pulled from Docker Hub
    """
    And the exit status should be 1
  
  Examples:
    | type    | iamge |
    | commerce | liferay/liferay-commerce:foo |
    | nightly | mdelapenya/liferay-portal-nightlies:foo |
    | release | mdelapenya/liferay-portal:foo |