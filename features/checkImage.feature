Feature: CheckImage command
  As a newcomer to lpn
  I want to be able to check if the images representing nightly builds or releases are present

  Scenario Outline: CheckImage command
    Given I run `lpn pull <type> -t latest`
    When I run `lpn checkImage <type> -t latest`
    Then the output should contain:
    """
    <message>
    """
    And the exit status should be 0

  Examples:
    | type    | message |
    | commerce | The image [liferay/liferay-commerce:latest] has been pulled from Docker Hub |
    | nightly | The image [mdelapenya/liferay-portal-nightlies:latest] has been pulled from Docker Hub |
    | release | The image [mdelapenya/liferay-portal:latest] has been pulled from Docker Hub |

  Scenario Outline: CheckImage command when an image is not found
    When I run `lpn checkImage <type> -t foo`
    Then the output should contain:
    """
    <message>
    """
    And the exit status should be 1
  
  Examples:
    | type    | message |
    | commerce | The image [liferay/liferay-commerce:foo] has NOT been pulled from Docker Hub |
    | nightly | The image [mdelapenya/liferay-portal-nightlies:foo] has NOT been pulled from Docker Hub |
    | release | The image [mdelapenya/liferay-portal:foo] has NOT been pulled from Docker Hub |