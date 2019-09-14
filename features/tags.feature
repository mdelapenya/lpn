Feature: Tags command
  As a newcomer to lpn
  I want to be able to list the available tags for an image type

  Scenario Outline: List available tags without pagination
    When I run `lpn tags <type>`
    Then the output should contain:
    """
    elements in page 1 of
    """
    And the output should contain:
    """
    currentPage=1
    """
    And the output should contain:
    """
    elements=
    """
    And the output should contain:
    """
    images=
    """
    And the output should contain:
    """
    totalPages=
    """
    And the exit status should be 0

  Examples:
    | type    |
    | ce      |
    | commerce |
    | dxp     |
    | nightly |
    | release |

  Scenario Outline: List available tags with pagination
    When I run `lpn tags <type> -p 2 -s 2`
    Then the output should contain:
    """
    elements in page 2 of
    """
    And the output should contain:
    """
    currentPage=2
    """
    And the exit status should be 0

  Examples:
    | type    |
    | ce      |
    | commerce |
    | dxp     |
    | nightly |
    | release |

  Scenario Outline: List available tags with size
    When I run `lpn tags <type> -s 2`
    Then the output should contain:
    """
    showing 2 elements in page 1 of
    """
    And the exit status should be 0

  Examples:
    | type    |
    | ce      |
    | commerce |
    | dxp     |
    | nightly |
    | release |

  Scenario Outline: Inform that the combination of page and size produces a not found resource error
    When I run `lpn tags <type> -p 103`
    Then the output should contain:
    """
    There are no available tags for that pagination. Please use --page and --size arguments to filter properly
    """
    And the exit status should be 0

  Examples:
    | type    |
    | ce      |
    | commerce |
    | dxp     |
    | nightly |
    | release |