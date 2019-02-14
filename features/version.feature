Feature: Version command
  As a newcomer to lpn
  I want to be able to check the version

  Scenario: Version command
    Given I run `lpn version`
    Then the exit status should be 0
    And the output should contain:
    """
    0.8.0
    """
    And the output should contain:
    """
    Client version:
    """
    And the output should contain:
    """
    Server version:
    """
    And the output should contain:
    """
    Go version:
    """
    And the output should contain:
    """
    lpn (Liferay Portal Nook) v
    """