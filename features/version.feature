Feature: Version command
  As a newcomer to lpn
  I want to be able to check the version

  Scenario: Version command
    Given I run `lpn version`
    Then the exit status should be 0
    And the output should contain:
    """
    0.13.0
    """
    And the output should contain:
    """
    dockerClient=
    """
    And the output should contain:
    """
    dockerServer=
    """
    And the output should contain:
    """
    golang=
    """
    And the output should contain:
    """
    lpn (Liferay Portal Nook) v
    """